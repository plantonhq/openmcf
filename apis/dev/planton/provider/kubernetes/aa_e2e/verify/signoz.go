package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// SignozVerifier checks a SigNoz install to the point a customer could
// run their observability on it: the server StatefulSet ready, the
// ingestion collector rolled out, the bundled ClickHouse installation
// reconciled (bundled arm), the server's own health endpoint OK (it
// answers only with a working telemetry-store connection), and a LIVE
// product-grade round-trip — the first admin user REGISTERED through the
// product's own API, a session opened, a span pushed over OTLP/HTTP to
// the collector, and the trace retrieved BY ID through the authenticated
// query API (an observability platform that cannot ingest and answer for
// a trace is not an observability platform).
//
// The behavioral-state scenario (recognized by name) additionally DELETES
// the server pod after the first query, waits for a REPLACEMENT pod (a
// new UID — status flapping Ready on the dying pod is not recovery),
// signs in AGAIN with the same credentials (users/dashboards live in the
// server's SQLite on the PVC — the state proof), and re-queries the trace
// (telemetry lives in ClickHouse — the storage-separation proof).
type SignozVerifier struct {
	Namespace string
	Name      string
	// BundledClickHouse asserts the chart-owned ClickHouseInstallation
	// and the module-owned auth Secret (the composition handle).
	BundledClickHouse bool
	// StateProof switches on the pod-replacement re-login + re-query arm.
	StateProof bool
}

// signozBundledClickHouse reports whether the bundled ClickHouse deploys:
// the external_clickhouse arm absent means the empty-oneof default — the
// bundled appliance (both manifest key forms tolerated).
func signozBundledClickHouse(spec map[string]interface{}) bool {
	if _, ok := spec["external_clickhouse"]; ok {
		return false
	}
	if _, ok := spec["externalClickhouse"]; ok {
		return false
	}
	return true
}

// The verifier-owned first-admin credentials. Registration is the
// product's own first-run flow (the /api/v1/register route is open ONLY
// until the first user exists), so the verifier IS the first customer.
const (
	signozE2eEmail    = "e2e-proof@example.com"
	signozE2ePassword = "E2eProofPassw0rd!"
)

func (v *SignozVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] signoz %q in namespace %q\n", v.Name, v.Namespace)

	// The server (UI/API/ruler/alertmanager in one binary).
	if err := waitStatefulSetReady(ctx, kubeconfig, v.Name, v.Namespace, 15*time.Minute); err != nil {
		return errors.Wrap(err, "the signoz server statefulset never became ready")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "signoz server service not found")
	}

	// The ingestion collector.
	collector := v.Name + "-otel-collector"
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+collector, v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the otel-collector deployment never rolled out")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", collector, v.Namespace); err != nil {
		return errors.Wrap(err, "otel-collector service not found")
	}

	// The bundled ClickHouse: the chart-owned installation object and
	// the module-owned credential Secret (the composition handle the
	// outputs promise).
	if v.BundledClickHouse {
		if err := KubectlResourceExists(ctx, kubeconfig,
			"clickhouseinstallations.clickhouse.altinity.com", v.Name+"-clickhouse", v.Namespace); err != nil {
			return errors.Wrap(err, "the bundled ClickHouseInstallation not found")
		}
		if err := KubectlResourceExists(ctx, kubeconfig,
			"secret", v.Name+"-clickhouse-auth", v.Namespace); err != nil {
			return errors.Wrap(err, "the module-owned clickhouse auth secret not found")
		}
	}

	return v.proveIngestQuery(ctx, kubeconfig)
}

func (v *SignozVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.Name, v.Namespace)
}

// proveIngestQuery runs the product-grade round-trip: health → first-admin
// registration → session → OTLP push → authenticated trace-by-ID. On the
// state-proof lane it kills the server pod between the first and second
// query and re-authenticates.
func (v *SignozVerifier) proveIngestQuery(ctx context.Context, kubeconfig string) error {
	const apiPort = "18080"
	const otlpPort = "14318"

	apiCancel, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name, v.Namespace, apiPort+":8080")
	if err != nil {
		return errors.Wrap(err, "starting port-forward to the signoz API")
	}
	defer apiCancel()

	otlpCancel, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name+"-otel-collector", v.Namespace, otlpPort+":4318")
	if err != nil {
		return errors.Wrap(err, "starting port-forward to the collector's OTLP receiver")
	}
	defer otlpCancel()

	apiBase := "http://127.0.0.1:" + apiPort
	otlpBase := "http://127.0.0.1:" + otlpPort

	// Health answers only with a working telemetry-store connection —
	// the end-to-end wiring proof (server → ClickHouse), bundled or
	// external alike.
	if _, err := httpRoundTrip(ctx, http.MethodGet, apiBase+"/api/v1/health", "", "", 5*time.Minute); err != nil {
		return errors.Wrap(err, "the signoz health endpoint never answered")
	}
	fmt.Printf("  [verify] HEALTH: /api/v1/health OK (telemetry store connected)\n")

	// First-admin registration — the product's own first-run flow.
	registerBody := fmt.Sprintf(`{
      "name": "E2E Proof",
      "email": %q,
      "password": %q,
      "orgDisplayName": "e2e-proof",
      "orgName": "e2e-proof"
    }`, signozE2eEmail, signozE2ePassword)
	if _, err := httpRoundTrip(ctx, http.MethodPost, apiBase+"/api/v1/register",
		"application/json", registerBody, 3*time.Minute); err != nil {
		return errors.Wrap(err, "registering the first admin user")
	}
	fmt.Printf("  [verify] REGISTER: first admin user created through the product API\n")

	token, err := v.login(ctx, apiBase)
	if err != nil {
		return errors.Wrap(err, "opening a session as the registered admin")
	}
	fmt.Printf("  [verify] LOGIN: session opened (access token issued)\n")

	// Push one span with a deterministic trace ID over OTLP/HTTP.
	traceId := fmt.Sprintf("%032x", time.Now().UnixNano())
	spanId := fmt.Sprintf("%016x", time.Now().UnixNano())
	startNs := time.Now().UnixNano()
	pushBody := fmt.Sprintf(`{
      "resourceSpans": [{
        "resource": {"attributes": [{"key": "service.name", "value": {"stringValue": "e2e-proof"}}]},
        "scopeSpans": [{
          "spans": [{
            "traceId": "%s",
            "spanId": "%s",
            "name": "e2e-proof-span",
            "kind": 1,
            "startTimeUnixNano": "%d",
            "endTimeUnixNano": "%d"
          }]
        }]
      }]
    }`, traceId, spanId, startNs, startNs+int64(time.Millisecond))
	if _, err := httpRoundTrip(ctx, http.MethodPost, otlpBase+"/v1/traces",
		"application/json", pushBody, 5*time.Minute); err != nil {
		return errors.Wrap(err, "pushing the proof span to the collector over OTLP")
	}
	fmt.Printf("  [verify] PUSH: proof span (trace %s) accepted by the collector\n", traceId)

	if err := v.queryTrace(ctx, apiBase, token, traceId, spanId, "QUERY"); err != nil {
		return err
	}

	if !v.StateProof {
		return nil
	}

	// The state proof: kill ONLY the server pod (the component label
	// fences the collector and ClickHouse out of the blast radius), wait
	// for a replacement, then re-authenticate and re-query — users
	// survive on the SQLite PVC, telemetry survives in ClickHouse.
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
		"app.kubernetes.io/instance="+v.Name+",app.kubernetes.io/component=signoz", 10*time.Minute); err != nil {
		return errors.Wrap(err, "the signoz server pod did not recover after deletion")
	}
	apiCancel()
	apiCancel2, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name, v.Namespace, apiPort+":8080")
	if err != nil {
		return errors.Wrap(err, "re-establishing the API port-forward after the pod kill")
	}
	defer apiCancel2()

	token, err = v.login(ctx, apiBase)
	if err != nil {
		return errors.Wrap(err, "re-authenticating after the pod replacement — the registered user should have survived on the state volume")
	}
	fmt.Printf("  [verify] STATE: re-login succeeded AFTER pod replacement — users survived on the PVC\n")

	return v.queryTrace(ctx, apiBase, token, traceId, spanId, "PERSISTENCE")
}

// login opens a session through the product's sessions API and returns
// the bearer access token.
func (v *SignozVerifier) login(ctx context.Context, apiBase string) (string, error) {
	loginBody := fmt.Sprintf(`{"email": %q, "password": %q}`, signozE2eEmail, signozE2ePassword)
	deadline := time.Now().Add(3 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		// The sessions route at the pin: POST /api/v2/sessions/email_password
		// (request PostableEmailPasswordSession, response GettableToken —
		// verified in the app source at v0.133.0).
		body, err := httpRoundTrip(ctx, http.MethodPost, apiBase+"/api/v2/sessions/email_password",
			"application/json", loginBody, 30*time.Second)
		if err == nil {
			if token := extractAccessToken(body); token != "" {
				return token, nil
			}
			lastErr = errors.Errorf("session response carried no access token: %s", firstLines(body, 3))
		} else {
			lastErr = err
		}
		time.Sleep(5 * time.Second)
	}
	return "", lastErr
}

// queryTrace retrieves the trace by ID through the authenticated query
// API until the proof span appears (ingestion flushes through ClickHouse
// batches before a trace is queryable).
func (v *SignozVerifier) queryTrace(ctx context.Context, apiBase, token, traceId, spanId, verb string) error {
	deadline := time.Now().Add(6 * time.Minute)
	var lastBody string
	for time.Now().Before(deadline) {
		body, err := httpBearerRoundTrip(ctx, http.MethodGet,
			apiBase+"/api/v1/traces/"+traceId, token, 1*time.Minute)
		if err == nil && strings.Contains(strings.ToLower(body), strings.ToLower(spanId)) {
			suffix := ""
			if verb == "PERSISTENCE" {
				suffix = " AFTER pod replacement — telemetry survived in ClickHouse"
			}
			fmt.Printf("  [verify] %s: trace %s retrieved by ID through the authenticated API%s\n", verb, traceId, suffix)
			return nil
		}
		if body != "" {
			lastBody = body
		} else if err != nil {
			lastBody = err.Error()
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("the proof trace was never retrievable by ID: %s", firstLines(lastBody, 3))
}

// extractAccessToken digs the access token out of the sessions API
// response, tolerating both the {status, data: {...}} envelope and a bare
// object.
func extractAccessToken(body string) string {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return ""
	}
	if token, _ := parsed["accessToken"].(string); token != "" {
		return token
	}
	if data, _ := parsed["data"].(map[string]interface{}); data != nil {
		if token, _ := data["accessToken"].(string); token != "" {
			return token
		}
	}
	return ""
}

// httpBearerRoundTrip is a single authenticated request with a bearer
// token (the shared httpRoundTripAuth speaks only basic auth).
func httpBearerRoundTrip(ctx context.Context, method, url, token string, budget time.Duration) (string, error) {
	deadline := time.Now().Add(budget)
	var lastOut string
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(nil))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			raw, readErr := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if readErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return string(raw), nil
			}
			lastOut = string(raw)
			lastErr = errors.Errorf("HTTP %d: %s", resp.StatusCode, firstLines(string(raw), 2))
		} else {
			lastErr = err
		}
		time.Sleep(5 * time.Second)
	}
	if lastErr == nil {
		lastErr = errors.New("request never succeeded within the budget")
	}
	_ = lastOut
	return lastOut, lastErr
}

package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// SignozVerifier checks a SigNoz install to the point a customer could
// run their observability on it: the server StatefulSet ready, the
// ingestion collector rolled out, the server's own health endpoint OK
// (it answers only with a working connection to the COMPOSED ClickHouse
// — this component installs no database), and a LIVE product-grade
// round-trip — the first admin user REGISTERED through the product's
// own API, a session opened, a span pushed over OTLP/HTTP to the
// collector, and the trace retrieved BY ID through the authenticated
// query API (an observability platform that cannot ingest and answer
// for a trace is not an observability platform).
//
// The behavioral-state scenario (recognized by name) additionally DELETES
// the server pod after the first query, waits for a REPLACEMENT pod (a
// new UID — status flapping Ready on the dying pod is not recovery),
// signs in AGAIN with the same credentials (users/dashboards live in the
// server's SQLite on the PVC — the state proof), and re-queries the trace
// (telemetry lives in the composed ClickHouse — the storage-separation
// proof).
type SignozVerifier struct {
	Namespace string
	Name      string
	// StateProof switches on the pod-replacement re-login + re-query arm.
	StateProof bool
	// Receiver posture, asserted against the collector Service's RENDERED
	// ports (nil = the scenario does not declare the toggle, nothing
	// asserted). The port toggles and the pipeline receiver lists derive
	// from ONE spec field, so a Service port present/absent as declared is
	// the observable half of that single-derivation contract.
	ExpectZipkinPort *bool
	ExpectJaegerPort *bool
}

// SignozReceiverPosture derives the receiver-port expectations from the
// scenario's spec: a toggle the manifest DECLARES becomes an assertion
// against the collector Service's rendered ports; an absent toggle
// asserts nothing (the chart default stays upstream's business). Both
// manifest key forms tolerated.
func SignozReceiverPosture(spec map[string]interface{}) (zipkin, jaeger *bool) {
	collector, _ := spec["otel_collector"].(map[string]interface{})
	if collector == nil {
		collector, _ = spec["otelCollector"].(map[string]interface{})
	}
	if collector == nil {
		return nil, nil
	}
	lookup := func(keys ...string) *bool {
		for _, k := range keys {
			if raw, ok := collector[k]; ok {
				if b, ok := raw.(bool); ok {
					return &b
				}
			}
		}
		return nil
	}
	return lookup("zipkin_receiver_enabled", "zipkinReceiverEnabled"),
		lookup("jaeger_receiver_enabled", "jaegerReceiverEnabled")
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
	if err := v.assertReceiverPorts(ctx, kubeconfig, collector); err != nil {
		return err
	}

	return v.proveIngestQuery(ctx, kubeconfig)
}

// assertReceiverPorts checks the collector Service's rendered port NAMES
// against the scenario's declared receiver toggles. The chart names each
// Service port after its values key (zipkin, jaeger-thrift, jaeger-grpc,
// jaeger-compact), so a toggle that failed to move the port map is
// directly observable here — and because the pipeline receiver lists
// derive from the same spec field, this is the rendered-surface half of
// the "ports and pipelines move together" contract.
func (v *SignozVerifier) assertReceiverPorts(ctx context.Context, kubeconfig, service string) error {
	if v.ExpectZipkinPort == nil && v.ExpectJaegerPort == nil {
		return nil
	}
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "service", service, "-n", v.Namespace,
		"-o", "jsonpath={.spec.ports[*].name}").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading the collector service ports: %s", string(out))
	}
	ports := strings.Fields(string(out))
	has := func(prefix string) bool {
		for _, p := range ports {
			if strings.HasPrefix(p, prefix) {
				return true
			}
		}
		return false
	}
	if v.ExpectZipkinPort != nil && has("zipkin") != *v.ExpectZipkinPort {
		return errors.Errorf("receiver posture violated: zipkin port present=%v, declared enabled=%v (rendered ports: %v)",
			has("zipkin"), *v.ExpectZipkinPort, ports)
	}
	if v.ExpectJaegerPort != nil && has("jaeger") != *v.ExpectJaegerPort {
		return errors.Errorf("receiver posture violated: jaeger port present=%v, declared enabled=%v (rendered ports: %v)",
			has("jaeger"), *v.ExpectJaegerPort, ports)
	}
	fmt.Printf("  [verify] RECEIVERS: collector service ports match the declared posture (%v)\n", ports)
	return nil
}

// VerifyAbsent asserts the release's workloads are gone. The release
// ships no operator and no CRDs (the composed ClickHouse is a separate
// component with its own lifecycle), so uninstall is ordinary object
// deletion — nothing survives by design.
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

	apiBase := "http://127.0.0.1:" + apiPort
	otlpBase := "http://127.0.0.1:" + otlpPort

	// Health answers only with a working telemetry-store connection —
	// the end-to-end wiring proof (server → ClickHouse), bundled or
	// external alike.
	if _, err := httpRoundTrip(ctx, http.MethodGet, apiBase+"/api/v1/health", "", "", 5*time.Minute); err != nil {
		return errors.Wrap(err, "the signoz health endpoint never answered")
	}
	fmt.Printf("  [verify] HEALTH: /api/v1/health OK (telemetry store connected)\n")

	// First-admin registration — the product's own first-run flow. A
	// register failure falls through to login instead of failing the
	// lane: the route flips SetupCompleted the moment the first user row
	// commits, so a retry after a response lost mid-read answers 4xx
	// "self-registration is disabled" forever after — the only honest
	// test of whether that first attempt actually took is a login with
	// the same credentials (verified against the app source at the pin).
	registerBody := fmt.Sprintf(`{
      "name": "E2E Proof",
      "email": %q,
      "password": %q,
      "orgDisplayName": "e2e-proof",
      "orgName": "e2e-proof"
    }`, signozE2eEmail, signozE2ePassword)
	_, registerErr := httpRoundTrip(ctx, http.MethodPost, apiBase+"/api/v1/register",
		"application/json", registerBody, 3*time.Minute)
	if registerErr == nil {
		fmt.Printf("  [verify] REGISTER: first admin user created through the product API\n")
	} else {
		fmt.Printf("  [verify] REGISTER: did not complete cleanly — attempting login with the same credentials (a partial success disables re-registration)\n")
	}

	token, err := v.login(ctx, apiBase)
	if err != nil {
		if registerErr != nil {
			return errors.Wrapf(registerErr, "registering the first admin user (login with the same credentials also failed: %v)", err)
		}
		return errors.Wrap(err, "opening a session as the registered admin")
	}
	fmt.Printf("  [verify] LOGIN: session opened (access token issued)\n")

	// Push one span with a deterministic trace ID over OTLP/HTTP. The
	// collector tunnel is established FRESH here — never earlier, and
	// re-established between attempts: kubectl port-forward is a
	// single-pod tunnel that dies silently when its pod restarts or the
	// stream drops, and a tunnel opened before the minutes of health/
	// register/login above was dead by push time (verified live — every
	// retry read connection refused on the local port).
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
	pushDeadline := time.Now().Add(5 * time.Minute)
	for {
		var pushErr error
		otlpCancel, err := startPortForward(ctx, kubeconfig,
			"svc/"+v.Name+"-otel-collector", v.Namespace, otlpPort+":4318")
		if err != nil {
			pushErr = errors.Wrap(err, "starting port-forward to the collector's OTLP receiver")
		} else {
			_, pushErr = httpRoundTrip(ctx, http.MethodPost, otlpBase+"/v1/traces",
				"application/json", pushBody, 30*time.Second)
			otlpCancel()
		}
		if pushErr == nil {
			break
		}
		if time.Now().After(pushDeadline) {
			return errors.Wrap(pushErr, "pushing the proof span to the collector over OTLP")
		}
		time.Sleep(10 * time.Second)
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
// the bearer access token. Login is a TWO-step flow at the pin,
// mirroring the product's own login page: resolve the user's org
// through the OPEN session-context route, then post the email_password
// session WITH that orgId — PostableEmailPasswordSession rejects a
// missing orgId at unmarshal ("orgID is required", caught live; the
// request type's own UnmarshalJSON enforces it before authentication
// runs).
func (v *SignozVerifier) login(ctx context.Context, apiBase string) (string, error) {
	deadline := time.Now().Add(3 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		orgID, err := v.sessionOrgID(ctx, apiBase)
		if err != nil {
			lastErr = errors.Wrap(err, "resolving the org through the session-context route")
			time.Sleep(5 * time.Second)
			continue
		}
		loginBody := fmt.Sprintf(`{"email": %q, "password": %q, "orgId": %q}`,
			signozE2eEmail, signozE2ePassword, orgID)
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

// sessionOrgID resolves the org the e2e user belongs to through the
// open session-context route (GET /api/v2/sessions/context?email=&ref=)
// — the same pre-login discovery the product's login page performs.
// The ref parameter only needs to parse as a URL.
func (v *SignozVerifier) sessionOrgID(ctx context.Context, apiBase string) (string, error) {
	u := apiBase + "/api/v2/sessions/context?email=" + url.QueryEscape(signozE2eEmail) +
		"&ref=" + url.QueryEscape(apiBase)
	body, err := httpRoundTrip(ctx, http.MethodGet, u, "", "", 30*time.Second)
	if err != nil {
		return "", err
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return "", errors.Wrap(err, "parsing the session context")
	}
	payload := parsed
	if data, _ := parsed["data"].(map[string]interface{}); data != nil {
		payload = data
	}
	orgs, _ := payload["orgs"].([]interface{})
	if len(orgs) == 0 {
		return "", errors.Errorf("session context lists no orgs for the e2e user: %s", firstLines(body, 3))
	}
	first, _ := orgs[0].(map[string]interface{})
	id, _ := first["id"].(string)
	if id == "" {
		return "", errors.Errorf("session context org carried no id: %s", firstLines(body, 3))
	}
	return id, nil
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

package verify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TrinoVerifier checks a Trino cluster to the point a customer could
// point a SQL client at it: coordinator and worker rollouts, THE AUTH
// GATE (an anonymous statement submission REJECTED — upstream ships NO
// authentication and that posture never deploys from this kind), and
// THE QUERY PROOF on every lane — a real query executes through the
// coordinator's own statement API as the module-generated admin and
// returns the expected result through REAL workers (an engine that
// cannot answer SQL is not a query engine).
//
// The full-surface scenario (recognized by name) adds THE FEDERATION
// PROOF: the composed PostgreSQL catalog answers through Trino, and a
// CROSS-CATALOG query joins the in-image tpch catalog against the
// composed database in one statement — the product's whole point,
// proven in one query.
//
// The behavioral scenario (recognized by name) adds THE RECOVERY
// PROOF: the coordinator pod is deleted (UID-verified replacement) and
// a fresh session answers the same query — the engine is stateless by
// design and recovery is a property, not an accident.
//
// Destroy is clean by design: a plain Helm release (no CRDs) plus the
// module-owned Secrets — everything leaves with the resource.
type TrinoVerifier struct {
	Namespace string
	Name      string
	// AdminUsername logs into the statement API (spec default "trino").
	AdminUsername string
	// AuthEnabled gates the credentialed arms (the auth-disabled shape
	// is never a shipped scenario, but the verifier stays honest).
	AuthEnabled bool
	// WorkerReplicas gates the worker-rollout wait (a zero-worker
	// single-node shape has no worker Deployment to wait for).
	WorkerReplicas int
	// FederationProof switches the composed-catalog arms on
	// (full-surface lanes); CatalogName names the composed catalog.
	FederationProof bool
	CatalogName     string
	// RecoveryProof switches the behavioral arm on.
	RecoveryProof bool
}

// trinoAdminUsername reads spec.auth.admin_username ("" = "trino").
func trinoAdminUsername(spec map[string]interface{}) string {
	if auth, ok := spec["auth"].(map[string]interface{}); ok {
		for _, key := range []string{"admin_username", "adminUsername"} {
			if username, ok := auth[key].(string); ok && username != "" {
				return username
			}
		}
	}
	return "trino"
}

// trinoAuthEnabled reads spec.auth.enabled (absent = true — the secured
// default).
func trinoAuthEnabled(spec map[string]interface{}) bool {
	if auth, ok := spec["auth"].(map[string]interface{}); ok {
		if enabled, ok := auth["enabled"].(bool); ok {
			return enabled
		}
	}
	return true
}

// trinoWorkerReplicas reads spec.workers.replicas (absent = 2 — the
// chart default the module keeps).
func trinoWorkerReplicas(spec map[string]interface{}) int {
	if workers, ok := spec["workers"].(map[string]interface{}); ok {
		if replicas, ok := workers["replicas"].(float64); ok {
			return int(replicas)
		}
	}
	return 2
}

const trinoApiPort = "18083"

func trinoBaseUrl() string {
	return "http://127.0.0.1:" + trinoApiPort
}

func (v *TrinoVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] trino %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-coordinator", v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the trino coordinator never rolled out")
	}
	if v.WorkerReplicas > 0 {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-worker", v.Namespace, 10*time.Minute); err != nil {
			return errors.Wrap(err, "the trino workers never rolled out")
		}
	}

	cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name+"-coordinator", trinoApiPort, "8080")
	if err != nil {
		return errors.Wrap(err, "opening the tunnel to the trino coordinator")
	}

	password := ""
	if v.AuthEnabled {
		// THE AUTH GATE before any credentialed call.
		if err := v.proveAuthGate(ctx); err != nil {
			cancel()
			return err
		}
		password, err = v.adminPassword(ctx, kubeconfig)
		if err != nil {
			cancel()
			return err
		}
	}

	// THE QUERY PROOF: the in-image tpch catalog's tiny nation table
	// has exactly 25 rows — a full plan/schedule/execute round-trip
	// through real workers.
	if err := v.proveQuery(ctx, password,
		"SELECT count(*) FROM tpch.tiny.nation", "25", "the tpch query proof"); err != nil {
		cancel()
		return err
	}
	fmt.Printf("  [verify] THE QUERY PROOF: tpch.tiny.nation answered 25 rows through the statement API\n")

	if v.FederationProof {
		catalog := v.CatalogName
		// THE FEDERATION PROOF: the composed PostgreSQL catalog
		// answers through Trino (information_schema always exists —
		// no test data required)…
		schemataCount, err := v.runScalarQuery(ctx, password,
			fmt.Sprintf("SELECT count(*) FROM %s.information_schema.schemata", catalog))
		if err != nil {
			cancel()
			return errors.Wrapf(err, "the federation proof: querying the composed %q catalog", catalog)
		}
		fmt.Printf("  [verify] THE FEDERATION PROOF: composed catalog %q answered (%s schemata)\n", catalog, schemataCount)
		// …and one CROSS-CATALOG statement joins the in-image sample
		// data against the composed database — federation in a single
		// query, the product's whole point.
		if _, err := v.runScalarQuery(ctx, password,
			fmt.Sprintf("SELECT count(*) FROM tpch.tiny.nation n CROSS JOIN %s.information_schema.schemata s", catalog)); err != nil {
			cancel()
			return errors.Wrap(err, "the cross-catalog join proof")
		}
		fmt.Printf("  [verify] THE FEDERATION PROOF: cross-catalog join (tpch × %s) answered\n", catalog)
	}

	if !v.RecoveryProof {
		cancel()
		return nil
	}

	// THE RECOVERY PROOF: drop the tunnel across the replacement
	// window (fresh-tunnel-per-phase), replace the coordinator pod
	// UID-verified, then answer the same query from a fresh session —
	// the engine is stateless by design.
	cancel()
	fmt.Printf("  [verify] THE RECOVERY PROOF: replacing the coordinator pod…\n")
	selector := "app.kubernetes.io/instance=" + v.Name + ",app.kubernetes.io/component=coordinator"
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace, selector, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the coordinator pod was never replaced")
	}
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-coordinator", v.Namespace, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the replaced coordinator never became ready")
	}
	cancel, err = openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name+"-coordinator", trinoApiPort, "8080")
	if err != nil {
		return errors.Wrap(err, "re-opening the tunnel after the coordinator replacement")
	}
	defer cancel()
	if err := v.proveQuery(ctx, password,
		"SELECT count(*) FROM tpch.tiny.nation", "25", "the post-replacement query proof"); err != nil {
		return err
	}
	fmt.Printf("  [verify] THE RECOVERY PROOF: the replaced coordinator answered the same query\n")
	return nil
}

// VerifyAbsent asserts destroy left nothing behind: a plain Helm
// release plus module-owned Secrets — no CRDs, no designed survivors.
func (v *TrinoVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify-absent] trino %q in namespace %q\n", v.Name, v.Namespace)
	for _, deployment := range []string{v.Name + "-coordinator", v.Name + "-worker"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", deployment, v.Namespace); err != nil {
			return err
		}
	}
	for _, secret := range []string{v.Name + "-auth", v.Name + "-internal"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "secret", secret, v.Namespace); err != nil {
			return err
		}
	}
	return nil
}

// adminPassword reads the module-generated admin credential from the
// exported `<name>-auth` Secret (key password).
func (v *TrinoVerifier) adminPassword(ctx context.Context, kubeconfig string) (string, error) {
	secretName := v.Name + "-auth"
	b64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.password}")
	if err != nil {
		return "", errors.Wrapf(err, "reading secret %q key password", secretName)
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	password := strings.TrimSpace(string(raw))
	if password == "" {
		return "", errors.New("admin password Secret was empty")
	}
	return password, nil
}

// proveAuthGate submits a statement WITHOUT credentials and requires a
// 401 — upstream's open server (any caller queries as any user) must
// be dead.
func (v *TrinoVerifier) proveAuthGate(ctx context.Context) error {
	status, _, err := v.statementRequest(ctx, "SELECT 1", "", "")
	if err != nil {
		return errors.Wrap(err, "the anonymous statement submission")
	}
	if status != http.StatusUnauthorized {
		return errors.Errorf("THE AUTH GATE FAILED: anonymous statement submission returned %d, want 401 — the open server posture must never ship", status)
	}
	fmt.Printf("  [verify] THE AUTH GATE: anonymous statement submission rejected (401)\n")
	return nil
}

// proveQuery runs a scalar query and asserts the expected value.
func (v *TrinoVerifier) proveQuery(ctx context.Context, password, query, want, label string) error {
	got, err := v.runScalarQuery(ctx, password, query)
	if err != nil {
		return errors.Wrap(err, label)
	}
	if got != want {
		return errors.Errorf("%s: query %q answered %q, want %q", label, query, got, want)
	}
	return nil
}

// runScalarQuery executes one query through the statement API (POST
// /v1/statement, then follow nextUri until FINISHED) and returns the
// single scalar result. Retries the submission during coordinator
// warm-up.
func (v *TrinoVerifier) runScalarQuery(ctx context.Context, password, query string) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		if attempt > 0 {
			time.Sleep(6 * time.Second)
		}
		value, err := v.runScalarQueryOnce(ctx, password, query)
		if err == nil {
			return value, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func (v *TrinoVerifier) runScalarQueryOnce(ctx context.Context, password, query string) (string, error) {
	status, body, err := v.statementRequest(ctx, query, v.AdminUsername, password)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", errors.Errorf("statement submission returned %d: %s", status, trinoTruncate(body))
	}

	for {
		var page struct {
			NextUri string          `json:"nextUri"`
			Data    [][]interface{} `json:"data"`
			Stats   struct {
				State string `json:"state"`
			} `json:"stats"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal([]byte(body), &page); err != nil {
			return "", errors.Wrap(err, "decoding the statement page")
		}
		if page.Error != nil {
			return "", errors.Errorf("the query failed: %s", page.Error.Message)
		}
		if len(page.Data) > 0 && len(page.Data[0]) > 0 {
			return fmt.Sprintf("%v", page.Data[0][0]), nil
		}
		if page.NextUri == "" {
			if page.Stats.State == "FINISHED" {
				return "", errors.New("the query finished without data")
			}
			return "", errors.Errorf("the query ended in state %s without data", page.Stats.State)
		}
		// Follow the coordinator-relative nextUri through the tunnel.
		nextPath := page.NextUri
		if idx := strings.Index(nextPath, "/v1/"); idx >= 0 {
			nextPath = nextPath[idx:]
		}
		status, body, err = v.request(ctx, http.MethodGet, nextPath, v.AdminUsername, password, nil)
		if err != nil {
			return "", err
		}
		if status != http.StatusOK {
			return "", errors.Errorf("following nextUri returned %d: %s", status, trinoTruncate(body))
		}
	}
}

// statementRequest POSTs one query to /v1/statement (basic auth when a
// username is given; X-Trino-User always names the session user).
func (v *TrinoVerifier) statementRequest(ctx context.Context, query, username, password string) (int, string, error) {
	return v.request(ctx, http.MethodPost, "/v1/statement", username, password, []byte(query))
}

func (v *TrinoVerifier) request(ctx context.Context, method, path, username, password string, payload []byte) (int, string, error) {
	req, err := http.NewRequestWithContext(ctx, method, trinoBaseUrl()+path, strings.NewReader(string(payload)))
	if err != nil {
		return 0, "", err
	}
	if username != "" {
		req.SetBasicAuth(username, password)
		req.Header.Set("X-Trino-User", username)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(responseBody), nil
}

// trinoTruncate keeps error payloads readable.
func trinoTruncate(body string) string {
	if len(body) > 300 {
		return body[:300] + "…"
	}
	return body
}

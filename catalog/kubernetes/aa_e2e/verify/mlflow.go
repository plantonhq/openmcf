package verify

import (
	"bytes"
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

// MlflowVerifier checks an MLflow tracking server to the point a
// customer could point MLFLOW_TRACKING_URI at it: the server rolled out,
// its own /health contract answering, THE AUTH GATE (an anonymous
// tracking-API read REJECTED, and upstream's public default password
// asserted DEAD — the module-generated credential is the only way in),
// and THE TRACKING PROOF on every lane — a real experiment is created
// through the REST API as the admin user, a run logs a parameter and a
// metric, an artifact uploads THROUGH THE SERVER'S OWN PROXY (the
// credential-free client story) and downloads back byte-identical, and
// the run finishes FINISHED (a tracker that cannot track a run is not a
// tracker). Proof artifacts stay inside the lane's own experiment.
//
// The behavioral-durability scenario (recognized by name) adds THE
// STATE PROOF: the server pod is deleted (UID-verified replacement) and
// a fresh session must find the same experiment, the same run's metric,
// and re-download the same artifact bytes — tracking state lives in the
// composed PostgreSQL and the artifacts in the composed object store,
// never in the pod.
//
// Destroy is clean by design: MLflow is module-owned manifests —
// everything the module created (Deployment, Service, Secrets, PVCs,
// the gc CronJob, the ServiceMonitor) leaves with the resource; no CRDs.
type MlflowVerifier struct {
	Namespace string
	Name      string
	// AdminUsername logs into the API (spec default "admin").
	AdminUsername string
	// AuthEnabled gates the credentialed arms (the auth-disabled shape
	// is never a shipped scenario, but the verifier stays honest).
	AuthEnabled bool
	// StateProof switches the behavioral arm on.
	StateProof bool
}

// mlflowAdminUsername reads spec.auth.admin_username ("" = "admin").
func mlflowAdminUsername(spec map[string]interface{}) string {
	if auth, ok := spec["auth"].(map[string]interface{}); ok {
		for _, key := range []string{"admin_username", "adminUsername"} {
			if username, ok := auth[key].(string); ok && username != "" {
				return username
			}
		}
	}
	return "admin"
}

// mlflowAuthEnabled reads spec.auth.enabled (absent = true — the
// secured default).
func mlflowAuthEnabled(spec map[string]interface{}) bool {
	if auth, ok := spec["auth"].(map[string]interface{}); ok {
		if enabled, ok := auth["enabled"].(bool); ok {
			return enabled
		}
	}
	return true
}

const mlflowApiPort = "18082"

func mlflowBaseUrl() string {
	return "http://127.0.0.1:" + mlflowApiPort
}

func (v *MlflowVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] mlflow %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the mlflow deployment never rolled out")
	}

	cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name, mlflowApiPort, "5000")
	if err != nil {
		return errors.Wrap(err, "opening the tunnel to the mlflow service")
	}

	// The server's own health contract answers unauthenticated even
	// with auth on.
	if err := v.proveHealth(ctx); err != nil {
		cancel()
		return err
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

	// THE TRACKING PROOF.
	experimentName := "e2e-proof-" + v.Name
	artifactBytes := []byte("mlflow-e2e-artifact-" + v.Name + "\n")
	runId, err := v.proveTracking(ctx, password, experimentName, artifactBytes, "first")
	if err != nil {
		cancel()
		return err
	}

	if !v.StateProof {
		cancel()
		return nil
	}

	// THE STATE PROOF: drop the tunnel across the replacement window
	// (fresh-tunnel-per-phase), replace the server pod UID-verified,
	// then re-read everything from a fresh session.
	cancel()
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
		"app.kubernetes.io/name=mlflow,app.kubernetes.io/instance="+v.Name, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the mlflow pod did not recover after deletion")
	}
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the mlflow deployment never rolled out after the replacement")
	}

	cancel, err = openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name, mlflowApiPort, "5000")
	if err != nil {
		return errors.Wrap(err, "re-establishing the tunnel after the replacement")
	}
	defer cancel()

	if err := v.proveStateSurvived(ctx, password, experimentName, runId, artifactBytes); err != nil {
		return err
	}
	fmt.Printf("  [verify] STATE: the experiment, its run's metric and its artifact bytes all survived a UID-verified pod replacement — tracking state lives in the database and the object store\n")
	return nil
}

func (v *MlflowVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return err
	}
	fmt.Printf("  [verify] DESTROY: the mlflow deployment and service are gone (module-owned manifests, no CRDs — destroy is clean by design)\n")
	return nil
}

// adminPassword reads the module-generated admin credential from the
// exported `<name>-admin-auth` Secret (key password) and asserts
// upstream's public default is dead.
func (v *MlflowVerifier) adminPassword(ctx context.Context, kubeconfig string) (string, error) {
	secretName := v.Name + "-admin-auth"
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
	// Upstream's public example credential must never ship.
	if password == "password1234" {
		return "", errors.New("admin password is upstream's public default — the module-generated Secret was not wired")
	}
	return password, nil
}

// mlflowRequest performs one HTTP request with optional basic auth.
func mlflowRequest(ctx context.Context, method, path, username, password string, payload []byte, contentType string) (int, string, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	request, err := http.NewRequestWithContext(ctx, method, mlflowBaseUrl()+path, body)
	if err != nil {
		return 0, "", err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if username != "" {
		request.SetBasicAuth(username, password)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return 0, "", err
	}
	responseBody, _ := io.ReadAll(response.Body)
	response.Body.Close()
	return response.StatusCode, string(responseBody), nil
}

// proveHealth waits for the server's own /health contract.
func (v *MlflowVerifier) proveHealth(ctx context.Context) error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastStatus int
	for time.Now().Before(deadline) {
		status, _, err := mlflowRequest(ctx, "GET", "/health", "", "", nil, "")
		if err == nil && status == http.StatusOK {
			fmt.Printf("  [verify] HEALTH: /health answers OK\n")
			return nil
		}
		if err == nil {
			lastStatus = status
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("/health never answered OK (last status %d)", lastStatus)
}

// proveAuthGate asserts an anonymous tracking-API read is rejected AND
// upstream's public default credential is dead on the wire.
func (v *MlflowVerifier) proveAuthGate(ctx context.Context) error {
	status, _, err := mlflowRequest(ctx, "GET", "/api/2.0/mlflow/experiments/search?max_results=1", "", "", nil, "")
	if err != nil {
		return errors.Wrap(err, "the anonymous tracking-API probe never answered")
	}
	if status != http.StatusUnauthorized && status != http.StatusForbidden {
		return errors.Errorf("THE AUTH GATE FAILED: anonymous experiments/search answered %d, expected 401/403", status)
	}
	// The upstream example credential must be dead even as a guess.
	status, _, err = mlflowRequest(ctx, "GET", "/api/2.0/mlflow/experiments/search?max_results=1", "admin", "password1234", nil, "")
	if err != nil {
		return errors.Wrap(err, "the default-credential probe never answered")
	}
	if status == http.StatusOK {
		return errors.New("THE AUTH GATE FAILED: upstream's admin/password1234 default signed in")
	}
	fmt.Printf("  [verify] AUTH GATE: anonymous read rejected and upstream's default credential is dead\n")
	return nil
}

// proveTracking runs the full tracking round-trip and returns the run id.
func (v *MlflowVerifier) proveTracking(ctx context.Context, password, experimentName string, artifactBytes []byte, phase string) (string, error) {
	// Create (or reuse) the lane's experiment.
	experimentId, err := v.ensureExperiment(ctx, password, experimentName)
	if err != nil {
		return "", err
	}

	// Create a run.
	payload, _ := json.Marshal(map[string]interface{}{
		"experiment_id": experimentId,
		"run_name":      "e2e-" + phase,
		"start_time":    time.Now().UnixMilli(),
	})
	status, body, err := mlflowRequest(ctx, "POST", "/api/2.0/mlflow/runs/create", v.AdminUsername, password, payload, "application/json")
	if err != nil || status != http.StatusOK {
		return "", errors.Errorf("creating a run failed (%d): %s", status, firstLines(body, 3))
	}
	var runResponse struct {
		Run struct {
			Info struct {
				RunId string `json:"run_id"`
			} `json:"info"`
		} `json:"run"`
	}
	if err := json.Unmarshal([]byte(body), &runResponse); err != nil || runResponse.Run.Info.RunId == "" {
		return "", errors.Errorf("the run-create response carried no run_id: %s", firstLines(body, 3))
	}
	runId := runResponse.Run.Info.RunId

	// Log a parameter and a metric.
	paramPayload, _ := json.Marshal(map[string]interface{}{
		"run_id": runId, "key": "e2e_phase", "value": phase,
	})
	if status, body, err = mlflowRequest(ctx, "POST", "/api/2.0/mlflow/runs/log-parameter", v.AdminUsername, password, paramPayload, "application/json"); err != nil || status != http.StatusOK {
		return "", errors.Errorf("logging a parameter failed (%d): %s", status, firstLines(body, 2))
	}
	metricPayload, _ := json.Marshal(map[string]interface{}{
		"run_id": runId, "key": "e2e_accuracy", "value": 0.99,
		"timestamp": time.Now().UnixMilli(), "step": 1,
	})
	if status, body, err = mlflowRequest(ctx, "POST", "/api/2.0/mlflow/runs/log-metric", v.AdminUsername, password, metricPayload, "application/json"); err != nil || status != http.StatusOK {
		return "", errors.Errorf("logging a metric failed (%d): %s", status, firstLines(body, 2))
	}

	// THE ARTIFACT PROOF through the server's own proxy: upload, then
	// download byte-identical — the credential-free client story.
	artifactPath := "/api/2.0/mlflow-artifacts/artifacts/" + experimentId + "/" + runId + "/artifacts/e2e-proof.txt"
	if status, body, err = mlflowRequest(ctx, "PUT", artifactPath, v.AdminUsername, password, artifactBytes, "application/octet-stream"); err != nil || status != http.StatusOK {
		return "", errors.Errorf("the proxied artifact upload failed (%d): %s", status, firstLines(body, 3))
	}
	status, downloaded, err := mlflowRequest(ctx, "GET", artifactPath, v.AdminUsername, password, nil, "")
	if err != nil || status != http.StatusOK {
		return "", errors.Errorf("the proxied artifact download failed (%d)", status)
	}
	if downloaded != string(artifactBytes) {
		return "", errors.New("the downloaded artifact differs from the uploaded bytes")
	}

	// Finish the run.
	finishPayload, _ := json.Marshal(map[string]interface{}{
		"run_id": runId, "status": "FINISHED", "end_time": time.Now().UnixMilli(),
	})
	if status, body, err = mlflowRequest(ctx, "POST", "/api/2.0/mlflow/runs/update", v.AdminUsername, password, finishPayload, "application/json"); err != nil || status != http.StatusOK {
		return "", errors.Errorf("finishing the run failed (%d): %s", status, firstLines(body, 2))
	}

	fmt.Printf("  [verify] TRACKING (%s): experiment %q — run %s logged a param + metric, an artifact round-tripped through the server proxy, and the run FINISHED\n",
		phase, experimentName, runId)
	return runId, nil
}

// ensureExperiment creates the experiment or resolves an existing one
// by name.
func (v *MlflowVerifier) ensureExperiment(ctx context.Context, password, experimentName string) (string, error) {
	payload, _ := json.Marshal(map[string]string{"name": experimentName})
	status, body, err := mlflowRequest(ctx, "POST", "/api/2.0/mlflow/experiments/create", v.AdminUsername, password, payload, "application/json")
	if err != nil {
		return "", errors.Wrap(err, "the experiment-create request never answered")
	}
	if status == http.StatusOK {
		var created struct {
			ExperimentId string `json:"experiment_id"`
		}
		if json.Unmarshal([]byte(body), &created) == nil && created.ExperimentId != "" {
			return created.ExperimentId, nil
		}
	}
	// Already exists (a re-run) — resolve by name.
	status, body, err = mlflowRequest(ctx, "GET", "/api/2.0/mlflow/experiments/get-by-name?experiment_name="+experimentName, v.AdminUsername, password, nil, "")
	if err != nil || status != http.StatusOK {
		return "", errors.Errorf("resolving experiment %q failed (%d): %s", experimentName, status, firstLines(body, 3))
	}
	var got struct {
		Experiment struct {
			ExperimentId string `json:"experiment_id"`
		} `json:"experiment"`
	}
	if err := json.Unmarshal([]byte(body), &got); err != nil || got.Experiment.ExperimentId == "" {
		return "", errors.Errorf("the get-by-name response carried no experiment_id: %s", firstLines(body, 3))
	}
	return got.Experiment.ExperimentId, nil
}

// proveStateSurvived re-reads the pre-replacement state from a fresh
// session: the experiment resolves, the run's metric answers, and the
// artifact bytes match.
func (v *MlflowVerifier) proveStateSurvived(ctx context.Context, password, experimentName, runId string, artifactBytes []byte) error {
	if err := v.proveHealth(ctx); err != nil {
		return err
	}
	experimentId, err := v.ensureExperiment(ctx, password, experimentName)
	if err != nil {
		return errors.Wrap(err, "the experiment did not survive the replacement")
	}
	status, body, err := mlflowRequest(ctx, "GET", "/api/2.0/mlflow/runs/get?run_id="+runId, v.AdminUsername, password, nil, "")
	if err != nil || status != http.StatusOK {
		return errors.Errorf("the run did not survive the replacement (%d): %s", status, firstLines(body, 3))
	}
	if !strings.Contains(body, "e2e_accuracy") {
		return errors.New("the run survived but its metric is missing")
	}
	artifactPath := "/api/2.0/mlflow-artifacts/artifacts/" + experimentId + "/" + runId + "/artifacts/e2e-proof.txt"
	status, downloaded, err := mlflowRequest(ctx, "GET", artifactPath, v.AdminUsername, password, nil, "")
	if err != nil || status != http.StatusOK {
		return errors.Errorf("the artifact did not survive the replacement (%d)", status)
	}
	if downloaded != string(artifactBytes) {
		return errors.New("the artifact survived but its bytes differ")
	}
	return nil
}

package verify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// AirflowVerifier checks an Airflow installation to the point a customer
// could run their pipelines on it: the Airflow 3 component set rolled
// out (API server, scheduler, dag-processor, triggerer — plus the Celery
// workers and bundled Redis when a Celery executor is declared), the
// API server's own health contract answering (metadatabase + scheduler
// healthy — the composed database wiring proven end to end), THE AUTH
// GATE (an unauthenticated API read REJECTED — the module-generated
// admin credential is the only way in, and the chart's admin/admin
// default must be dead), and THE DAG PROOF on every lane — a real DAG
// is unpaused, triggered through the REST API as the admin user, and
// polled to a SUCCESSFUL DagRun (an orchestrator that cannot run a DAG
// to success is not an orchestrator). On KubernetesExecutor lanes the
// proof exercises the whole story: API → scheduler → per-task pods →
// state in the composed PostgreSQL; on Celery lanes it additionally
// proves broker → worker dispatch.
//
// The behavioral-dag-delivery scenario (recognized by name) runs THE
// DELIVERY PROOF instead of the example DAG: the verifier writes its
// own marker DAG into the shared dags volume (through the dag-processor
// pod — the parser's own mount), waits for the DAG PROCESSOR to parse
// it into the API, triggers it to success, then DELETES the scheduler
// pod (UID-verified replacement) and triggers a SECOND run — scheduling
// state lives in the database, never in the pod.
//
// Destroy is clean by design: Airflow installs no CRDs — everything
// leaves with the release (the module-owned credential Secrets travel
// with the module).
type AirflowVerifier struct {
	Namespace string
	Name      string
	// CeleryEnabled gates the worker/redis rollout checks (derived
	// from the manifest's executor — the chart's own substring test).
	CeleryEnabled bool
	// BundledRedis gates the Redis StatefulSet rollout check.
	BundledRedis bool
	// AdminUsername logs into the API (spec default "admin").
	AdminUsername string
	// DeliveryProof switches the behavioral arm on.
	DeliveryProof bool
}

// airflowExecutor reads spec.executor ("" = the spec default,
// KubernetesExecutor).
func airflowExecutor(spec map[string]interface{}) string {
	if executor, ok := spec["executor"].(string); ok {
		return executor
	}
	return "KubernetesExecutor"
}

// airflowCeleryEnabled mirrors the chart's own pairing test.
func airflowCeleryEnabled(spec map[string]interface{}) bool {
	executor := airflowExecutor(spec)
	return strings.Contains(executor, "CeleryExecutor") || strings.Contains(executor, "CeleryKubernetesExecutor")
}

// airflowBundledRedis reads whether the broker's bundled_redis arm is
// declared (both manifest key forms tolerated).
func airflowBundledRedis(spec map[string]interface{}) bool {
	for _, brokerKey := range []string{"broker"} {
		if broker, ok := spec[brokerKey].(map[string]interface{}); ok {
			for _, key := range []string{"bundled_redis", "bundledRedis"} {
				if _, ok := broker[key]; ok {
					return true
				}
			}
		}
	}
	return false
}

// airflowAdminUsername reads spec.admin_user.username ("" = "admin").
func airflowAdminUsername(spec map[string]interface{}) string {
	for _, adminKey := range []string{"admin_user", "adminUser"} {
		if admin, ok := spec[adminKey].(map[string]interface{}); ok {
			if username, ok := admin["username"].(string); ok && username != "" {
				return username
			}
		}
	}
	return "admin"
}

const airflowApiPort = "18080"

func (v *AirflowVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] airflow %q in namespace %q\n", v.Name, v.Namespace)

	// The Airflow 3 component set rolls out. The first wait absorbs the
	// image pull + migration-Job budget (the post-install hooks run
	// before Helm's wait returns, but rollout re-asserts from here).
	for _, deployment := range []string{"api-server", "scheduler", "dag-processor"} {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-"+deployment, v.Namespace, 15*time.Minute); err != nil {
			return errors.Wrapf(err, "the %s deployment never rolled out", deployment)
		}
	}
	// The triggerer is a StatefulSet at the chart default (persistence
	// on).
	if err := kubectlRolloutStatus(ctx, kubeconfig, "statefulset/"+v.Name+"-triggerer", v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the triggerer statefulset never rolled out")
	}
	if v.CeleryEnabled {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "statefulset/"+v.Name+"-worker", v.Namespace, 10*time.Minute); err != nil {
			return errors.Wrap(err, "the celery worker statefulset never rolled out")
		}
	}
	if v.BundledRedis {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "statefulset/"+v.Name+"-redis", v.Namespace, 10*time.Minute); err != nil {
			return errors.Wrap(err, "the bundled redis statefulset never rolled out")
		}
	}

	cancel, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name+"-api-server", v.Namespace, airflowApiPort+":8080")
	if err != nil {
		return errors.Wrap(err, "starting port-forward to the api server")
	}
	defer cancel()

	// THE AUTH GATE — before any credentialed call: the dags API must
	// reject anonymous reads (source-verified authenticated route).
	if err := v.proveAuthGate(ctx); err != nil {
		return err
	}

	token, err := v.login(ctx, kubeconfig)
	if err != nil {
		return err
	}

	// The API server's own health contract: metadatabase healthy (the
	// composed database wiring, end to end) and the scheduler
	// heartbeating.
	if err := v.proveHealth(ctx, token); err != nil {
		return err
	}

	if v.DeliveryProof {
		return v.proveDagDelivery(ctx, kubeconfig, token)
	}

	// THE DAG PROOF — the bundled example DAG (load_examples is
	// declared on the proof scenarios) runs to a successful DagRun.
	if err := v.proveDagRun(ctx, token, "example_bash_operator", "e2e-proof"); err != nil {
		return err
	}
	return nil
}

func (v *AirflowVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	for _, deployment := range []string{"api-server", "scheduler"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name+"-"+deployment, v.Namespace); err != nil {
			return err
		}
	}
	fmt.Printf("  [verify] DESTROY: the airflow deployments are gone (Airflow installs no CRDs — destroy is clean by design)\n")
	return nil
}

// adminPassword reads the module-generated admin credential from the
// exported `<name>-admin-auth` Secret (key password).
func (v *AirflowVerifier) adminPassword(ctx context.Context, kubeconfig string) (string, error) {
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
	// The chart's documented default must never ship.
	if password == "admin" {
		return "", errors.New("admin password is the chart's public default — the module-generated Secret was not wired")
	}
	return password, nil
}

// proveAuthGate asserts an anonymous read of the dags API is REJECTED
// (401 — every /api/v2 route except monitoring requires a token,
// source-verified at the pin).
func (v *AirflowVerifier) proveAuthGate(ctx context.Context) error {
	status, _, err := airflowApiRequest(ctx, "GET", "http://127.0.0.1:"+airflowApiPort+"/api/v2/dags", "", "", 3*time.Minute)
	if err != nil {
		return errors.Wrap(err, "the anonymous dags-API probe never answered")
	}
	if status != http.StatusUnauthorized && status != http.StatusForbidden {
		return errors.Errorf("THE AUTH GATE FAILED: anonymous GET /api/v2/dags answered %d, expected 401/403", status)
	}
	fmt.Printf("  [verify] AUTH GATE: anonymous API read rejected (%d)\n", status)
	return nil
}

// login exchanges the admin credential for a JWT at the api server's
// own token endpoint (POST /auth/token — source-verified; the FAB auth
// manager is the chart's default posture).
func (v *AirflowVerifier) login(ctx context.Context, kubeconfig string) (string, error) {
	password, err := v.adminPassword(ctx, kubeconfig)
	if err != nil {
		return "", err
	}
	payload, _ := json.Marshal(map[string]string{
		"username": v.AdminUsername,
		"password": password,
	})
	status, body, err := airflowApiRequest(ctx, "POST", "http://127.0.0.1:"+airflowApiPort+"/auth/token", "", string(payload), 3*time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "the token endpoint never answered")
	}
	if status < 200 || status >= 300 {
		return "", errors.Errorf("login as %q failed with %d: %s", v.AdminUsername, status, firstLines(body, 2))
	}
	var parsed struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil || parsed.AccessToken == "" {
		return "", errors.Errorf("the token endpoint answered without an access_token: %s", firstLines(body, 2))
	}
	fmt.Printf("  [verify] LOGIN: admin JWT issued from the module-generated credential\n")
	return parsed.AccessToken, nil
}

// proveHealth asserts the api server's health contract: metadatabase
// AND scheduler healthy (GET /api/v2/monitor/health).
func (v *AirflowVerifier) proveHealth(ctx context.Context, token string) error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastBody string
	for time.Now().Before(deadline) {
		status, body, err := airflowApiRequest(ctx, "GET", "http://127.0.0.1:"+airflowApiPort+"/api/v2/monitor/health", token, "", 30*time.Second)
		if err == nil && status == http.StatusOK {
			var parsed struct {
				Metadatabase struct {
					Status string `json:"status"`
				} `json:"metadatabase"`
				Scheduler struct {
					Status string `json:"status"`
				} `json:"scheduler"`
			}
			if json.Unmarshal([]byte(body), &parsed) == nil &&
				parsed.Metadatabase.Status == "healthy" && parsed.Scheduler.Status == "healthy" {
				fmt.Printf("  [verify] HEALTH: metadatabase healthy (the composed database answers) and the scheduler is heartbeating\n")
				return nil
			}
			lastBody = body
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("the health API never reported metadatabase+scheduler healthy: %s", firstLines(lastBody, 3))
}

// proveDagRun unpauses a DAG, triggers a run through the REST API, and
// polls it to state=success — THE DAG PROOF.
func (v *AirflowVerifier) proveDagRun(ctx context.Context, token, dagId, runPrefix string) error {
	base := "http://127.0.0.1:" + airflowApiPort + "/api/v2/dags/" + dagId

	// The DAG must be parsed before it can run — poll the dag-processor's
	// output through the API (examples parse at boot; a delivered DAG
	// parses within the processor's refresh loop).
	deadline := time.Now().Add(5 * time.Minute)
	for {
		status, _, err := airflowApiRequest(ctx, "GET", base, token, "", 30*time.Second)
		if err == nil && status == http.StatusOK {
			break
		}
		if time.Now().After(deadline) {
			return errors.Errorf("DAG %q never appeared in the API — the dag processor did not parse it", dagId)
		}
		time.Sleep(10 * time.Second)
	}

	// Unpause (DAGs are born paused unless configured otherwise).
	status, body, err := airflowApiRequest(ctx, "PATCH", base+"?update_mask=is_paused", token, `{"is_paused": false}`, 2*time.Minute)
	if err != nil || status < 200 || status >= 300 {
		return errors.Errorf("unpausing DAG %q failed (%d): %s", dagId, status, firstLines(body, 2))
	}

	runId := fmt.Sprintf("%s-%d", runPrefix, time.Now().Unix())
	payload := fmt.Sprintf(`{"dag_run_id": %q, "logical_date": null}`, runId)
	status, body, err = airflowApiRequest(ctx, "POST", base+"/dagRuns", token, payload, 2*time.Minute)
	if err != nil || status < 200 || status >= 300 {
		return errors.Errorf("triggering DAG %q failed (%d): %s", dagId, status, firstLines(body, 2))
	}
	fmt.Printf("  [verify] SUBMIT: DagRun %q triggered on DAG %q\n", runId, dagId)

	// Poll the run to success. Task pods (KubernetesExecutor) or worker
	// dispatch (Celery) happen inside this budget.
	pollDeadline := time.Now().Add(10 * time.Minute)
	var lastState string
	for time.Now().Before(pollDeadline) {
		status, body, err = airflowApiRequest(ctx, "GET", base+"/dagRuns/"+runId, token, "", 30*time.Second)
		if err == nil && status == http.StatusOK {
			var parsed struct {
				State string `json:"state"`
			}
			if json.Unmarshal([]byte(body), &parsed) == nil {
				lastState = parsed.State
				if parsed.State == "success" {
					fmt.Printf("  [verify] RUN: DagRun %q on %q reached state=success — the orchestrator runs pipelines\n", runId, dagId)
					return nil
				}
				if parsed.State == "failed" {
					return errors.Errorf("DagRun %q on %q FAILED — the engine cannot run its own DAG", runId, dagId)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("DagRun %q on %q never reached success (last state %q)", runId, dagId, lastState)
}

// airflowMarkerDag is the verifier-owned pipeline for THE DELIVERY
// PROOF — one bash task echoing a marker (the standard provider ships
// in the official image; the legacy import path covers older lines).
const airflowMarkerDag = `
try:
    from airflow.providers.standard.operators.bash import BashOperator
except ImportError:  # older provider layout
    from airflow.operators.bash import BashOperator
from airflow.sdk import DAG

with DAG(dag_id="e2e_delivery_proof", schedule=None, catchup=False) as dag:
    BashOperator(task_id="echo_marker", bash_command="echo e2e-delivery-ok")
`

// proveDagDelivery writes the marker DAG into the shared dags volume
// through the dag-processor pod (the parser's own mount), waits for it
// to parse, runs it to success, then replaces the SCHEDULER pod
// (UID-verified) and runs it again — scheduling state lives in the
// database, never in the pod.
func (v *AirflowVerifier) proveDagDelivery(ctx context.Context, kubeconfig, token string) error {
	// Write through kubectl exec — the dag-processor mounts the shared
	// dags volume read-write and parses what lands there.
	writeCmd := fmt.Sprintf("cat > /opt/airflow/dags/e2e_delivery_proof.py <<'PYEOF'\n%s\nPYEOF", airflowMarkerDag)
	execArgs := []string{
		"--kubeconfig", kubeconfig, "exec", "-n", v.Namespace,
		"deploy/" + v.Name + "-dag-processor", "--", "bash", "-c", writeCmd,
	}
	if out, err := exec.CommandContext(ctx, "kubectl", execArgs...).CombinedOutput(); err != nil {
		return errors.Wrapf(err, "writing the marker DAG into the shared dags volume: %s", string(out))
	}
	fmt.Printf("  [verify] DELIVERY: marker DAG written into the shared dags volume\n")

	if err := v.proveDagRun(ctx, token, "e2e_delivery_proof", "delivery-proof"); err != nil {
		return errors.Wrap(err, "the delivered DAG never ran to success")
	}

	// THE DURABILITY ARM: replace the scheduler, then a fresh run must
	// still schedule and succeed.
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
		"component=scheduler,release="+v.Name, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the scheduler pod did not recover after deletion")
	}
	if err := v.proveDagRun(ctx, token, "e2e_delivery_proof", "post-replacement"); err != nil {
		return errors.Wrap(err, "a DagRun triggered AFTER the scheduler replacement should still succeed")
	}
	fmt.Printf("  [verify] DURABILITY: the delivered DAG ran to success again after a UID-verified scheduler replacement — scheduling state lives in the database\n")
	return nil
}

// airflowApiRequest is one HTTP request with optional bearer token and
// JSON body, retried until the budget expires on transport errors only
// (HTTP error statuses return immediately — they are answers).
func airflowApiRequest(ctx context.Context, method, url, token, body string, budget time.Duration) (int, string, error) {
	deadline := time.Now().Add(budget)
	var lastErr error
	for {
		req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
		if err != nil {
			return 0, "", err
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			raw, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr == nil {
				return resp.StatusCode, string(raw), nil
			}
			lastErr = readErr
		} else {
			lastErr = err
		}
		if time.Now().After(deadline) {
			return 0, "", lastErr
		}
		time.Sleep(5 * time.Second)
	}
}

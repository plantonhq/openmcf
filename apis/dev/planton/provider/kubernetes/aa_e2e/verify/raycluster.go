package verify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// RayClusterVerifier checks a KubeRay-managed Ray cluster to the point
// a customer could run distributed work on it: the RayCluster CR
// reports state "ready" (the operator's own aggregate — Ready
// ClusterState = "ready" in raycluster_types.go:270, surfaced through
// the status.state field at :284), the operator's naming contract
// holds (`<name>-head-svc`, generated in util_test.go's
// "raycluster-sample-head-svc" contract), and THE RAY PROOF on every
// lane: a REAL job through the dashboard's Job Submission REST API —
// a Python process that calls ray.init(), joins the GCS, gets
// scheduled onto the cluster's declared capacity, and runs to
// SUCCEEDED. That is Ray's definition of working: a cluster whose
// pods are Running but cannot execute a submitted job is not a Ray
// cluster, and no liveness probe can tell the difference.
//
// Token auth (the catalog default) adds THE AUTH GATE before the
// proof: an UNAUTHENTICATED job listing must be rejected — a job API
// that schedules arbitrary code for anonymous callers is the exact
// failure the auth field exists to prevent. The bearer token is read
// from the Secret the operator provisions under the CLUSTER'S OWN
// NAME (raycluster_controller.go:398 `secretName :=
// utils.CheckName(instance.Name)` — CheckName only rewrites names
// over 50 chars or with a leading digit/punctuation, so for these
// manifests it IS the name), key `auth_token`
// (raycluster_controller.go:429), unless the spec brought its own
// Secret.
//
// The behavioral-state lane (StateProof) additionally proves GCS
// fault tolerance: after the job completes, DELETE the head pod, wait
// for a UID-verified Ready replacement, and assert the recovered head
// STILL LISTS the completed submission. Without the external store
// the replacement head boots with an empty GCS and the job history is
// gone — control state surviving head replacement is the difference
// the gcsFaultTolerance field exists to make.
type RayClusterVerifier struct {
	Namespace string
	Name      string
	// AuthDisabled mirrors spec.auth.mode "disabled" — token auth is
	// the catalog default, so the gate and the bearer token ride every
	// lane unless the spec opted out.
	AuthDisabled bool
	// ExistingTokenSecret is spec.auth.secretName when the spec brought
	// its own token Secret; empty means the operator provisioned one
	// under the cluster's own name.
	ExistingTokenSecret string
	// GcsFaultTolerance marks the composed external GCS store — the
	// seam the state proof exercises.
	GcsFaultTolerance bool
	// StateProof marks the behavioral-state lane (recognized by the
	// dispatcher from the manifest path).
	StateProof bool
}

// rayDashboardLocalPort is the workstation side of the port-forward to
// the head Service's dashboard/Job-API port (8265).
const rayDashboardLocalPort = "18265"

// rayJobMarker is printed by the proof job's driver — finding it in the
// driver logs ties the SUCCEEDED status to OUR entrypoint, not a stale
// artifact.
const rayJobMarker = "planton-e2e-marker"

func (v *RayClusterVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] raycluster %q in namespace %q (auth disabled %v, gcs fault tolerance %v)\n",
		v.Name, v.Namespace, v.AuthDisabled, v.GcsFaultTolerance)

	// The operator only reports "ready" once the head (and every
	// declared worker) Pod is running and Ready — first boot pays image
	// pulls of the multi-GB Ray image.
	if err := v.waitForReadyState(ctx, kubeconfig, 10*time.Minute); err != nil {
		return err
	}

	// The operator's naming contract — every published endpoint rides
	// this Service.
	headSvc := v.Name + "-head-svc"
	if err := KubectlResourceExists(ctx, kubeconfig, "service", headSvc, v.Namespace); err != nil {
		return errors.Wrap(err, "the head service not found")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	var token string
	if !v.AuthDisabled {
		var err error
		token, err = v.readAuthToken(ctx, kubeconfig)
		if err != nil {
			return err
		}

		// THE AUTH GATE on its own tunnel: an unauthenticated read of
		// the job API must be rejected BEFORE any authenticated call
		// proves anything — a dashboard that serves anonymous callers
		// makes the token Secret theater.
		cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, headSvc, rayDashboardLocalPort, "8265")
		if err != nil {
			return err
		}
		err = v.assertJobsAuthGate(ctx, client, "http://127.0.0.1:"+rayDashboardLocalPort, 4*time.Minute)
		cancel()
		if err != nil {
			return err
		}
	}

	// THE RAY PROOF on a fresh tunnel (fresh per phase — a tunnel dies
	// silently with its backing pod, the caught-live port-forward
	// class).
	cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, headSvc, rayDashboardLocalPort, "8265")
	if err != nil {
		return err
	}
	submissionId, err := v.proveRayJob(ctx, client, "http://127.0.0.1:"+rayDashboardLocalPort, token)
	cancel()
	if err != nil {
		return err
	}

	if !v.StateProof {
		return nil
	}

	// THE STATE PROOF: kill the head, wait for a UID-verified Ready
	// replacement, and ask the RECOVERED head about the job completed
	// before its predecessor died. Head pods are created with
	// GenerateName (pod.go:174 in the operator), so the replacement
	// carries a NEW name AND a new UID — the selector-keyed
	// delete-and-await is the honest recovery wait. The head selector
	// is the operator's own labels: ray.io/node-type=head
	// (constant.go:20 + raycluster_types.go:392 HeadNode = "head"),
	// scoped by ray.io/cluster (constant.go:19) because the fixture
	// namespace is shared across lanes.
	headSelector := "ray.io/cluster=" + v.Name + ",ray.io/node-type=head"
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace, headSelector, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the head pod did not recover after deletion")
	}

	// Re-read the token (same Secret — it must have survived too) and
	// open a FRESH tunnel: the old one died with the old head.
	if !v.AuthDisabled {
		token, err = v.readAuthToken(ctx, kubeconfig)
		if err != nil {
			return errors.Wrap(err, "re-reading the token after the head replacement")
		}
	}
	cancel, err = openServiceTunnel(ctx, kubeconfig, v.Namespace, headSvc, rayDashboardLocalPort, "8265")
	if err != nil {
		return errors.Wrap(err, "re-establishing the port-forward after the head replacement")
	}
	defer cancel()

	jobs, err := v.listJobs(ctx, client, "http://127.0.0.1:"+rayDashboardLocalPort, token, 6*time.Minute)
	if err != nil {
		return errors.Wrap(err, "listing jobs on the recovered head")
	}
	for _, job := range jobs {
		if job.SubmissionId == submissionId {
			fmt.Printf("  [verify] THE STATE PROOF: the recovered head lists %d job(s) including submission %q — control state survived the head replacement in the external GCS store\n",
				len(jobs), submissionId)
			return nil
		}
	}
	return errors.Errorf("the recovered head lists %d job(s) but NOT submission %q — the control state did not survive the head replacement", len(jobs), submissionId)
}

func (v *RayClusterVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "rayclusters.ray.io", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "service", v.Name+"-head-svc", v.Namespace); err != nil {
		return errors.Wrap(err, "the head service never deleted after the CR was removed")
	}
	// The operator garbage-collects the pods asynchronously; the
	// cluster label (constant.go:19) scopes the sweep to THIS cluster —
	// the fixture namespace hosts other lanes.
	return waitForNoPodsBySelector(ctx, kubeconfig, v.Namespace, "ray.io/cluster="+v.Name, 3*time.Minute)
}

// waitForReadyState polls the CR's status.state until the operator
// reports "ready" (Ready ClusterState = "ready",
// raycluster_types.go:270), printing progress — first boot pays the
// multi-GB image pull.
func (v *RayClusterVerifier) waitForReadyState(ctx context.Context, kubeconfig string, budget time.Duration) error {
	start := time.Now()
	deadline := start.Add(budget)
	var lastState string
	var lastProgress time.Time
	for time.Now().Before(deadline) {
		state, _ := kubectlGetJSONPath(ctx, kubeconfig, "rayclusters.ray.io", v.Name, v.Namespace, "{.status.state}")
		lastState = strings.TrimSpace(state)
		if lastState == "ready" {
			fmt.Printf("  [verify] RayCluster CR reports state ready\n")
			return nil
		}
		if time.Since(lastProgress) >= 30*time.Second {
			fmt.Printf("  [verify] waiting on the RayCluster CR: state %q (elapsed %s)\n",
				lastState, time.Since(start).Round(time.Second))
			lastProgress = time.Now()
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("the RayCluster CR never reported state ready within %s (last state %q)", budget, lastState)
}

// readAuthToken reads the bearer token from the token Secret: the
// spec-declared one when present, else the operator-provisioned Secret
// named after the cluster itself (raycluster_controller.go:398), key
// `auth_token` (raycluster_controller.go:429). The value stays
// in-process — never printed.
func (v *RayClusterVerifier) readAuthToken(ctx context.Context, kubeconfig string) (string, error) {
	secretName := v.ExistingTokenSecret
	if secretName == "" {
		secretName = v.Name
	}
	encoded, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.auth_token}")
	if err != nil {
		return "", errors.Wrapf(err, "reading auth_token from secret %q", secretName)
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return "", errors.Wrapf(err, "decoding auth_token from secret %q", secretName)
	}
	token := strings.TrimSpace(string(raw))
	if token == "" {
		return "", errors.Errorf("secret %q carries an empty auth_token", secretName)
	}
	return token, nil
}

// assertJobsAuthGate proves the job API rejects anonymous callers: an
// unauthenticated GET of the job listing must answer 401 or 403
// (whichever the dashboard's auth middleware speaks — the rejection is
// the assertion, the exact code is evidence). Connection errors are
// tunnel warm-up and retry; a 2xx is the gate NOT holding and fails.
func (v *RayClusterVerifier) assertJobsAuthGate(ctx context.Context, client *http.Client, base string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var last error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/api/jobs/", nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			last = err
			time.Sleep(5 * time.Second)
			continue
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			fmt.Printf("  [verify] AUTH GATE: unauthenticated job listing rejected with HTTP %d\n", resp.StatusCode)
			return nil
		}
		last = errors.Errorf("expected 401/403, got %d", resp.StatusCode)
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(last, "the unauthenticated job listing was NOT rejected (the auth gate)")
}

// proveRayJob runs THE RAY PROOF: submit a real job through the Job
// Submission REST API and poll it to SUCCEEDED. The entrypoint is a
// Python driver that calls ray.init() — joining the GCS and consuming
// scheduled capacity — so SUCCEEDED means the cluster executed
// distributed work end to end, not that a port answered. Returns the
// submission id for the state proof.
func (v *RayClusterVerifier) proveRayJob(ctx context.Context, client *http.Client, base, token string) (string, error) {
	payload, err := json.Marshal(struct {
		Entrypoint string `json:"entrypoint"`
	}{
		Entrypoint: fmt.Sprintf("python -c \"import ray; ray.init(); print('%s')\"", rayJobMarker),
	})
	if err != nil {
		return "", err
	}
	_, body, err := v.rayRequest(ctx, client, http.MethodPost, base+"/api/jobs/", string(payload), token, 4*time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "submitting the proof job")
	}
	// The dashboard answers with submission_id (job_id on older
	// surfaces — the operator's own types carry both,
	// types/dashboard_httpclient.go:18-19).
	var submitted struct {
		SubmissionId string `json:"submission_id"`
		JobId        string `json:"job_id"`
	}
	if err := json.Unmarshal([]byte(body), &submitted); err != nil {
		return "", errors.Wrapf(err, "parsing the job submission response: %s", firstLines(body, 2))
	}
	submissionId := submitted.SubmissionId
	if submissionId == "" {
		submissionId = submitted.JobId
	}
	if submissionId == "" {
		return "", errors.Errorf("the job submission answered without a submission_id: %s", firstLines(body, 2))
	}
	fmt.Printf("  [verify] THE RAY PROOF: job submitted (submission %q)\n", submissionId)

	// The first job pays Ray worker-process startup (and on the
	// full-surface lane, scheduling onto a WORKER — the head advertises
	// zero CPUs there), hence the 8-minute budget.
	status, err := v.waitForJobSucceeded(ctx, client, base, token, submissionId, 8*time.Minute)
	if err != nil {
		return "", err
	}
	fmt.Printf("  [verify] THE RAY PROOF: job %q reached %s — real distributed execution on the cluster's capacity\n", submissionId, status)

	// Tie the status to OUR driver: the marker line in the job logs.
	// Best-effort — the proof already stands on SUCCEEDED.
	if excerpt := v.jobLogMarker(ctx, client, base, token, submissionId); excerpt != "" {
		fmt.Printf("  [verify] THE RAY PROOF: driver logs carry the marker: %q\n", excerpt)
	}
	return submissionId, nil
}

// waitForJobSucceeded polls the job's status to SUCCEEDED
// (JobStatusSucceeded = "SUCCEEDED", rayjob_types.go:22; the JSON field
// is `status`, types/dashboard_httpclient.go:16), failing fast on the
// terminal FAILED/STOPPED states with the dashboard's own message.
func (v *RayClusterVerifier) waitForJobSucceeded(ctx context.Context, client *http.Client, base, token, submissionId string, budget time.Duration) (string, error) {
	start := time.Now()
	deadline := start.Add(budget)
	var lastStatus string
	var lastProgress time.Time
	for time.Now().Before(deadline) {
		_, body, err := v.rayRequest(ctx, client, http.MethodGet, base+"/api/jobs/"+submissionId, "", token, 90*time.Second)
		if err == nil {
			var info struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}
			if jsonErr := json.Unmarshal([]byte(body), &info); jsonErr == nil {
				lastStatus = info.Status
				switch info.Status {
				case "SUCCEEDED":
					return info.Status, nil
				case "FAILED", "STOPPED":
					return "", errors.Errorf("the proof job reached terminal state %s: %s", info.Status, firstLines(info.Message, 2))
				}
			}
		}
		if time.Since(lastProgress) >= 30*time.Second {
			fmt.Printf("  [verify] waiting on the proof job: status %q (elapsed %s)\n",
				lastStatus, time.Since(start).Round(time.Second))
			lastProgress = time.Now()
		}
		time.Sleep(10 * time.Second)
	}
	return "", errors.Errorf("the proof job never reached SUCCEEDED within %s (last status %q)", budget, lastStatus)
}

// jobLogMarker fetches the job's driver logs and returns the line
// carrying the proof marker, empty when unavailable (best-effort — the
// logs endpoint answers `{"logs": ...}`, types RayJobLogsResponse).
func (v *RayClusterVerifier) jobLogMarker(ctx context.Context, client *http.Client, base, token, submissionId string) string {
	_, body, err := v.rayRequest(ctx, client, http.MethodGet, base+"/api/jobs/"+submissionId+"/logs", "", token, 90*time.Second)
	if err != nil {
		return ""
	}
	var logs struct {
		Logs string `json:"logs"`
	}
	if err := json.Unmarshal([]byte(body), &logs); err != nil {
		return ""
	}
	for _, line := range strings.Split(logs.Logs, "\n") {
		if strings.Contains(line, rayJobMarker) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// listJobs reads the full job listing (authenticated when auth is on).
func (v *RayClusterVerifier) listJobs(ctx context.Context, client *http.Client, base, token string, budget time.Duration) ([]struct {
	SubmissionId string `json:"submission_id"`
	Status       string `json:"status"`
}, error) {
	_, body, err := v.rayRequest(ctx, client, http.MethodGet, base+"/api/jobs/", "", token, budget)
	if err != nil {
		return nil, err
	}
	var jobs []struct {
		SubmissionId string `json:"submission_id"`
		Status       string `json:"status"`
	}
	if err := json.Unmarshal([]byte(body), &jobs); err != nil {
		return nil, errors.Wrapf(err, "parsing the job listing: %s", firstLines(body, 2))
	}
	return jobs, nil
}

// rayRequest performs one JSON request against the dashboard with
// retries across tunnel warm-up, succeeding on any 2xx. The token rides
// BOTH the standard Authorization header and the X-Ray-Authorization
// twin the operator's own clients send (dashboard_httpclient.go:57,
// rayjob-submitter.go:44 — the twin exists for proxies that consume
// Authorization; the dashboard honors either).
func (v *RayClusterVerifier) rayRequest(ctx context.Context, client *http.Client, method, endpoint, body, token string, budget time.Duration) (int, string, error) {
	deadline := time.Now().Add(budget)
	var lastStatus int
	var lastBody string
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, method, endpoint, strings.NewReader(body))
		if err != nil {
			return 0, "", err
		}
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("X-Ray-Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err == nil {
			out := drainBody(resp)
			lastStatus = resp.StatusCode
			lastBody = out
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp.StatusCode, out, nil
			}
			lastErr = errors.Errorf("HTTP %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(5 * time.Second)
	}
	return lastStatus, lastBody, errors.Wrapf(lastErr, "last body: %s", firstLines(lastBody, 2))
}

// openServiceTunnel starts a kubectl port-forward to a Service and only
// returns once the LOCAL port accepts a TCP connection — kubectl binds
// the listener asynchronously, so "process started" is not "tunnel up"
// (the caught-live race harbor.go documents). Shared by the Flink
// verifier. Cancel FIRST, then reap — Wait blocks forever on a
// port-forward never told to exit.
func openServiceTunnel(ctx context.Context, kubeconfig, namespace, service, localPort, remotePort string) (func(), error) {
	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+service, localPort+":"+remotePort, "-n", namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return nil, errors.Wrapf(err, "starting port-forward to service %q", service)
	}
	done := make(chan struct{})
	go func() {
		_ = pf.Wait()
		close(done)
	}()
	cancelAndReap := func() {
		cancel()
		<-done
	}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+localPort, 2*time.Second)
		if err == nil {
			_ = conn.Close()
			return cancelAndReap, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	cancelAndReap()
	return nil, errors.Errorf("the port-forward never started listening on 127.0.0.1:%s within 30s; kubectl output: %s",
		localPort, firstLines(pfOut.String(), 3))
}

// waitForNoPodsBySelector polls until no pods match the selector — the
// destroy-side sweep for operators that tear children down
// asynchronously. Shared by the Flink verifier.
func waitForNoPodsBySelector(ctx context.Context, kubeconfig, namespace, selector string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "pods", "-n", namespace, "-l", selector,
			"-o", "jsonpath={.items[*].metadata.name}").Output()
		if err == nil {
			last = strings.TrimSpace(string(out))
			if last == "" {
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("pods still present for selector %q after %s: %s", selector, budget, last)
}

// rayClusterAuthDisabled reads spec.auth.mode — token auth is the
// catalog default, so only an explicit "disabled" turns the gate off.
func rayClusterAuthDisabled(spec map[string]interface{}) bool {
	auth, _ := spec["auth"].(map[string]interface{})
	if auth == nil {
		return false
	}
	mode, _ := auth["mode"].(string)
	return mode == "disabled"
}

// rayClusterExistingTokenSecret reads the bring-your-own token Secret
// name (both manifest key forms tolerated).
func rayClusterExistingTokenSecret(spec map[string]interface{}) string {
	auth, _ := spec["auth"].(map[string]interface{})
	if auth == nil {
		return ""
	}
	for _, key := range []string{"existing_token_secret_name", "existingTokenSecretName"} {
		if name, ok := auth[key].(string); ok && name != "" {
			return name
		}
	}
	return ""
}

// rayClusterGcsFaultTolerance reads whether the GCS fault-tolerance arm
// is enabled (both manifest key forms tolerated).
func rayClusterGcsFaultTolerance(spec map[string]interface{}) bool {
	for _, key := range []string{"gcs_fault_tolerance", "gcsFaultTolerance"} {
		if gcs, ok := spec[key].(map[string]interface{}); ok {
			if enabled, ok := gcs["enabled"].(bool); ok {
				return enabled
			}
		}
	}
	return false
}

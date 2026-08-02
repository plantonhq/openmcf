package verify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// LocustVerifier checks a Locust cluster to the point a customer could
// run a load test with it: master and worker rollouts, THE AUTH GATE
// (an anonymous stats read bounced to the login — upstream ships the
// web UI OPEN, able to fire load at any reachable host, and that
// posture never deploys from this kind), THE LOGIN PROOF (a wrong
// password refused; the module-generated credential signs in through
// the platform-managed login backend), and THE SWARM PROOF on every
// lane — a real distributed test started through the master's own REST
// API drives real requests through real workers at the composed
// fixture target with zero failures (a load tester that cannot swarm
// is not a load tester).
//
// The behavioral scenario (recognized by name) adds THE RECONNECT
// PROOF: the master pod is deleted (UID-verified replacement), the
// PRE-REPLACEMENT session cookie still authenticates (the Flask
// session-signing key is module-generated and STABLE by design — a
// per-start random would log every user out on every roll), the
// workers re-register with the replacement through the stable Service,
// and a SECOND swarm drives requests again — coordination recovers by
// design, not by accident.
//
// Destroy is clean by design: a plain Helm release plus the
// module-owned ConfigMaps and Secret — everything leaves with the
// resource.
type LocustVerifier struct {
	Namespace string
	Name      string
	// WebLoginEnabled gates the credentialed arms (headless runs and
	// explicitly disabled logins ship no login machinery).
	WebLoginEnabled bool
	// Username signs into the login backend (spec default "locust").
	Username string
	// WorkerReplicas gates the worker-rollout wait and the
	// worker-registration assertions.
	WorkerReplicas int
	// ReconnectProof switches the behavioral arm on.
	ReconnectProof bool
}

// locustWebLoginEnabled reads the secured default: web_ui_auth.enabled
// absent = true, and a headless run never starts the web UI at all.
func locustWebLoginEnabled(spec map[string]interface{}) bool {
	if loadTest, ok := spec["load_test"].(map[string]interface{}); ok {
		for _, key := range []string{"headless"} {
			if headless, ok := loadTest[key].(bool); ok && headless {
				return false
			}
		}
	}
	if auth, ok := spec["web_ui_auth"].(map[string]interface{}); ok {
		if enabled, ok := auth["enabled"].(bool); ok {
			return enabled
		}
	}
	return true
}

// locustUsername reads spec.web_ui_auth.username ("" = "locust").
func locustUsername(spec map[string]interface{}) string {
	if auth, ok := spec["web_ui_auth"].(map[string]interface{}); ok {
		if username, ok := auth["username"].(string); ok && username != "" {
			return username
		}
	}
	return "locust"
}

// locustWorkerReplicas reads spec.workers.replicas (absent = 1 — the
// chart default the module keeps).
func locustWorkerReplicas(spec map[string]interface{}) int {
	if workers, ok := spec["workers"].(map[string]interface{}); ok {
		if replicas, ok := workers["replicas"].(float64); ok {
			return int(replicas)
		}
	}
	return 1
}

const locustApiPort = "18089"

func locustBaseUrl() string {
	return "http://127.0.0.1:" + locustApiPort
}

func (v *LocustVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] locust %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-master", v.Namespace, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the locust master never rolled out")
	}
	if v.WorkerReplicas > 0 {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-worker", v.Namespace, 8*time.Minute); err != nil {
			return errors.Wrap(err, "the locust workers never rolled out")
		}
	}

	// The master Service is named exactly the fullname (chart
	// master-service.yaml at the pin): only the DEPLOYMENTS carry the
	// -master/-worker suffixes.
	cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name, locustApiPort, "8089")
	if err != nil {
		return errors.Wrap(err, "opening the tunnel to the locust master")
	}

	// One client per session: the login rides a signed session cookie.
	client, err := locustHTTPClient()
	if err != nil {
		cancel()
		return err
	}

	if v.WebLoginEnabled {
		// THE AUTH GATE before any credentialed call.
		if err := v.proveAuthGate(ctx, client); err != nil {
			cancel()
			return err
		}
		password, err := v.loginPassword(ctx, kubeconfig)
		if err != nil {
			cancel()
			return err
		}
		// A wrong password must be REFUSED before the real one is
		// accepted (a backend that lets anything in is worse than
		// none — it looks secure).
		if err := v.proveLogin(ctx, client, v.Username, "wrong-"+password[:8], false); err != nil {
			cancel()
			return err
		}
		if err := v.proveLogin(ctx, client, v.Username, password, true); err != nil {
			cancel()
			return err
		}
		fmt.Printf("  [verify] THE LOGIN PROOF: wrong password refused; the generated credential signed in\n")
	}

	// THE SWARM PROOF: a real distributed test through the master's
	// own REST API — workers registered, requests flowing, zero
	// failures against the composed fixture target.
	if err := v.proveSwarm(ctx, client, "the swarm proof"); err != nil {
		cancel()
		return err
	}

	if !v.ReconnectProof {
		cancel()
		return nil
	}

	// THE RECONNECT PROOF: drop the tunnel across the replacement
	// window (fresh-tunnel-per-phase), replace the master pod
	// UID-verified, then reuse the SAME session against the
	// replacement — the stable session-signing key keeps operators
	// signed in across restarts — and swarm again once the workers
	// re-register through the stable Service.
	cancel()
	fmt.Printf("  [verify] THE RECONNECT PROOF: replacing the master pod…\n")
	selector := "app.kubernetes.io/instance=" + v.Name + ",component=master"
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace, selector, 6*time.Minute); err != nil {
		return errors.Wrap(err, "the master pod was never replaced")
	}
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-master", v.Namespace, 6*time.Minute); err != nil {
		return errors.Wrap(err, "the replaced master never became ready")
	}
	cancel, err = openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name, locustApiPort, "8089")
	if err != nil {
		return errors.Wrap(err, "re-opening the tunnel after the master replacement")
	}
	defer cancel()
	if err := v.proveSwarm(ctx, client, "the post-replacement swarm proof"); err != nil {
		return err
	}
	fmt.Printf("  [verify] THE RECONNECT PROOF: the same session swarmed again through the replaced master\n")
	return nil
}

// VerifyAbsent asserts destroy left nothing behind: a plain Helm
// release plus module-owned ConfigMaps and Secret — no CRDs, no
// designed survivors.
func (v *LocustVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify-absent] locust %q in namespace %q\n", v.Name, v.Namespace)
	for _, deployment := range []string{v.Name + "-master", v.Name + "-worker"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", deployment, v.Namespace); err != nil {
			return err
		}
	}
	for _, configMap := range []string{v.Name + "-locustfile", v.Name + "-lib", v.Name + "-web-auth"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "configmap", configMap, v.Namespace); err != nil {
			return err
		}
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "secret", v.Name+"-auth", v.Namespace); err != nil {
		return err
	}
	return nil
}

// loginPassword reads the module-generated credential from the
// exported `<name>-auth` Secret (key password).
func (v *LocustVerifier) loginPassword(ctx context.Context, kubeconfig string) (string, error) {
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
		return "", errors.New("web-ui password Secret was empty")
	}
	return password, nil
}

// proveAuthGate reads the stats API WITHOUT a session and requires the
// login bounce — the open UI (any caller starts load tests) must be
// dead.
func (v *LocustVerifier) proveAuthGate(ctx context.Context, client *http.Client) error {
	status, _, location, err := locustRequest(ctx, client, http.MethodGet, "/stats/requests", nil)
	if err != nil {
		return errors.Wrap(err, "the anonymous stats read")
	}
	if status != http.StatusFound || !strings.Contains(location, "/login") {
		return errors.Errorf("THE AUTH GATE FAILED: anonymous stats read returned %d (location %q), want the login bounce — the open web UI must never ship", status, location)
	}
	fmt.Printf("  [verify] THE AUTH GATE: anonymous stats read bounced to the login\n")
	return nil
}

// proveLogin submits the login form (the platform-managed backend's
// own route) and asserts acceptance or refusal by where the redirect
// lands: the index on success, back to /login on refusal.
func (v *LocustVerifier) proveLogin(ctx context.Context, client *http.Client, username, password string, wantSuccess bool) error {
	form := url.Values{"username": {username}, "password": {password}}
	status, _, location, err := locustRequest(ctx, client, http.MethodPost, "/planton/login", form)
	if err != nil {
		return errors.Wrap(err, "the login submission")
	}
	if status != http.StatusFound {
		return errors.Errorf("the login submission returned %d, want a redirect", status)
	}
	succeeded := !strings.Contains(location, "/login")
	if succeeded != wantSuccess {
		if wantSuccess {
			return errors.Errorf("THE LOGIN PROOF FAILED: the generated credential was refused (redirected to %q)", location)
		}
		return errors.Errorf("THE LOGIN PROOF FAILED: a WRONG password was accepted (redirected to %q) — the backend authenticates nothing", location)
	}
	if wantSuccess {
		// The session must now read the gated API.
		statsStatus, _, _, err := locustRequest(ctx, client, http.MethodGet, "/stats/requests", nil)
		if err != nil {
			return errors.Wrap(err, "the authenticated stats read")
		}
		if statsStatus != http.StatusOK {
			return errors.Errorf("the authenticated stats read returned %d, want 200", statsStatus)
		}
	}
	return nil
}

// locustStats is the /stats/requests contract the proofs read (web.py
// at the pin).
type locustStats struct {
	State     string `json:"state"`
	UserCount int    `json:"user_count"`
	Workers   []struct {
		Id string `json:"id"`
	} `json:"workers"`
	Stats []struct {
		Name        string  `json:"name"`
		NumRequests float64 `json:"num_requests"`
		NumFailures float64 `json:"num_failures"`
	} `json:"stats"`
}

// proveSwarm starts a distributed test through the REST API, polls the
// stats until real requests flowed through the registered workers with
// zero failures, then stops the test.
func (v *LocustVerifier) proveSwarm(ctx context.Context, client *http.Client, label string) error {
	// The workers register with the master asynchronously — wait for
	// the full fleet before starting (a swarm with zero workers
	// "runs" but generates nothing).
	if err := v.awaitWorkers(ctx, client); err != nil {
		return errors.Wrap(err, label)
	}

	form := url.Values{"user_count": {"10"}, "spawn_rate": {"5"}}
	status, body, _, err := locustRequest(ctx, client, http.MethodPost, "/swarm", form)
	if err != nil {
		return errors.Wrapf(err, "%s: starting the swarm", label)
	}
	if status != http.StatusOK {
		return errors.Errorf("%s: starting the swarm returned %d: %s", label, status, locustTruncate(body))
	}
	var swarmResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &swarmResponse); err != nil {
		return errors.Wrapf(err, "%s: decoding the swarm response", label)
	}
	if !swarmResponse.Success {
		return errors.Errorf("%s: the swarm refused to start: %s", label, swarmResponse.Message)
	}

	// Poll for real traffic: aggregated requests over the threshold
	// with ZERO failures (a failing swarm would prove connectivity
	// problems, a missing env var, or a broken script — not load).
	deadline := time.Now().Add(3 * time.Minute)
	for {
		if time.Now().After(deadline) {
			return errors.Errorf("%s: the swarm never generated the expected traffic", label)
		}
		time.Sleep(5 * time.Second)

		stats, err := v.readStats(ctx, client)
		if err != nil {
			return errors.Wrapf(err, "%s: reading the stats", label)
		}
		var requests, failures float64
		for _, entry := range stats.Stats {
			if entry.Name == "Aggregated" {
				requests, failures = entry.NumRequests, entry.NumFailures
			}
		}
		if failures > 0 {
			return errors.Errorf("%s: the swarm recorded %v failed requests — the target (or the test environment wiring) is broken", label, failures)
		}
		if requests >= 20 && stats.UserCount > 0 {
			fmt.Printf("  [verify] THE SWARM PROOF (%s): %v requests, 0 failures, %d users across %d workers\n",
				label, requests, stats.UserCount, len(stats.Workers))
			break
		}
	}

	// Stop the test — the lane leaves no swarm running behind it.
	status, body, _, err = locustRequest(ctx, client, http.MethodGet, "/stop", nil)
	if err != nil {
		return errors.Wrapf(err, "%s: stopping the swarm", label)
	}
	if status != http.StatusOK {
		return errors.Errorf("%s: stopping the swarm returned %d: %s", label, status, locustTruncate(body))
	}
	return nil
}

// awaitWorkers polls the stats until every declared worker registered.
func (v *LocustVerifier) awaitWorkers(ctx context.Context, client *http.Client) error {
	deadline := time.Now().Add(4 * time.Minute)
	for {
		stats, err := v.readStats(ctx, client)
		if err == nil && len(stats.Workers) >= v.WorkerReplicas {
			fmt.Printf("  [verify] %d worker(s) registered with the master\n", len(stats.Workers))
			return nil
		}
		if time.Now().After(deadline) {
			registered := 0
			if err == nil {
				registered = len(stats.Workers)
			}
			return errors.Errorf("only %d of %d workers registered with the master", registered, v.WorkerReplicas)
		}
		time.Sleep(5 * time.Second)
	}
}

func (v *LocustVerifier) readStats(ctx context.Context, client *http.Client) (*locustStats, error) {
	status, body, _, err := locustRequest(ctx, client, http.MethodGet, "/stats/requests", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("the stats read returned %d: %s", status, locustTruncate(body))
	}
	var stats locustStats
	if err := json.Unmarshal([]byte(body), &stats); err != nil {
		return nil, errors.Wrap(err, "decoding the stats")
	}
	return &stats, nil
}

// locustHTTPClient builds the session-carrying client: a cookie jar
// (the login is a signed session cookie) with redirects NOT followed —
// the proofs assert on the redirect targets themselves.
func locustHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating the cookie jar")
	}
	return &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

// locustRequest performs one call against the tunneled master and
// returns the status, body and any redirect Location.
func locustRequest(ctx context.Context, client *http.Client, method, path string, form url.Values) (int, string, string, error) {
	var payload io.Reader
	if form != nil {
		payload = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, locustBaseUrl()+path, payload)
	if err != nil {
		return 0, "", "", err
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", "", err
	}
	return resp.StatusCode, string(responseBody), resp.Header.Get("Location"), nil
}

// locustTruncate keeps error payloads readable.
func locustTruncate(body string) string {
	if len(body) > 300 {
		return body[:300] + "…"
	}
	return body
}

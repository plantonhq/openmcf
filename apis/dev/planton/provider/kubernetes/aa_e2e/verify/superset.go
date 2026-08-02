package verify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// SupersetVerifier checks a Superset deployment to the point a customer
// could hand the URL to their analysts: the web rollout (and the worker
// rollout when the Celery arm is on), the server's own /health contract,
// THE AUTH GATE (an anonymous API read REJECTED and the chart's
// documented admin/admin default asserted DEAD — the module-generated
// credential is the only way in), a REAL login through Superset's own
// security API, and THE DASHBOARD PROOF on every lane — a dashboard is
// created through the REST API (JWT + CSRF, the same contract real
// clients use) and read back (a BI platform that cannot hold a
// dashboard is not a BI platform). Proof artifacts are swept.
//
// The behavioral-durability scenario (recognized by name) adds THE
// STATE PROOF: the web pod is deleted (UID-verified replacement) and a
// fresh session must sign in and find the same dashboard — BI state
// lives in the composed PostgreSQL, never in the pod.
//
// Destroy is clean by design: a plain Helm release (no CRDs) plus the
// module-owned Secrets — everything leaves with the resource.
type SupersetVerifier struct {
	Namespace string
	Name      string
	// AdminUsername logs into the security API (spec default "admin").
	AdminUsername string
	// WorkerEnabled gates the worker-rollout wait (web-only shapes have
	// no worker Deployment).
	WorkerEnabled bool
	// StateProof switches the behavioral arm on.
	StateProof bool
}

// supersetAdminUsername reads spec.init.admin.username ("" = "admin").
func supersetAdminUsername(spec map[string]interface{}) string {
	if init, ok := spec["init"].(map[string]interface{}); ok {
		if admin, ok := init["admin"].(map[string]interface{}); ok {
			if username, ok := admin["username"].(string); ok && username != "" {
				return username
			}
		}
	}
	return "admin"
}

// supersetWorkerEnabled mirrors the module's resolution: workers run
// when the cache is declared and the worker is not explicitly disabled.
func supersetWorkerEnabled(spec map[string]interface{}) bool {
	if _, ok := spec["cache"].(map[string]interface{}); !ok {
		return false
	}
	if worker, ok := spec["worker"].(map[string]interface{}); ok {
		if enabled, ok := worker["enabled"].(bool); ok {
			return enabled
		}
	}
	return true
}

const supersetApiPort = "18084"

func supersetBaseUrl() string {
	return "http://127.0.0.1:" + supersetApiPort
}

func (v *SupersetVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] superset %q in namespace %q\n", v.Name, v.Namespace)

	// The web rollout gates on the init Job's schema migration against
	// the composed database (readiness needs a migrated schema).
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 12*time.Minute); err != nil {
		return errors.Wrap(err, "the superset web deployment never rolled out")
	}
	if v.WorkerEnabled {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-worker", v.Namespace, 10*time.Minute); err != nil {
			return errors.Wrap(err, "the superset worker deployment never rolled out")
		}
	}

	cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name, supersetApiPort, "8088")
	if err != nil {
		return errors.Wrap(err, "opening the tunnel to the superset service")
	}

	if err := v.proveHealth(ctx); err != nil {
		cancel()
		return err
	}

	// THE AUTH GATE before any credentialed call.
	if err := v.proveAuthGate(ctx); err != nil {
		cancel()
		return err
	}

	password, err := v.adminPassword(ctx, kubeconfig)
	if err != nil {
		cancel()
		return err
	}

	// THE DASHBOARD PROOF.
	dashboardTitle := "e2e-proof-" + v.Name
	session, dashboardId, err := v.proveDashboard(ctx, password, dashboardTitle)
	if err != nil {
		cancel()
		return err
	}
	fmt.Printf("  [verify] THE DASHBOARD PROOF: dashboard %q created (id %s) and read back\n", dashboardTitle, dashboardId)

	if !v.StateProof {
		// Sweep the proof artifact and finish.
		if err := session.deleteDashboard(ctx, dashboardId); err != nil {
			cancel()
			return err
		}
		cancel()
		return nil
	}

	// THE STATE PROOF: drop the tunnel across the replacement window
	// (fresh-tunnel-per-phase), replace the web pod UID-verified, then
	// sign in fresh and find the same dashboard — BI state lives in
	// the composed database.
	cancel()
	fmt.Printf("  [verify] THE STATE PROOF: replacing the web pod…\n")
	// The chart stamps every component's pods with the standard
	// app.kubernetes.io labels plus a per-component
	// `app.kubernetes.io/component` (web/worker/…) — VERIFIED LIVE on
	// the running pods at the pin. (An earlier chart line stamped
	// legacy app/release labels instead — read the label shape from
	// the clone checked out at the CHART tag, not the app tag.)
	selector := "app.kubernetes.io/instance=" + v.Name + ",app.kubernetes.io/component=web"
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace, selector, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the web pod was never replaced")
	}
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the replaced web pod never became ready")
	}
	cancel, err = openServiceTunnel(ctx, kubeconfig, v.Namespace, v.Name, supersetApiPort, "8088")
	if err != nil {
		return errors.Wrap(err, "re-opening the tunnel after the web replacement")
	}
	defer cancel()

	freshSession, err := v.login(ctx, password)
	if err != nil {
		return errors.Wrap(err, "the post-replacement sign-in")
	}
	title, err := freshSession.dashboardTitle(ctx, dashboardId)
	if err != nil {
		return errors.Wrap(err, "re-reading the dashboard after the replacement")
	}
	if title != dashboardTitle {
		return errors.Errorf("THE STATE PROOF FAILED: dashboard %s carries title %q after the replacement, want %q", dashboardId, title, dashboardTitle)
	}
	fmt.Printf("  [verify] THE STATE PROOF: dashboard %q survived the web replacement\n", dashboardTitle)
	return freshSession.deleteDashboard(ctx, dashboardId)
}

// VerifyAbsent asserts destroy left nothing behind: a plain Helm
// release plus module-owned Secrets — no CRDs, no designed survivors.
func (v *SupersetVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify-absent] superset %q in namespace %q\n", v.Name, v.Namespace)
	for _, deployment := range []string{v.Name, v.Name + "-worker"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", deployment, v.Namespace); err != nil {
			return err
		}
	}
	for _, secret := range []string{v.Name + "-env", v.Name + "-secret-key", v.Name + "-admin-auth"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "secret", secret, v.Namespace); err != nil {
			return err
		}
	}
	return nil
}

// adminPassword reads the module-generated admin credential from the
// exported `<name>-admin-auth` Secret (key password).
func (v *SupersetVerifier) adminPassword(ctx context.Context, kubeconfig string) (string, error) {
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
	if password == "admin" {
		return "", errors.New("admin password is the chart's documented default — the module-generated Secret was not wired")
	}
	return password, nil
}

// proveHealth polls the server's own health contract (it answers
// unauthenticated) while gunicorn warms up.
func (v *SupersetVerifier) proveHealth(ctx context.Context) error {
	var lastErr error
	for attempt := 0; attempt < 20; attempt++ {
		if attempt > 0 {
			time.Sleep(6 * time.Second)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, supersetBaseUrl()+"/health", nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK && strings.Contains(strings.ToUpper(string(body)), "OK") {
			fmt.Printf("  [verify] /health answered OK\n")
			return nil
		}
		lastErr = errors.Errorf("/health returned %d: %s", resp.StatusCode, string(body))
	}
	return errors.Wrap(lastErr, "the health contract never answered")
}

// proveAuthGate requires an anonymous API read to be rejected AND the
// chart's documented admin/admin default to be dead on the wire.
func (v *SupersetVerifier) proveAuthGate(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, supersetBaseUrl()+"/api/v1/me/", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "the anonymous API read")
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		return errors.Errorf("THE AUTH GATE FAILED: anonymous /api/v1/me/ returned %d, want 401", resp.StatusCode)
	}
	fmt.Printf("  [verify] THE AUTH GATE: anonymous API read rejected (401)\n")

	// The chart's documented default credential must be dead — unless
	// the operator legitimately chose "admin" as the username, in
	// which case only the PASSWORD must differ (checked at read time).
	status, _, err := supersetJsonRequest(ctx, http.DefaultClient, http.MethodPost, "/api/v1/security/login",
		map[string]interface{}{"username": "admin", "password": "admin", "provider": "db", "refresh": true}, nil)
	if err != nil {
		return errors.Wrap(err, "probing the default credential")
	}
	if status == http.StatusOK {
		return errors.New("THE AUTH GATE FAILED: the chart's documented admin/admin default signs in — the generated credential was not wired")
	}
	fmt.Printf("  [verify] THE AUTH GATE: the documented admin/admin default is dead (%d)\n", status)
	return nil
}

// supersetSession carries the signed-in API context: the JWT, the CSRF
// token and the cookie jar Flask pairs it with.
type supersetSession struct {
	client      *http.Client
	accessToken string
	csrfToken   string
}

// login signs in through Superset's own security API and captures the
// CSRF pairing for write calls.
func (v *SupersetVerifier) login(ctx context.Context, password string) (*supersetSession, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	session := &supersetSession{client: &http.Client{Jar: jar}}

	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		if attempt > 0 {
			time.Sleep(6 * time.Second)
		}
		status, body, err := supersetJsonRequest(ctx, session.client, http.MethodPost, "/api/v1/security/login",
			map[string]interface{}{"username": v.AdminUsername, "password": password, "provider": "db", "refresh": true}, nil)
		if err != nil {
			lastErr = err
			continue
		}
		if status != http.StatusOK {
			lastErr = errors.Errorf("login returned %d: %s", status, body)
			continue
		}
		var payload struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.Unmarshal([]byte(body), &payload); err != nil {
			return nil, errors.Wrap(err, "decoding the login response")
		}
		if payload.AccessToken == "" {
			return nil, errors.New("login answered without an access token")
		}
		session.accessToken = payload.AccessToken

		// The CSRF token pairs with the session cookie the jar now
		// holds — write calls need BOTH (Flask's pairing contract).
		status, body, err = supersetJsonRequest(ctx, session.client, http.MethodGet, "/api/v1/security/csrf_token/", nil, session)
		if err != nil {
			return nil, err
		}
		if status != http.StatusOK {
			return nil, errors.Errorf("the csrf token read returned %d: %s", status, body)
		}
		var csrf struct {
			Result string `json:"result"`
		}
		if err := json.Unmarshal([]byte(body), &csrf); err != nil {
			return nil, errors.Wrap(err, "decoding the csrf response")
		}
		session.csrfToken = csrf.Result
		fmt.Printf("  [verify] signed in as %q through the security API\n", v.AdminUsername)
		return session, nil
	}
	return nil, errors.Wrap(lastErr, "signing in through the security API")
}

// proveDashboard creates a dashboard through the REST API and reads it
// back.
func (v *SupersetVerifier) proveDashboard(ctx context.Context, password, title string) (*supersetSession, string, error) {
	session, err := v.login(ctx, password)
	if err != nil {
		return nil, "", err
	}
	status, body, err := supersetJsonRequest(ctx, session.client, http.MethodPost, "/api/v1/dashboard/",
		map[string]interface{}{"dashboard_title": title, "published": true}, session)
	if err != nil {
		return nil, "", errors.Wrap(err, "creating the proof dashboard")
	}
	if status != http.StatusCreated {
		return nil, "", errors.Errorf("creating the proof dashboard returned %d: %s", status, body)
	}
	var created struct {
		Id json.Number `json:"id"`
	}
	if err := json.Unmarshal([]byte(body), &created); err != nil {
		return nil, "", errors.Wrap(err, "decoding the dashboard-create response")
	}
	dashboardId := created.Id.String()

	got, err := session.dashboardTitle(ctx, dashboardId)
	if err != nil {
		return nil, "", err
	}
	if got != title {
		return nil, "", errors.Errorf("dashboard %s reads back title %q, want %q", dashboardId, got, title)
	}
	return session, dashboardId, nil
}

// dashboardTitle reads one dashboard's title back through the API.
func (s *supersetSession) dashboardTitle(ctx context.Context, dashboardId string) (string, error) {
	status, body, err := supersetJsonRequest(ctx, s.client, http.MethodGet, "/api/v1/dashboard/"+dashboardId, nil, s)
	if err != nil {
		return "", errors.Wrap(err, "reading the dashboard back")
	}
	if status != http.StatusOK {
		return "", errors.Errorf("reading the dashboard returned %d: %s", status, body)
	}
	var payload struct {
		Result struct {
			DashboardTitle string `json:"dashboard_title"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return "", errors.Wrap(err, "decoding the dashboard response")
	}
	return payload.Result.DashboardTitle, nil
}

// deleteDashboard sweeps the proof artifact.
func (s *supersetSession) deleteDashboard(ctx context.Context, dashboardId string) error {
	status, body, err := supersetJsonRequest(ctx, s.client, http.MethodDelete, "/api/v1/dashboard/"+dashboardId, nil, s)
	if err != nil {
		return errors.Wrap(err, "sweeping the proof dashboard")
	}
	if status != http.StatusOK {
		return errors.Errorf("sweeping the proof dashboard returned %d: %s", status, body)
	}
	fmt.Printf("  [verify] proof dashboard %s swept\n", dashboardId)
	return nil
}

// supersetJsonRequest performs one JSON API call. A non-nil session
// adds the JWT bearer, the CSRF header and the Referer Flask-Talisman
// setups expect on write calls.
func supersetJsonRequest(ctx context.Context, client *http.Client, method, path string, payload map[string]interface{}, session *supersetSession) (int, string, error) {
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return 0, "", err
		}
		body = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, supersetBaseUrl()+path, body)
	if err != nil {
		return 0, "", err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if session != nil {
		req.Header.Set("Authorization", "Bearer "+session.accessToken)
		if session.csrfToken != "" {
			req.Header.Set("X-CSRFToken", session.csrfToken)
			req.Header.Set("Referer", supersetBaseUrl()+path)
		}
	}
	resp, err := client.Do(req)
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

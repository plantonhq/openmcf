package verify

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// JupyterHubVerifier checks a JupyterHub installation to the point a
// customer could hand it to a team: the hub and proxy rolled out (the
// chart's CHART-FIXED bare names — hub, proxy, proxy-public), THE AUTH
// GATE (an anonymous hub-API read REJECTED and a wrong-password sign-in
// refused — the chart's own default of any-username-no-password must be
// dead), a REAL SIGN-IN through the login form with the module-generated
// shared password, and THE SPAWN PROOF on every lane — the signed-in
// user's "start my server" creates a REAL single-user pod
// (jupyter-<username>) through KubeSpawner and the hub reports the
// server ready (a notebook platform that cannot spawn a notebook is not
// a notebook platform). The proof artifacts are swept: the server is
// stopped through the API and the user's runtime home PVC
// (claim-<username>) — DESIGNED to survive server stops and destroys —
// is deleted by the verifier for shared-cluster hygiene.
//
// The behavioral-spawn scenario (recognized by name) adds THE STATE
// PROOF: after the spawn proof, the hub pod is deleted (UID-verified
// replacement) and a FRESH sign-in must find the same user account and
// spawn again — hub state (users, server records) lives in the hub
// database (sqlite PVC or composed PostgreSQL), never in the pod.
//
// Destroy: the release resources leave with the release (JupyterHub
// installs no CRDs). Runtime `claim-*` user PVCs survive BY DESIGN —
// VerifyAbsent treats survivors as designed, deletes them, and reports
// the sweep.
type JupyterHubVerifier struct {
	Namespace string
	Name      string
	// SpawnUsername is the identity the proof signs in as (letters and
	// digits only — it becomes part of the user pod's name).
	SpawnUsername string
	// StateProof switches the behavioral arm on.
	StateProof bool
}

const jupyterhubProxyPort = "18081"

// jupyterhubBaseUrl is the tunneled front door (proxy-public:80).
func jupyterhubBaseUrl() string {
	return "http://127.0.0.1:" + jupyterhubProxyPort
}

func (v *JupyterHubVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] jupyterhub %q in namespace %q\n", v.Name, v.Namespace)

	// The core tiers roll out — bare chart-fixed names. The first wait
	// absorbs the pre-puller window on lanes where the hook is off
	// (with the hook on, Helm's own wait already paid for the image
	// pull).
	for _, deployment := range []string{"hub", "proxy"} {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+deployment, v.Namespace, 15*time.Minute); err != nil {
			return errors.Wrapf(err, "the %s deployment never rolled out", deployment)
		}
	}

	// The tunnel backs proxy-public — the ONLY front door (login and
	// user traffic both enter here); listen-waited.
	cancel, err := openServiceTunnel(ctx, kubeconfig, v.Namespace, "proxy-public", jupyterhubProxyPort, "80")
	if err != nil {
		return errors.Wrap(err, "opening the tunnel to proxy-public")
	}

	// THE AUTH GATE, half one: an anonymous hub-API read is rejected.
	if err := v.proveAnonymousRejected(ctx); err != nil {
		cancel()
		return err
	}

	password, err := v.sharedPassword(ctx, kubeconfig)
	if err != nil {
		cancel()
		return err
	}

	// THE AUTH GATE, half two: the WRONG password must not sign in —
	// the chart's own open-door default (any username, NO password)
	// must be dead.
	if err := v.proveWrongPasswordRejected(ctx); err != nil {
		cancel()
		return err
	}

	// A REAL sign-in with the module-generated credential, then THE
	// SPAWN PROOF.
	session, err := v.login(ctx, password)
	if err != nil {
		cancel()
		return err
	}
	if err := v.proveSpawn(ctx, kubeconfig, session, "first"); err != nil {
		cancel()
		return err
	}

	if !v.StateProof {
		err := v.stopServerAndSweep(ctx, kubeconfig, session)
		cancel()
		return err
	}

	// THE STATE PROOF: stop the server, drop the tunnel across the
	// replacement window (fresh-tunnel-per-phase — the shared
	// single-half-open-pipe lesson), replace the hub pod UID-verified,
	// then a FRESH sign-in must find the same account and spawn again —
	// hub state lives in the database, never in the pod.
	if err := v.stopServer(ctx, session); err != nil {
		cancel()
		return err
	}
	cancel()
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
		"component=hub,app=jupyterhub", 10*time.Minute); err != nil {
		return errors.Wrap(err, "the hub pod did not recover after deletion")
	}
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/hub", v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the hub deployment never rolled out after the replacement")
	}

	cancel, err = openServiceTunnel(ctx, kubeconfig, v.Namespace, "proxy-public", jupyterhubProxyPort, "80")
	if err != nil {
		return errors.Wrap(err, "re-establishing the tunnel after the hub replacement")
	}
	defer cancel()
	session, err = v.login(ctx, password)
	if err != nil {
		return errors.Wrap(err, "signing in again after the hub replacement")
	}
	if err := v.proveSpawn(ctx, kubeconfig, session, "post-replacement"); err != nil {
		return errors.Wrap(err, "a spawn AFTER the hub replacement should still succeed")
	}
	fmt.Printf("  [verify] STATE: the same account signed in and spawned again after a UID-verified hub replacement — hub state lives in the database\n")
	return v.stopServerAndSweep(ctx, kubeconfig, session)
}

func (v *JupyterHubVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	for _, deployment := range []string{"hub", "proxy"} {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", deployment, v.Namespace); err != nil {
			return err
		}
	}
	// Runtime user-home PVCs (claim-<username>) survive destroy BY
	// DESIGN — they are the hub's children, not the release's. Treat
	// survivors as designed, then sweep them for shared-cluster
	// hygiene.
	out, err := kubectlGetJSONPath(ctx, kubeconfig, "pvc", "", v.Namespace, "{range .items[*]}{.metadata.name}{\"\\n\"}{end}")
	if err == nil {
		for _, pvcName := range strings.Split(strings.TrimSpace(out), "\n") {
			if strings.HasPrefix(pvcName, "claim-") {
				fmt.Printf("  [verify] DESTROY: runtime user PVC %q survived (designed) — sweeping\n", pvcName)
				_ = kubectlDeleteResource(ctx, kubeconfig, "pvc", pvcName, v.Namespace)
			}
		}
	}
	fmt.Printf("  [verify] DESTROY: the hub and proxy are gone (JupyterHub installs no CRDs — destroy is clean by design)\n")
	return nil
}

// sharedPassword reads the module-generated shared sign-in credential
// from the exported `<name>-auth` Secret (key password).
func (v *JupyterHubVerifier) sharedPassword(ctx context.Context, kubeconfig string) (string, error) {
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
		return "", errors.New("shared password Secret was empty — the secured-by-default contract was not wired")
	}
	return password, nil
}

// jupyterhubSession is an authenticated browser-shaped session: the
// cookie jar carries the hub's signed session cookie and the XSRF
// cookie every state-changing request must echo (JupyterHub requires
// the _xsrf pairing on POSTs since 4.1).
type jupyterhubSession struct {
	client *http.Client
}

// xsrfToken returns the current _xsrf cookie value for the hub path.
func (s *jupyterhubSession) xsrfToken() string {
	base, _ := url.Parse(jupyterhubBaseUrl() + "/hub/")
	for _, cookie := range s.client.Jar.Cookies(base) {
		if cookie.Name == "_xsrf" {
			return cookie.Value
		}
	}
	return ""
}

// newJupyterhubSession primes a session: GET /hub/login sets the _xsrf
// cookie the login POST must echo.
func newJupyterhubSession(ctx context.Context) (*jupyterhubSession, error) {
	jar, _ := cookiejar.New(nil)
	session := &jupyterhubSession{client: &http.Client{Jar: jar, Timeout: 60 * time.Second}}
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		request, _ := http.NewRequestWithContext(ctx, "GET", jupyterhubBaseUrl()+"/hub/login", nil)
		response, err := session.client.Do(request)
		if err == nil {
			_, _ = io.Copy(io.Discard, response.Body)
			response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return session, nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return nil, errors.New("the login page never answered through proxy-public")
}

// login signs in through the hub's own login form as SpawnUsername.
func (v *JupyterHubVerifier) login(ctx context.Context, password string) (*jupyterhubSession, error) {
	session, err := newJupyterhubSession(ctx)
	if err != nil {
		return nil, err
	}
	if err := session.postLogin(ctx, v.SpawnUsername, password); err != nil {
		return nil, err
	}
	// Authenticated proof: the home page answers 200 for a signed-in
	// session (an anonymous request redirects to the login page).
	status, body, err := session.request(ctx, "GET", "/hub/home", nil, "")
	if err != nil {
		return nil, errors.Wrap(err, "the home page never answered after login")
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("sign-in as %q did not stick: /hub/home answered %d: %s", v.SpawnUsername, status, firstLines(body, 2))
	}
	fmt.Printf("  [verify] LOGIN: %q signed in with the module-generated shared password\n", v.SpawnUsername)
	return session, nil
}

// postLogin submits the login form (username/password + the XSRF echo).
func (s *jupyterhubSession) postLogin(ctx context.Context, username, password string) error {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	if xsrf := s.xsrfToken(); xsrf != "" {
		form.Set("_xsrf", xsrf)
	}
	request, _ := http.NewRequestWithContext(ctx, "POST", jupyterhubBaseUrl()+"/hub/login", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := s.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "the login form never answered")
	}
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	// A successful login lands on a 200 page after redirects; a
	// rejected one re-renders the login form with an error and no
	// session cookie — the /hub/home probe distinguishes them.
	if response.StatusCode >= 500 {
		return errors.Errorf("the login form answered %d: %s", response.StatusCode, firstLines(string(body), 2))
	}
	return nil
}

// request performs a session request with the XSRF header echoed (the
// hub accepts the cookie value as the X-XSRFToken header on API calls).
func (s *jupyterhubSession) request(ctx context.Context, method, path string, payload io.Reader, contentType string) (int, string, error) {
	request, err := http.NewRequestWithContext(ctx, method, jupyterhubBaseUrl()+path, payload)
	if err != nil {
		return 0, "", err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if xsrf := s.xsrfToken(); xsrf != "" {
		request.Header.Set("X-XSRFToken", xsrf)
	}
	response, err := s.client.Do(request)
	if err != nil {
		return 0, "", err
	}
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	return response.StatusCode, string(body), nil
}

// proveAnonymousRejected asserts an anonymous hub-API read is rejected.
func (v *JupyterHubVerifier) proveAnonymousRejected(ctx context.Context) error {
	client := &http.Client{Timeout: 30 * time.Second}
	deadline := time.Now().Add(3 * time.Minute)
	var lastStatus int
	for time.Now().Before(deadline) {
		request, _ := http.NewRequestWithContext(ctx, "GET", jupyterhubBaseUrl()+"/hub/api/users", nil)
		response, err := client.Do(request)
		if err == nil {
			_, _ = io.Copy(io.Discard, response.Body)
			response.Body.Close()
			lastStatus = response.StatusCode
			if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
				fmt.Printf("  [verify] AUTH GATE: anonymous hub-API read rejected (%d)\n", response.StatusCode)
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("THE AUTH GATE FAILED: anonymous GET /hub/api/users answered %d, expected 401/403", lastStatus)
}

// proveWrongPasswordRejected asserts a wrong-password sign-in never
// authenticates.
func (v *JupyterHubVerifier) proveWrongPasswordRejected(ctx context.Context) error {
	session, err := newJupyterhubSession(ctx)
	if err != nil {
		return err
	}
	if err := session.postLogin(ctx, v.SpawnUsername, "definitely-not-the-password"); err != nil {
		return err
	}
	status, _, err := session.request(ctx, "GET", "/hub/home", nil, "")
	if err != nil {
		return errors.Wrap(err, "the post-login probe never answered")
	}
	if status == http.StatusOK {
		return errors.New("THE AUTH GATE FAILED: a WRONG password signed in — the chart's open-door default is live")
	}
	fmt.Printf("  [verify] AUTH GATE: wrong-password sign-in refused (the open-door chart default is dead)\n")
	return nil
}

// proveSpawn starts the signed-in user's default server and proves the
// REAL single-user pod runs: KubeSpawner creates `jupyter-<username>`
// and the hub reports the server ready.
func (v *JupyterHubVerifier) proveSpawn(ctx context.Context, kubeconfig string, session *jupyterhubSession, phase string) error {
	// POST /hub/spawn submits the (option-less) spawn form for the
	// session's own user.
	form := url.Values{}
	if xsrf := session.xsrfToken(); xsrf != "" {
		form.Set("_xsrf", xsrf)
	}
	status, body, err := session.request(ctx, "POST", "/hub/spawn/"+v.SpawnUsername, strings.NewReader(form.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return errors.Wrap(err, "the spawn request never answered")
	}
	if status >= 500 {
		return errors.Errorf("the spawn request answered %d: %s", status, firstLines(body, 2))
	}

	// The pod truth: KubeSpawner's user pod reaches Running.
	podName := "jupyter-" + v.SpawnUsername
	deadline := time.Now().Add(10 * time.Minute)
	for time.Now().Before(deadline) {
		podPhase, err := kubectlGetJSONPath(ctx, kubeconfig, "pod", podName, v.Namespace, "{.status.phase}")
		if err == nil && podPhase == "Running" {
			fmt.Printf("  [verify] SPAWN (%s): the user server pod %q is Running — KubeSpawner spawned a real notebook server\n", phase, podName)
			return v.awaitServerReady(ctx, session)
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("THE SPAWN PROOF FAILED (%s): the user pod %q never reached Running", phase, podName)
}

// awaitServerReady polls the hub's own record until the user's default
// server reports ready.
func (v *JupyterHubVerifier) awaitServerReady(ctx context.Context, session *jupyterhubSession) error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastBody string
	for time.Now().Before(deadline) {
		status, body, err := session.request(ctx, "GET", "/hub/api/users/"+v.SpawnUsername, nil, "")
		if err == nil && status == http.StatusOK &&
			(strings.Contains(body, "\"ready\": true") || strings.Contains(body, "\"ready\":true")) {
			fmt.Printf("  [verify] SPAWN: the hub reports the server ready\n")
			return nil
		}
		lastBody = body
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("the hub never reported the server ready: %s", firstLines(lastBody, 3))
}

// stopServer stops the user's default server through the hub API.
func (v *JupyterHubVerifier) stopServer(ctx context.Context, session *jupyterhubSession) error {
	status, body, err := session.request(ctx, "DELETE", "/hub/api/users/"+v.SpawnUsername+"/server", nil, "")
	if err != nil {
		return errors.Wrap(err, "the stop-server request never answered")
	}
	if status >= 300 && status != http.StatusNotFound {
		return errors.Errorf("stopping the server answered %d: %s", status, firstLines(body, 2))
	}
	return nil
}

// stopServerAndSweep stops the server and deletes the user's runtime
// home PVC — proof artifacts leave the shared cluster with the lane.
func (v *JupyterHubVerifier) stopServerAndSweep(ctx context.Context, kubeconfig string, session *jupyterhubSession) error {
	if err := v.stopServer(ctx, session); err != nil {
		return err
	}
	_ = kubectlDeleteResource(ctx, kubeconfig, "pvc", "claim-"+v.SpawnUsername, v.Namespace)
	fmt.Printf("  [verify] SWEEP: the proof server stopped and the runtime home PVC removed\n")
	return nil
}

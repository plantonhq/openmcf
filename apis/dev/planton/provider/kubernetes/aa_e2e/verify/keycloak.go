package verify

import (
	"context"
	"crypto/tls"
	"encoding/base64"
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

// KeycloakVerifier checks an operator-managed Keycloak server to the
// point a customer could sign in: the Keycloak CR reports Ready (its
// boolean-status condition — and Ready=true with HasErrors warnings is a
// LEGITIMATE state the operator documents, so HasErrors is deliberately
// not treated as failure), the StatefulSet reaches the declared instance
// count, the operator's naming contract holds (`<name>-service`,
// `<name>-discovery`, the `<name>-initial-admin` credential Secret) —
// and THE PRODUCT PROOF on every lane: a real admin login (the
// password grant on the master realm with the bootstrap credentials)
// followed by an authenticated admin-API read. An identity server that
// cannot authenticate its own admin is not an identity server.
//
// The behavioral-durability scenario (recognized by name) additionally
// creates a verifier-owned realm, DELETES pod 0, waits for the
// replacement to serve, and re-reads the realm through a fresh login —
// configuration surviving pod replacement is the database-durability
// proof (and exactly what the dev-file sandbox vendors cannot do). The
// verifier sweeps its realm afterwards.
type KeycloakVerifier struct {
	Namespace string
	Name      string
	// Instances is spec.instances (default 1) — the StatefulSet ready
	// target.
	Instances int
	// Https selects the scheme/port of the listener the spec
	// configured (TLS declared -> https, else the http_enabled
	// listener).
	Https bool
	// Port is the listener port (spec override or the 8443/8080
	// defaults).
	Port int
	// AdminSecretName is the bootstrap-admin Secret to read
	// credentials from — `<name>-initial-admin` unless the spec brought
	// its own.
	AdminSecretName string
	// GeneratedAdminSecret marks the operator-generated Secret (its
	// existence is part of the naming contract to assert).
	GeneratedAdminSecret bool
	Behavioral           bool
}

func (v *KeycloakVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] keycloak %q in namespace %q (instances %d, https %v)\n", v.Name, v.Namespace, v.Instances, v.Https)

	// The operator gates the StatefulSet on the admin Secret and first
	// boot runs schema migrations (upstream budgets 10 minutes of
	// startup probe) — poll the CR's own Ready condition first.
	if err := v.waitForReadyCondition(ctx, kubeconfig, 15*time.Minute); err != nil {
		return err
	}

	// The operator's naming contract.
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-service", v.Namespace); err != nil {
		return errors.Wrap(err, "the client service not found")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-discovery", v.Namespace); err != nil {
		return errors.Wrap(err, "the JGroups discovery service not found")
	}
	if v.GeneratedAdminSecret {
		if err := KubectlResourceExists(ctx, kubeconfig, "secret", v.AdminSecretName, v.Namespace); err != nil {
			return errors.Wrap(err, "the operator-generated initial-admin secret not found")
		}
	}

	ready, err := kubectlGetJSONPath(ctx, kubeconfig, "statefulset", v.Name, v.Namespace, "{.status.readyReplicas}")
	if err != nil {
		return errors.Wrap(err, "reading the StatefulSet ready count")
	}
	if ready != fmt.Sprintf("%d", v.Instances) {
		return errors.Errorf("StatefulSet has %s ready replicas, declared %d", ready, v.Instances)
	}

	username, password, err := v.adminCredentials(ctx, kubeconfig)
	if err != nil {
		return err
	}

	return v.proveAdminLogin(ctx, kubeconfig, username, password)
}

func (v *KeycloakVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "keycloaks.k8s.keycloak.org", v.Name, v.Namespace); err != nil {
		return err
	}
	// The operator (destroy-ordered after its CRs) garbage-collects the
	// children; the StatefulSet is the load-bearing one to assert.
	deadline := time.Now().Add(3 * time.Minute)
	for {
		err := KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.Name, v.Namespace)
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return errors.Wrap(err, "the server StatefulSet never deleted after the CR was removed")
		}
		time.Sleep(5 * time.Second)
	}
}

// waitForReadyCondition polls the Keycloak CR until its Ready condition
// reports true. NOTE the operator's conditions carry BOOLEAN statuses
// (not the usual "True" strings) — match both spellings defensively.
func (v *KeycloakVerifier) waitForReadyCondition(ctx context.Context, kubeconfig string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastReady, lastMessage string
	for time.Now().Before(deadline) {
		ready, _ := kubectlGetJSONPath(ctx, kubeconfig, "keycloaks.k8s.keycloak.org", v.Name, v.Namespace,
			`{.status.conditions[?(@.type=="Ready")].status}`)
		lastReady = ready
		if ready == "true" || ready == "True" {
			fmt.Printf("  [verify] Keycloak CR reports Ready\n")
			return nil
		}
		// Surface the operator's own diagnosis while waiting — it
		// inlines crashing pods' log tails into the condition message.
		msg, _ := kubectlGetJSONPath(ctx, kubeconfig, "keycloaks.k8s.keycloak.org", v.Name, v.Namespace,
			`{.status.conditions[?(@.type=="Ready")].message}`)
		if msg != "" {
			lastMessage = msg
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("the Keycloak CR never reported Ready (last status %q, message %q)",
		lastReady, firstLines(lastMessage, 3))
}

// adminCredentials reads the bootstrap admin username/password from the
// credential Secret (the exported handle). The values stay in-process —
// never printed.
func (v *KeycloakVerifier) adminCredentials(ctx context.Context, kubeconfig string) (string, string, error) {
	read := func(key string) (string, error) {
		encoded, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", v.AdminSecretName, v.Namespace,
			fmt.Sprintf("{.data.%s}", key))
		if err != nil {
			return "", errors.Wrapf(err, "reading %q from the admin secret", key)
		}
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", errors.Wrapf(err, "decoding %q from the admin secret", key)
		}
		return string(decoded), nil
	}
	username, err := read("username")
	if err != nil {
		return "", "", err
	}
	password, err := read("password")
	if err != nil {
		return "", "", err
	}
	if username == "" || password == "" {
		return "", "", errors.Errorf("the admin secret %q carries empty credentials", v.AdminSecretName)
	}
	return username, password, nil
}

// proveAdminLogin performs the product proof through a port-forward to
// the client Service: the OIDC password grant on the master realm
// (admin-cli), then an authenticated admin-API realm read. The
// behavioral arm inserts a verifier-owned realm, replaces pod 0, and
// re-reads the realm through a FRESH login and tunnel.
func (v *KeycloakVerifier) proveAdminLogin(ctx context.Context, kubeconfig, username, password string) error {
	// E2E TLS lanes serve self-signed certificates — trust is not what
	// this proof asserts.
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	login := func(base string) (string, error) {
		form := url.Values{
			"grant_type": {"password"},
			"client_id":  {"admin-cli"},
			"username":   {username},
			"password":   {password},
		}
		_, body, err := v.httpForm(ctx, client, base+"/realms/master/protocol/openid-connect/token", form, 6*time.Minute)
		if err != nil {
			return "", errors.Wrap(err, "the admin password grant failed")
		}
		var token struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.Unmarshal([]byte(body), &token); err != nil || token.AccessToken == "" {
			return "", errors.New("the token endpoint answered without an access_token")
		}
		return token.AccessToken, nil
	}

	realm := fmt.Sprintf("e2e-proof-%d", time.Now().Unix())

	err := v.withServiceTunnel(ctx, kubeconfig, func(base string) error {
		token, err := login(base)
		if err != nil {
			return err
		}
		fmt.Printf("  [verify] PRODUCT PROOF: admin password grant succeeded on the master realm\n")

		status, body, err := v.httpJSON(ctx, client, http.MethodGet, base+"/admin/realms", "", token, 2*time.Minute, 200)
		if err != nil {
			return errors.Wrapf(err, "the authenticated admin-API read failed (HTTP %d)", status)
		}
		if !strings.Contains(body, `"master"`) {
			return errors.New("the realm list answered without the master realm")
		}
		fmt.Printf("  [verify] PRODUCT PROOF: authenticated admin API served the realm list\n")

		if v.Behavioral {
			if _, _, err := v.httpJSON(ctx, client, http.MethodPost, base+"/admin/realms",
				fmt.Sprintf(`{"realm": %q, "enabled": true}`, realm), token, 2*time.Minute, 201); err != nil {
				return errors.Wrap(err, "creating the verifier-owned realm")
			}
			fmt.Printf("  [verify] DURABILITY: verifier-owned realm %q created\n", realm)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if !v.Behavioral {
		return nil
	}

	// THE DURABILITY PROOF: configuration lives in the database, so a
	// replaced pod must serve the realm created before its death.
	// Capture the doomed pod's uid first — a terminating pod keeps
	// reporting Ready and the StatefulSet's status lags the deletion,
	// so the recovery is only real once a pod with a NEW uid is Ready.
	pod0 := v.Name + "-0"
	oldUid, err := podUid(ctx, kubeconfig, v.Namespace, pod0)
	if err != nil {
		return err
	}
	fmt.Printf("  [verify] DURABILITY: deleting pod %q\n", pod0)
	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", pod0, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
		return errors.Wrapf(err, "deleting pod 0: %s", string(out))
	}
	if err := waitForPodReplaced(ctx, kubeconfig, v.Namespace, pod0, oldUid, true, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the replacement pod never became Ready")
	}
	if err := v.waitForReadyStatefulSet(ctx, kubeconfig, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the StatefulSet never recovered after the pod loss")
	}

	// Fresh tunnel, fresh login — the old tunnel died with its pod.
	return v.withServiceTunnel(ctx, kubeconfig, func(base string) error {
		token, err := login(base)
		if err != nil {
			return errors.Wrap(err, "re-login after the pod replacement failed")
		}
		if _, _, err := v.httpJSON(ctx, client, http.MethodGet, base+"/admin/realms/"+realm, "", token, 4*time.Minute, 200); err != nil {
			return errors.Wrap(err, "the verifier-owned realm did NOT survive the pod replacement")
		}
		fmt.Printf("  [verify] DURABILITY: realm %q served after pod replacement — configuration survived in the database\n", realm)
		// Sweep the verifier-owned realm.
		if _, _, err := v.httpJSON(ctx, client, http.MethodDelete, base+"/admin/realms/"+realm, "", token, 2*time.Minute, 204); err != nil {
			return errors.Wrap(err, "sweeping the verifier-owned realm")
		}
		return nil
	})
}

func (v *KeycloakVerifier) waitForReadyStatefulSet(ctx context.Context, kubeconfig string, budget time.Duration) error {
	want := fmt.Sprintf("%d", v.Instances)
	deadline := time.Now().Add(budget)
	var last string
	for time.Now().Before(deadline) {
		ready, _ := kubectlGetJSONPath(ctx, kubeconfig, "statefulset", v.Name, v.Namespace, "{.status.readyReplicas}")
		last = ready
		if ready == want {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("StatefulSet never returned to %s ready replicas (last %q)", want, last)
}

// withServiceTunnel runs fn with a fresh port-forward to the client
// Service on the configured listener port. Fresh per call: a tunnel dies
// silently with its backing pod (the caught-live port-forward class).
func (v *KeycloakVerifier) withServiceTunnel(ctx context.Context, kubeconfig string, fn func(base string) error) error {
	const localPort = "18443"

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.Name+"-service", fmt.Sprintf("%s:%d", localPort, v.Port), "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the client service")
	}
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	scheme := "http"
	if v.Https {
		scheme = "https"
	}
	return fn(fmt.Sprintf("%s://127.0.0.1:%s", scheme, localPort))
}

// httpForm POSTs a form (the token endpoint's content type), retrying
// across tunnel warm-up until 200 or the budget expires.
func (v *KeycloakVerifier) httpForm(ctx context.Context, client *http.Client, endpoint string, form url.Values, budget time.Duration) (int, string, error) {
	deadline := time.Now().Add(budget)
	var lastStatus int
	var lastBody string
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
		if err != nil {
			return 0, "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err == nil {
			body := drainBody(resp)
			lastStatus = resp.StatusCode
			lastBody = body
			if resp.StatusCode == http.StatusOK {
				return resp.StatusCode, body, nil
			}
			lastErr = errors.Errorf("HTTP %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(5 * time.Second)
	}
	return lastStatus, lastBody, errors.Wrapf(lastErr, "last body: %s", firstLines(lastBody, 2))
}

// drainBody reads and closes a response body.
func drainBody(resp *http.Response) string {
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(body)
}

// keycloakScenarioShape pulls the verifier's inputs out of a
// KubernetesKeycloak scenario manifest: the instance count, which
// listener the spec configured (TLS declared -> https and its port,
// else the plain-HTTP listener), and the bootstrap-admin Secret
// (spec-declared, or the operator-generated `<name>-initial-admin`).
// Scenario manifests use the snake_case field convention; the camelCase
// twin is tolerated.
func keycloakScenarioShape(spec map[string]interface{}, name string) (instances int, https bool, port int, adminSecret string, generated bool) {
	instances = 1
	if n, ok := specInt(spec["instances"]); ok {
		instances = n
	}

	httpBlock := specNestedMap(spec, "http")
	if ref := specNestedMap(httpBlock, "tls_secret_name", "tlsSecretName"); ref != nil {
		if s, _ := ref["value"].(string); s != "" {
			https = true
		}
		if _, ok := ref["value_from"]; ok {
			https = true
		}
		if _, ok := ref["valueFrom"]; ok {
			https = true
		}
	}
	if https {
		port = 8443
		if n, ok := specInt(firstNonNil(httpBlock, "https_port", "httpsPort")); ok {
			port = n
		}
	} else {
		port = 8080
		if n, ok := specInt(firstNonNil(httpBlock, "http_port", "httpPort")); ok {
			port = n
		}
	}

	if s, _ := firstNonNil(spec, "bootstrap_admin_secret_name", "bootstrapAdminSecretName").(string); s != "" {
		adminSecret = s
		return
	}
	adminSecret = name + "-initial-admin"
	generated = true
	return
}

// specNestedMap reads spec[key] as a map, trying each key spelling.
func specNestedMap(spec map[string]interface{}, keys ...string) map[string]interface{} {
	if spec == nil {
		return nil
	}
	for _, key := range keys {
		if m, ok := spec[key].(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// firstNonNil returns the first present value among the key spellings.
func firstNonNil(spec map[string]interface{}, keys ...string) interface{} {
	if spec == nil {
		return nil
	}
	for _, key := range keys {
		if v, ok := spec[key]; ok {
			return v
		}
	}
	return nil
}

// httpJSON performs one JSON request with retries, succeeding on any of
// wantStatuses.
func (v *KeycloakVerifier) httpJSON(ctx context.Context, client *http.Client, method, endpoint, body, token string, budget time.Duration, wantStatuses ...int) (int, string, error) {
	deadline := time.Now().Add(budget)
	var lastStatus int
	var lastBody string
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, method, endpoint, strings.NewReader(body))
		if err != nil {
			return 0, "", err
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err == nil {
			out := drainBody(resp)
			lastStatus = resp.StatusCode
			lastBody = out
			for _, want := range wantStatuses {
				if resp.StatusCode == want {
					return resp.StatusCode, out, nil
				}
			}
			lastErr = errors.Errorf("HTTP %d (wanted one of %v)", resp.StatusCode, wantStatuses)
		} else {
			lastErr = err
		}
		time.Sleep(5 * time.Second)
	}
	return lastStatus, lastBody, errors.Wrapf(lastErr, "last body: %s", firstLines(lastBody, 2))
}

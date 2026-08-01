package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OpenFgaVerifier checks an OpenFGA deployment to the point an
// application could authorize against it: the server Deployment rolled
// out, the chart's naming contract holds (the `<name>` client Service on
// the HTTP API port) — and THE ZANZIBAR PROOF on every lane: a live
// store → authorization-model → tuple → Check round-trip, asserting BOTH
// decisions (the granted user is ALLOWED, an ungranted user is DENIED —
// an authorization engine that cannot say no is not an authorization
// engine).
//
// When the scenario declares pre-shared-key authentication the proof
// runs AS an authenticated client and additionally asserts THE AUTH
// GATE: the same API call WITHOUT the key is rejected with 401.
type OpenFgaVerifier struct {
	Namespace string
	Name      string
	// ApiKey is the first declared pre-shared key ("" = the
	// no-authentication posture).
	ApiKey string
}

func (v *OpenFgaVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] openfga %q in namespace %q (authn: %v)\n", v.Name, v.Namespace, v.ApiKey != "")

	// The migration init container gates every pod on the datastore
	// being reachable and migrated — the rollout wait absorbs it.
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the server deployment never rolled out")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "the client service not found")
	}

	return v.proveZanzibarRoundTrip(ctx, kubeconfig)
}

func (v *OpenFgaVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace)
}

// fgaModel is the smallest useful authorization model: documents with
// directly-assignable viewers.
const fgaModel = `{
  "schema_version": "1.1",
  "type_definitions": [
    {"type": "user"},
    {
      "type": "document",
      "relations": {"viewer": {"this": {}}},
      "metadata": {"relations": {"viewer": {"directly_related_user_types": [{"type": "user"}]}}}
    }
  ]
}`

// proveZanzibarRoundTrip drives the HTTP API over a port-forward:
// create a run-unique store, write the model, grant anne viewer on a
// document, then Check BOTH ways. With authentication declared, the
// unauthenticated rejection is asserted FIRST.
func (v *OpenFgaVerifier) proveZanzibarRoundTrip(ctx context.Context, kubeconfig string) error {
	const localPort = "18080"

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.Name, localPort+":8080", "-n", v.Namespace)
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

	base := "http://127.0.0.1:" + localPort
	client := &http.Client{Timeout: 30 * time.Second}

	// THE AUTH GATE: with keys declared, the API without a key must
	// answer 401 — proving the keys actually guard the surface.
	if v.ApiKey != "" {
		if _, _, err := v.fgaCall(ctx, client, http.MethodGet, base+"/stores", "", "", 4*time.Minute, 401); err != nil {
			return errors.Wrap(err, "the unauthenticated request was NOT rejected (the auth gate)")
		}
		fmt.Printf("  [verify] AUTH GATE: unauthenticated request rejected with 401\n")
	}

	// Store create → model write → tuple write → both Check decisions.
	storeName := fmt.Sprintf("e2e-proof-%d", time.Now().Unix())
	_, body, err := v.fgaCall(ctx, client, http.MethodPost, base+"/stores",
		fmt.Sprintf(`{"name": %q}`, storeName), v.ApiKey, 4*time.Minute, 201, 200)
	if err != nil {
		return errors.Wrap(err, "creating the proof store")
	}
	var store struct {
		Id string `json:"id"`
	}
	if err := json.Unmarshal([]byte(body), &store); err != nil || store.Id == "" {
		return errors.New("the store create answered without a store id")
	}

	_, body, err = v.fgaCall(ctx, client, http.MethodPost, base+"/stores/"+store.Id+"/authorization-models",
		fgaModel, v.ApiKey, 2*time.Minute, 201, 200)
	if err != nil {
		return errors.Wrap(err, "writing the authorization model")
	}
	var model struct {
		Id string `json:"authorization_model_id"`
	}
	if err := json.Unmarshal([]byte(body), &model); err != nil || model.Id == "" {
		return errors.New("the model write answered without an authorization_model_id")
	}

	if _, _, err := v.fgaCall(ctx, client, http.MethodPost, base+"/stores/"+store.Id+"/write",
		fmt.Sprintf(`{"writes": {"tuple_keys": [{"user": "user:anne", "relation": "viewer", "object": "document:budget"}]}, "authorization_model_id": %q}`, model.Id),
		v.ApiKey, 2*time.Minute, 200); err != nil {
		return errors.Wrap(err, "writing the relationship tuple")
	}

	check := func(user string) (bool, error) {
		_, body, err := v.fgaCall(ctx, client, http.MethodPost, base+"/stores/"+store.Id+"/check",
			fmt.Sprintf(`{"tuple_key": {"user": %q, "relation": "viewer", "object": "document:budget"}, "authorization_model_id": %q}`, user, model.Id),
			v.ApiKey, 2*time.Minute, 200)
		if err != nil {
			return false, err
		}
		var decision struct {
			Allowed bool `json:"allowed"`
		}
		if err := json.Unmarshal([]byte(body), &decision); err != nil {
			return false, errors.Wrap(err, "parsing the check response")
		}
		return decision.Allowed, nil
	}

	allowed, err := check("user:anne")
	if err != nil {
		return errors.Wrap(err, "checking the granted user")
	}
	if !allowed {
		return errors.New("anne was granted viewer but Check answered DENIED")
	}
	denied, err := check("user:bob")
	if err != nil {
		return errors.Wrap(err, "checking the ungranted user")
	}
	if denied {
		return errors.New("bob was never granted anything but Check answered ALLOWED")
	}
	fmt.Printf("  [verify] ZANZIBAR PROOF: store %q — anne ALLOWED, bob DENIED\n", storeName)

	// Sweep the verifier-owned store.
	if _, _, err := v.fgaCall(ctx, client, http.MethodDelete, base+"/stores/"+store.Id, "", v.ApiKey, 2*time.Minute, 204, 200); err != nil {
		return errors.Wrap(err, "sweeping the proof store")
	}
	return nil
}

// openFgaPresharedKey pulls the first declared pre-shared key out of a
// KubernetesOpenFga scenario manifest ("" = no authentication declared).
// Scenario manifests use the snake_case field convention.
func openFgaPresharedKey(spec map[string]interface{}) string {
	authn := specNestedMap(spec, "authn")
	preshared := specNestedMap(authn, "preshared")
	if preshared == nil {
		return ""
	}
	if keys, ok := preshared["keys"].([]interface{}); ok && len(keys) > 0 {
		key, _ := keys[0].(string)
		return key
	}
	return ""
}

// fgaCall performs one JSON API call with retries across the tunnel
// warm-up, succeeding on any of wantStatuses. The pre-shared key rides
// the Authorization: Bearer header (the server's preshared contract).
func (v *OpenFgaVerifier) fgaCall(ctx context.Context, client *http.Client, method, endpoint, body, apiKey string, budget time.Duration, wantStatuses ...int) (int, string, error) {
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
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
		resp, err := client.Do(req)
		if err == nil {
			out, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			lastStatus = resp.StatusCode
			lastBody = string(out)
			for _, want := range wantStatuses {
				if resp.StatusCode == want {
					return resp.StatusCode, lastBody, nil
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

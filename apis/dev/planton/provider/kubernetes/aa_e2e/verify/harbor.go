package verify

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
)

// HarborVerifier checks a Harbor registry to the point a customer could
// push and pull OCI artifacts against it: every stateless component
// rolled out (core, portal, registry, jobservice, nginx), the front-door
// Service present, and THE REGISTRY PROOF on every lane — login as the
// module-generated admin, create a project, OCI push/pull round-trip,
// asserting BOTH the granted path and THE AUTH GATE (an unauthenticated
// API call is rejected with 401).
//
// The behavioral-durability scenario (recognized by name) additionally
// DELETES the registry pod after the push, waits for a UID-verified
// replacement, re-establishes the port-forward, and pulls the artifact
// back byte-identical — registry blobs surviving pod loss through the
// filesystem PVC. The tunnel targets the nginx front-door Service, so
// the registry deletion does not kill it — the re-tunnel is deliberate
// defense against the dead-tunnel class (a mid-proof nginx restart
// would silently break the old tunnel, and a fresh one costs nothing).
//
// externalUrl in the scenario MUST match the verifier's port-forward
// address (http://127.0.0.1:18080): Harbor embeds it in the token-service
// URL returned to every OCI client, so push/pull FAIL AUTH when the
// dialed address and externalUrl disagree.
type HarborVerifier struct {
	Namespace  string
	Name       string
	Durability bool
}

// harborLocalPort is the workstation side of the port-forward. Kept in
// lockstep with every scenario's externalUrl.
const harborLocalPort = "18080"

func (v *HarborVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] harbor %q in namespace %q\n", v.Name, v.Namespace)

	// Stateless components self-ready; core's first-boot schema migration
	// is the long pole (startup probe budgets 60 minutes upstream — the
	// Helm wait already absorbed the normal case).
	for _, component := range []string{"core", "portal", "registry", "jobservice", "nginx"} {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-"+component, v.Namespace, 15*time.Minute); err != nil {
			return errors.Wrapf(err, "the %s deployment never rolled out", component)
		}
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "the front-door Service not found")
	}

	password, err := v.adminPassword(ctx, kubeconfig)
	if err != nil {
		return err
	}

	return v.proveRegistryRoundTrip(ctx, kubeconfig, password)
}

func (v *HarborVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name+"-core", v.Namespace)
}

// adminPassword reads the module-generated admin password from the
// exported `<name>-admin-auth` Secret (key HARBOR_ADMIN_PASSWORD).
func (v *HarborVerifier) adminPassword(ctx context.Context, kubeconfig string) (string, error) {
	secretName := v.Name + "-admin-auth"
	b64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", secretName, v.Namespace, "{.data.HARBOR_ADMIN_PASSWORD}")
	if err != nil {
		return "", errors.Wrapf(err, "reading secret %q HARBOR_ADMIN_PASSWORD", secretName)
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	password := strings.TrimSpace(string(raw))
	if password == "" {
		return "", errors.New("admin password Secret was empty")
	}
	// The chart's publicly documented default must never ship.
	if password == "Harbor12345" {
		return "", errors.New("admin password is the chart's public default Harbor12345 — the module-generated Secret was not wired")
	}
	return password, nil
}

// proveRegistryRoundTrip drives the Harbor API + OCI distribution over a
// port-forward to the front-door Service: assert the auth gate, create a
// run-unique project, push a random image, pull it back. The durability
// variant deletes the registry pod between push and pull and opens a
// fresh tunnel after replacement.
func (v *HarborVerifier) proveRegistryRoundTrip(ctx context.Context, kubeconfig, password string) error {
	cancel, err := v.startPortForward(ctx, kubeconfig)
	if err != nil {
		return err
	}

	base := "http://127.0.0.1:" + harborLocalPort
	client := &http.Client{Timeout: 30 * time.Second}

	// THE AUTH GATE: the API without credentials must answer 401 —
	// proving the generated admin password actually guards the surface.
	if err := v.waitForAuthGate(ctx, client, base, 4*time.Minute); err != nil {
		cancel()
		return err
	}
	fmt.Printf("  [verify] AUTH GATE: unauthenticated request rejected with 401\n")

	project := fmt.Sprintf("e2e%d", time.Now().Unix()%1000000)
	if err := v.createProject(ctx, client, base, password, project, 4*time.Minute); err != nil {
		cancel()
		return err
	}
	fmt.Printf("  [verify] PROJECT: %q created as admin\n", project)

	repo := "proof"
	tag := fmt.Sprintf("t%d", time.Now().Unix())
	refStr := fmt.Sprintf("127.0.0.1:%s/%s/%s:%s", harborLocalPort, project, repo, tag)
	digest, err := v.pushImage(ctx, refStr, password)
	if err != nil {
		cancel()
		return err
	}
	fmt.Printf("  [verify] OCI: pushed %s (digest %s)\n", refStr, digest)

	if v.Durability {
		// Close the tunnel across the replacement window and reopen fresh
		// after — see the type comment on the dead-tunnel defense.
		cancel()
		if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
			"app.kubernetes.io/instance="+v.Name+",app.kubernetes.io/component=registry", 10*time.Minute); err != nil {
			return errors.Wrap(err, "the registry pod did not recover after deletion")
		}
		// Fresh tunnel — the previous one died with the pod.
		cancel, err = v.startPortForward(ctx, kubeconfig)
		if err != nil {
			return errors.Wrap(err, "re-establishing the port-forward after registry replacement")
		}
	}
	defer cancel()

	gotDigest, err := v.pullImage(ctx, refStr, password)
	if err != nil {
		return err
	}
	if gotDigest != digest {
		return errors.Errorf("pulled digest %s does not match pushed digest %s", gotDigest, digest)
	}
	if v.Durability {
		fmt.Printf("  [verify] STATE: artifact survived registry pod replacement AND pulled byte-identical — blobs live on the PVC\n")
	} else {
		fmt.Printf("  [verify] OCI: pulled %s (digest match)\n", refStr)
	}
	return nil
}

func (v *HarborVerifier) startPortForward(ctx context.Context, kubeconfig string) (context.CancelFunc, error) {
	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.Name, harborLocalPort+":80", "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return nil, errors.Wrap(err, "starting port-forward to the front-door Service")
	}
	// Cancel FIRST — Wait blocks forever on a port-forward that is never
	// told to exit.
	done := make(chan struct{})
	go func() {
		_ = pf.Wait()
		close(done)
	}()
	return func() {
		cancel()
		<-done
	}, nil
}

// waitForAuthGate retries until Harbor answers 401 to an unauthenticated
// projects list (core finishes migrating before the API authenticates).
func (v *HarborVerifier) waitForAuthGate(ctx context.Context, client *http.Client, base string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var last error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/api/v2.0/projects", nil)
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
		if resp.StatusCode == http.StatusUnauthorized {
			return nil
		}
		last = errors.Errorf("expected 401, got %d", resp.StatusCode)
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(last, "the unauthenticated request was NOT rejected (the auth gate)")
}

func (v *HarborVerifier) createProject(ctx context.Context, client *http.Client, base, password, project string, budget time.Duration) error {
	body := fmt.Sprintf(`{"project_name":%q,"public":true}`, project)
	deadline := time.Now().Add(budget)
	var last error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/v2.0/projects", strings.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("admin", password)
		resp, err := client.Do(req)
		if err != nil {
			last = err
			time.Sleep(5 * time.Second)
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		// 201 created; 409 already exists (retry-safe).
		if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusConflict {
			return nil
		}
		last = errors.Errorf("create project: status %d: %s", resp.StatusCode, string(respBody))
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(last, "creating the proof project")
}

func (v *HarborVerifier) pushImage(ctx context.Context, refStr, password string) (string, error) {
	ref, err := name.ParseReference(refStr, name.Insecure)
	if err != nil {
		return "", errors.Wrap(err, "parsing push reference")
	}
	img, err := random.Image(1024, 1)
	if err != nil {
		return "", errors.Wrap(err, "building the proof image")
	}
	auth := &authn.Basic{Username: "admin", Password: password}
	if err := remote.Write(ref, img, remote.WithAuth(auth), remote.WithContext(ctx)); err != nil {
		return "", errors.Wrap(err, "the OCI push never succeeded")
	}
	digest, err := img.Digest()
	if err != nil {
		return "", err
	}
	return digest.String(), nil
}

func (v *HarborVerifier) pullImage(ctx context.Context, refStr, password string) (string, error) {
	ref, err := name.ParseReference(refStr, name.Insecure)
	if err != nil {
		return "", errors.Wrap(err, "parsing pull reference")
	}
	auth := &authn.Basic{Username: "admin", Password: password}
	img, err := remote.Image(ref, remote.WithAuth(auth), remote.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "the OCI pull never succeeded")
	}
	digest, err := img.Digest()
	if err != nil {
		return "", err
	}
	// Touch a config read so the pull is not a no-op against a stub.
	if _, err := img.ConfigName(); err != nil {
		return "", errors.Wrap(err, "reading the pulled image config")
	}
	return digest.String(), nil
}

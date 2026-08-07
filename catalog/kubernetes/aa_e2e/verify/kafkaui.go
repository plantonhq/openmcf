package verify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// KafkaUiVerifier checks a kafbat UI console: the Deployment Available and
// the Service present.
//
// When ClusterOnline is set (the behavioral-observe scenario, which chains
// a live Kafka fixture), the verifier additionally proves the console is
// actually OBSERVING: the UI's own /api/clusters endpoint must report the
// declared cluster with status "online" — a console that renders but
// cannot reach its cluster is exactly the failure mode install checks
// cannot see.
type KafkaUiVerifier struct {
	Namespace string
	// ServiceName is the chart-derived Service name.
	ServiceName string
	Port        int
	// ClusterOnline names the declared cluster the console must report
	// online (empty = install-grade assertions only).
	ClusterOnline string
	// Authenticated marks scenarios with login_form auth: the API then
	// answers 401/302 for anonymous calls, which IS the auth proof.
	Authenticated bool
}

func (v *KafkaUiVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kafbat UI %q in namespace %q\n", v.ServiceName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "deployment", v.ServiceName, v.Namespace,
		"condition=Available", 6*time.Minute); err != nil {
		return errors.Wrapf(err, "console deployment %q never became Available", v.ServiceName)
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.ServiceName, v.Namespace); err != nil {
		return errors.Wrap(err, "console service not found")
	}
	return v.probeApi(ctx, kubeconfig)
}

func (v *KafkaUiVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.ServiceName, v.Namespace); err != nil {
		return err
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "service", v.ServiceName, v.Namespace)
}

// probeApi drives the console's own API over a port-forward.
func (v *KafkaUiVerifier) probeApi(ctx context.Context, kubeconfig string) error {
	const localPort = "18082"
	port := v.Port
	if port == 0 {
		port = 80
	}

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.ServiceName, fmt.Sprintf("%s:%d", localPort, port), "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the console service")
	}
	// ONE deferred func, cancel FIRST — Wait blocks forever otherwise.
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	// The console's Spring app can answer the probe a little after the
	// deployment turns Available — retry rather than sleeping blind.
	client := &http.Client{
		// A login_form console 302-redirects anonymous API calls to the
		// login page; the redirect itself is the signal we assert on.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	deadline := time.Now().Add(3 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			"http://127.0.0.1:"+localPort+"/api/clusters", nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if v.Authenticated {
			// Anonymous access must be REFUSED — 401 or a redirect to
			// the login form both prove the auth gate is on.
			if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusFound {
				fmt.Printf("  [verify] AUTH-GATE: anonymous /api/clusters refused with HTTP %d — login is enforced\n", resp.StatusCode)
				return nil
			}
			lastErr = errors.Errorf("expected the auth gate to refuse anonymous access, got HTTP %d", resp.StatusCode)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = errors.Errorf("/api/clusters returned HTTP %d", resp.StatusCode)
			time.Sleep(5 * time.Second)
			continue
		}
		if v.ClusterOnline == "" {
			fmt.Printf("  [verify] console API serving (/api/clusters HTTP 200)\n")
			return nil
		}
		// The observe proof: the declared cluster listed AND online. The
		// API serializes the status enum UPPERCASE ("ONLINE") — verified
		// live against the running console.
		payload := string(body)
		if strings.Contains(payload, fmt.Sprintf(`"name":"%s"`, v.ClusterOnline)) &&
			strings.Contains(payload, `"status":"ONLINE"`) {
			fmt.Printf("  [verify] OBSERVING: console reports cluster %q online\n", v.ClusterOnline)
			return nil
		}
		lastErr = errors.Errorf("cluster %q not reported online yet: %s", v.ClusterOnline, firstLines(payload, 2))
		time.Sleep(10 * time.Second)
	}
	return errors.Wrapf(lastErr, "the console API probe never succeeded (port-forward output: %s)", pfOut.String())
}

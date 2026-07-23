package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// AuthzDenyBehavioralVerifier proves AuthorizationPolicy ENFORCEMENT in a
// real mesh: a meshed client's in-cluster request to the protected backend
// must come back 403 (Envoy's RBAC filter), and after the policy is
// destroyed the same request must succeed again — the full
// enforcement-and-release cycle, not just object existence.
type AuthzDenyBehavioralVerifier struct {
	Namespace  string
	PolicyName string
	// ClientDeployment is the meshed curl client the probe execs into.
	ClientDeployment string
	// BackendURL is the in-cluster URL of the protected Service.
	BackendURL string
}

func (v *AuthzDenyBehavioralVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral authz deny: policy %q must 403 %q\n", v.PolicyName, v.BackendURL)

	if err := KubectlResourceExists(ctx, kubeconfig,
		"authorizationpolicies.security.istio.io", v.PolicyName, v.Namespace); err != nil {
		return err
	}

	// Sidecar readiness gates everything: 2/2 containers on both workloads.
	for _, deploy := range []string{v.ClientDeployment, "e2e-authz-backend"} {
		if err := kubectlWait(ctx, kubeconfig, "deployment", deploy, v.Namespace,
			"condition=Available", 4*time.Minute); err != nil {
			return errors.Wrapf(err, "meshed workload %q not available", deploy)
		}
	}

	// Policy propagation to the sidecar takes seconds — poll for the DENY.
	return v.pollStatus(ctx, kubeconfig, "403", 3*time.Minute,
		"policy never enforced (request kept succeeding)")
}

func (v *AuthzDenyBehavioralVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig,
		"authorizationpolicies.security.istio.io", v.PolicyName, v.Namespace); err != nil {
		return err
	}
	// The release proof: with the policy gone the same request succeeds —
	// enforcement really came from the destroyed object. The fixtures (and
	// their namespace) outlive the component under test, so the probe still
	// runs.
	return v.pollStatus(ctx, kubeconfig, "200", 3*time.Minute,
		"request kept failing after the policy was destroyed")
}

// pollStatus execs curl inside the meshed client until the backend returns
// the wanted HTTP status.
func (v *AuthzDenyBehavioralVerifier) pollStatus(ctx context.Context, kubeconfig, want string, timeout time.Duration, failMsg string) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"exec", "deploy/"+v.ClientDeployment, "-n", v.Namespace, "-c", "app", "--",
			"curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "--max-time", "5", v.BackendURL)
		out, err := cmd.CombinedOutput()
		status := strings.TrimSpace(string(out))
		if err == nil && status == want {
			fmt.Printf("  [verify] in-mesh request returned HTTP %s — as required\n", status)
			return nil
		}
		last = fmt.Sprintf("status=%q err=%v", status, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("%s (last probe: %s)", failMsg, last)
}

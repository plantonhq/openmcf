package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// NetworkPolicyDenyBehavioralVerifier proves NetworkPolicy ENFORCEMENT on a
// cluster whose CNI actually enforces (Cilium): with the deny policy in
// place, a client pod's in-cluster request to the selected backend must
// FAIL (the packet is dropped — curl times out; there is no HTTP status
// because nothing answers), and after the policy is destroyed the same
// request must succeed — the full enforce-and-release cycle. On the default
// kind cluster this proof is impossible (kindnet ignores NetworkPolicy),
// which is exactly why these scenarios carry the cilium-cni cluster
// profile.
type NetworkPolicyDenyBehavioralVerifier struct {
	Namespace  string
	PolicyName string
	// ClientDeployment is the curl client the probe execs into.
	ClientDeployment string
	// BackendURL is the in-cluster URL of the policy-protected Service.
	BackendURL string
}

func (v *NetworkPolicyDenyBehavioralVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral netpol deny: policy %q must block %q\n", v.PolicyName, v.BackendURL)

	if err := KubectlResourceExists(ctx, kubeconfig, "networkpolicy", v.PolicyName, v.Namespace); err != nil {
		return err
	}

	for _, deploy := range []string{v.ClientDeployment, "e2e-netpol-backend"} {
		if err := kubectlWait(ctx, kubeconfig, "deployment", deploy, v.Namespace,
			"condition=Available", 3*time.Minute); err != nil {
			return errors.Wrapf(err, "workload %q not available", deploy)
		}
	}

	// Enforcement proof: the request must FAIL. Curl exits non-zero on a
	// dropped connection (--max-time turns the silent drop into a timeout).
	// Polling absorbs the seconds Cilium takes to compile the policy into
	// the datapath.
	return v.pollRequest(ctx, kubeconfig, false, 3*time.Minute,
		"policy never enforced (request kept succeeding — is the CNI enforcing NetworkPolicy?)")
}

func (v *NetworkPolicyDenyBehavioralVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "networkpolicy", v.PolicyName, v.Namespace); err != nil {
		return err
	}
	// The release proof: with the policy destroyed the same request
	// succeeds — the block really came from the destroyed object. The
	// fixtures (and their namespace) outlive the component under test, so
	// the probe still runs.
	return v.pollRequest(ctx, kubeconfig, true, 3*time.Minute,
		"request kept failing after the policy was destroyed")
}

// pollRequest execs curl inside the client until the request outcome matches
// wantSuccess.
func (v *NetworkPolicyDenyBehavioralVerifier) pollRequest(ctx context.Context, kubeconfig string, wantSuccess bool, timeout time.Duration, failMsg string) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"exec", "deploy/"+v.ClientDeployment, "-n", v.Namespace, "--",
			"curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "--max-time", "5", v.BackendURL)
		out, err := cmd.CombinedOutput()
		status := strings.TrimSpace(string(out))
		succeeded := err == nil && status == "200"
		if succeeded == wantSuccess {
			if wantSuccess {
				fmt.Printf("  [verify] in-cluster request returned HTTP 200 — traffic flows\n")
			} else {
				fmt.Printf("  [verify] in-cluster request failed (status=%q) — traffic blocked, as required\n", status)
			}
			return nil
		}
		last = fmt.Sprintf("status=%q err=%v", status, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("%s (last probe: %s)", failMsg, last)
}

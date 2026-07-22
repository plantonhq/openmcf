package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// CiliumInstallVerifier checks a Cilium installation all the way to a
// FUNCTIONING dataplane: the cilium agent DaemonSet fully rolled out, the
// cilium-operator Deployment Available, and — the load-bearing assertion —
// every node Ready. On the CNI-less cluster profile the nodes are NotReady
// BY DESIGN until Cilium installs, so the NotReady→Ready transition IS the
// proof that Cilium actually became the cluster's CNI, not just that its
// pods started.
type CiliumInstallVerifier struct {
	Namespace string
}

func (v *CiliumInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] cilium installation in namespace %q\n", v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for cilium", v.Namespace)
	}

	// The chart names its workloads fixed ("cilium" / "cilium-operator")
	// regardless of release name — one dataplane per cluster.
	if err := kubectlRolloutStatus(ctx, kubeconfig, "daemonset/cilium", v.Namespace, 4*time.Minute); err != nil {
		return errors.Wrap(err, "cilium agent daemonset never fully rolled out")
	}

	if err := kubectlWait(ctx, kubeconfig, "deployment", "cilium-operator", v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "cilium-operator deployment not available")
	}

	// Nodes Ready is the CNI contract: kubelet flips a node Ready only once
	// a CNI configuration exists and initializes. On the cilium-cni profile
	// this is the moment the cluster becomes schedulable at all.
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"wait", "node", "--all", "--for=condition=Ready", "--timeout=3m")
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Errorf("nodes never became Ready with cilium installed: %v: %s", err, string(out))
	}

	fmt.Printf("  [verify] cilium agent rolled out, operator available, all nodes Ready\n")
	return nil
}

func (v *CiliumInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "daemonset", "cilium", v.Namespace); err != nil {
		return err
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", "cilium-operator", v.Namespace)
}

// kubectlRolloutStatus waits for a workload's rollout to complete. DaemonSet
// readiness has no single kubectl-wait condition (numberReady is a numeric
// status field, not a condition), and `rollout status` is exactly the
// primitive that encodes "all desired pods updated and ready".
func kubectlRolloutStatus(ctx context.Context, kubeconfig, workload, namespace string, timeout time.Duration) error {
	args := []string{"--kubeconfig", kubeconfig, "rollout", "status", workload,
		"--timeout", timeout.String()}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("rollout status %s: %v: %s", workload, err, strings.TrimSpace(string(out)))
	}
	return nil
}

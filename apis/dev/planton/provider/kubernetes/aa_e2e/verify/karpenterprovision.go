package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// KarpenterNodePoolVerifier verifies a NodePool object and, in behavioral
// mode, proves Karpenter actually LAUNCHES machines: a verifier-owned driver
// pod that only the pool under test can satisfy (its nodeSelector targets
// the pool's template label; its toleration matches the pool's fencing
// taint) must trigger a real EC2 node launch, schedule onto the new node,
// and — after the driver is deleted — consolidation must reclaim the node.
// The driver is verifier-owned because it must appear AFTER the pool under
// test exists (fixtures deploy before the component).
type KarpenterNodePoolVerifier struct {
	Name string
	// Behavioral switches on the live node-launch proof (see above).
	Behavioral bool
}

// karpenterDriver* identify the verifier-owned driver workload and the pool
// template labels it schedules against — fixed so the behavioral scenario's
// template block and this proof cannot drift apart silently.
const (
	karpenterDriverNamespace = "e2e-karpenter-drive"
	karpenterDriverPod       = "e2e-karpenter-driver"
	karpenterPoolLabelKey    = "planton.dev/e2e-karpenter-pool"
	karpenterPoolLabelValue  = "behavioral"
	karpenterFenceTaintKey   = "planton.dev/e2e-karpenter-only"
)

func (v *KarpenterNodePoolVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] karpenter NodePool %q\n", v.Name)

	if err := KubectlResourceExists(ctx, kubeconfig, "nodepools.karpenter.sh", v.Name, ""); err != nil {
		return err
	}
	// A registered pool reaches Ready once its NodeClass resolves — the
	// object-grade contract (the controller accepted the pool).
	if err := kubectlWait(ctx, kubeconfig, "nodepools.karpenter.sh", v.Name, "",
		"condition=Ready", 3*time.Minute); err != nil {
		return errors.Wrap(err, "NodePool never reached Ready (NodeClass resolution or controller registration failed)")
	}

	if !v.Behavioral {
		return nil
	}
	return v.proveNodeLaunch(ctx, kubeconfig)
}

func (v *KarpenterNodePoolVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// Deleting the pool drains and terminates its nodes (Karpenter's own
	// finalizer); pool absence is the API-side contract, and the AWS audit
	// cross-checks instances by tag.
	return KubectlResourceAbsent(ctx, kubeconfig, "nodepools.karpenter.sh", v.Name, "")
}

// proveNodeLaunch runs the provision-and-consolidate loop. Every resource is
// verifier-owned and removed in defers so a failed assertion cannot strand a
// pod that would pin the launched node past the lane.
func (v *KarpenterNodePoolVerifier) proveNodeLaunch(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral node launch: a pending pod must produce a real node\n")

	if err := v.kubectl(ctx, kubeconfig, "create", "namespace", karpenterDriverNamespace); err != nil {
		return errors.Wrap(err, "failed to create driver namespace")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "namespace", karpenterDriverNamespace,
			"--ignore-not-found", "--wait=false")
	}()

	// The driver: schedulable ONLY on this pool's nodes (selector + the
	// fencing taint's toleration), sized well under a t3.medium so a single
	// cheap node satisfies it.
	driver := fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  nodeSelector:
    %s: %s
  tolerations:
    - key: %s
      operator: Equal
      value: "true"
      effect: NoSchedule
  containers:
    - name: driver
      image: registry.k8s.io/pause:3.10
      resources:
        requests:
          cpu: 100m
          memory: 64Mi
`, karpenterDriverPod, karpenterDriverNamespace,
		karpenterPoolLabelKey, karpenterPoolLabelValue, karpenterFenceTaintKey)
	driverFile, err := writeTempManifest(driver)
	if err != nil {
		return err
	}
	defer os.Remove(driverFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", driverFile); err != nil {
		return errors.Wrap(err, "failed to apply driver pod")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "pod", karpenterDriverPod,
			"-n", karpenterDriverNamespace, "--ignore-not-found", "--wait=false")
	}()

	// Node launch: EC2 provisioning + join + kubelet Ready typically lands
	// in 1-2 minutes for AL2023; 6 covers cold AMI paths.
	if err := v.waitForPoolNodes(ctx, kubeconfig, 1, 6*time.Minute); err != nil {
		return errors.Wrap(err, "Karpenter never launched a node for the pending driver pod")
	}
	if err := kubectlWait(ctx, kubeconfig, "pod", karpenterDriverPod, karpenterDriverNamespace,
		"condition=Ready", 3*time.Minute); err != nil {
		return errors.Wrap(err, "driver pod never became Ready on the launched node")
	}
	fmt.Printf("  [verify] node launched and driver scheduled — Karpenter provisioned real capacity\n")

	// The reclaim half: with the driver gone the node is empty, and the
	// pool's consolidateAfter (30s in the scenario) must remove it.
	if err := v.kubectl(ctx, kubeconfig, "delete", "pod", karpenterDriverPod,
		"-n", karpenterDriverNamespace, "--wait=true", "--timeout=2m"); err != nil {
		return errors.Wrap(err, "failed to delete driver pod")
	}
	if err := v.waitForPoolNodes(ctx, kubeconfig, 0, 6*time.Minute); err != nil {
		return errors.Wrap(err, "consolidation never reclaimed the empty node")
	}
	fmt.Printf("  [verify] empty node consolidated away — the full provision/reclaim loop is proven\n")
	return nil
}

// waitForPoolNodes polls until exactly want nodes carry the pool's template
// label (Karpenter stamps template labels onto the nodes it launches).
func (v *KarpenterNodePoolVerifier) waitForPoolNodes(ctx context.Context, kubeconfig string, want int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "nodes", "-l", karpenterPoolLabelKey+"="+karpenterPoolLabelValue,
			"-o", "jsonpath={.items[*].metadata.name}").CombinedOutput()
		names := strings.Fields(strings.TrimSpace(string(out)))
		if err == nil && len(names) == want {
			return nil
		}
		last = fmt.Sprintf("nodes=%v err=%v", names, err)
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("pool node count never reached %d (last: %s)", want, last)
}

func (v *KarpenterNodePoolVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}

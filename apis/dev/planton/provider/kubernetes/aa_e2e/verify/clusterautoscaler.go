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

// ClusterAutoscalerInstallVerifier checks a Cluster Autoscaler installation
// to the point the RECONCILE LOOP is provably running: the deployment
// Available, and the cluster-autoscaler-status ConfigMap present — the
// autoscaler writes that ConfigMap from inside its main loop on every scan,
// so its existence is the loop's own heartbeat, not a proxy. The kind lane
// runs the KWOK simulation arm.
//
// When Behavioral is set (the aws-asg-scaling scenario on the aws-eks
// profile), the verifier additionally proves REAL node-group scaling: a
// verifier-owned Deployment that only the batch's tainted min-0 "ca-scale"
// node group can host goes Pending, the autoscaler must scale the ASG up
// and the pods must schedule onto a real node; after the workload is
// deleted, scale-down (tuned tight via the spec's typed scaling block) must
// drain the group back to zero.
type ClusterAutoscalerInstallVerifier struct {
	Namespace string
	// Behavioral switches on the live ASG scale proof (see above).
	Behavioral bool
}

// caScale* identify the verifier-owned burst workload and the batch node
// group it targets — the labels/taints are the eksctl batch config's
// contract (aa_e2e/realcluster/aws-eks/cluster.eksctl.yaml).
const (
	caScaleNamespace     = "e2e-ca-scale"
	caScaleDeployment    = "e2e-ca-burst"
	caScaleNodeRoleLabel = "planton.dev/e2e-node-role"
	caScaleNodeRole      = "ca-scale"
	caScaleTaintKey      = "planton.dev/e2e-ca-scale"
)

func (v *ClusterAutoscalerInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] cluster-autoscaler installation in namespace %q\n", v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for cluster-autoscaler", v.Namespace)
	}

	// The deployment name embeds the cloud-provider arm: the chart's
	// fullname is <release>-<cloudProvider>-<chartName> (its name template
	// defaults to "<cloudProvider>-cluster-autoscaler", which never equals
	// the release name) — e.g. cluster-autoscaler-kwok-cluster-autoscaler
	// on the kind lane. Resolve it through the release's instance label
	// instead of hardcoding an arm here.
	deployName, err := v.deploymentName(ctx, kubeconfig)
	if err != nil {
		return err
	}
	if err := kubectlWait(ctx, kubeconfig, "deployment", deployName, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrapf(err, "cluster-autoscaler deployment %q not available", deployName)
	}

	// The status ConfigMap appears only after the loop completes a scan —
	// give the freshly-Available pod a moment to write it.
	deadline := time.Now().Add(2 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		if lastErr = KubectlResourceExists(ctx, kubeconfig, "configmap", "cluster-autoscaler-status", v.Namespace); lastErr == nil {
			fmt.Printf("  [verify] cluster-autoscaler-status ConfigMap present — the reconcile loop is running\n")
			if !v.Behavioral {
				return nil
			}
			return v.proveAsgScaling(ctx, kubeconfig)
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(lastErr, "cluster-autoscaler-status ConfigMap never materialized — the loop is not running")
}

// proveAsgScaling runs the scale-up/scale-down cycle against the batch's
// dedicated node group. Everything is verifier-owned and cleaned in defers;
// the final zero-node wait keeps a scaled-up ASG from outliving the lane
// (the cost- and audit-relevant residue).
func (v *ClusterAutoscalerInstallVerifier) proveAsgScaling(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral ASG scaling: a pending burst pod must produce a real node\n")

	if err := v.kubectl(ctx, kubeconfig, "create", "namespace", caScaleNamespace); err != nil {
		return errors.Wrap(err, "failed to create burst namespace")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "namespace", caScaleNamespace,
			"--ignore-not-found", "--wait=false")
	}()

	burst := fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
spec:
  replicas: 1
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      nodeSelector:
        %s: %s
      tolerations:
        - key: %s
          operator: Equal
          value: "true"
          effect: NoSchedule
      containers:
        - name: burst
          image: registry.k8s.io/pause:3.10
          resources:
            requests:
              cpu: 100m
              memory: 64Mi
`, caScaleDeployment, caScaleNamespace, caScaleDeployment, caScaleDeployment,
		caScaleNodeRoleLabel, caScaleNodeRole, caScaleTaintKey)
	burstFile, err := writeTempManifest(burst)
	if err != nil {
		return err
	}
	defer os.Remove(burstFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", burstFile); err != nil {
		return errors.Wrap(err, "failed to apply burst deployment")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "deployment", caScaleDeployment,
			"-n", caScaleNamespace, "--ignore-not-found", "--wait=false")
	}()

	// Scale-up: ASG launch + node join + Ready typically lands in 2-4
	// minutes for a managed node group; 8 covers cold paths.
	if err := v.waitForGroupNodes(ctx, kubeconfig, 1, 8*time.Minute); err != nil {
		return errors.Wrap(err, "the autoscaler never scaled the ca-scale node group up")
	}
	if err := kubectlWait(ctx, kubeconfig, "deployment", caScaleDeployment, caScaleNamespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "burst deployment never became Available on the scaled node")
	}
	fmt.Printf("  [verify] node group scaled up and burst pod scheduled — a real ASG scale-up\n")

	// Scale-down: with the burst gone the node is unneeded; the scenario
	// tunes unneeded_time/delay_after_add to 1m, so ~5 minutes covers the
	// scan cadence + drain + ASG terminate.
	if err := v.kubectl(ctx, kubeconfig, "delete", "deployment", caScaleDeployment,
		"-n", caScaleNamespace, "--wait=true", "--timeout=2m"); err != nil {
		return errors.Wrap(err, "failed to delete burst deployment")
	}
	if err := v.waitForGroupNodes(ctx, kubeconfig, 0, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the autoscaler never scaled the ca-scale node group back down")
	}
	fmt.Printf("  [verify] node group drained back to zero — the full scale-up/scale-down loop is proven\n")
	return nil
}

// waitForGroupNodes polls until exactly want READY nodes carry the ca-scale
// node-group label.
func (v *ClusterAutoscalerInstallVerifier) waitForGroupNodes(ctx context.Context, kubeconfig string, want int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "nodes", "-l", caScaleNodeRoleLabel+"="+caScaleNodeRole,
			"-o", "jsonpath={.items[*].metadata.name}").CombinedOutput()
		names := strings.Fields(strings.TrimSpace(string(out)))
		if err == nil && len(names) == want {
			return nil
		}
		last = fmt.Sprintf("nodes=%v err=%v", names, err)
		time.Sleep(15 * time.Second)
	}
	return errors.Errorf("ca-scale node count never reached %d (last: %s)", want, last)
}

func (v *ClusterAutoscalerInstallVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}

func (v *ClusterAutoscalerInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// Absence is asserted through the same instance label the existence
	// path resolves the deployment by — no arm-specific name to hardcode.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", "-n", v.Namespace,
		"-l", "app.kubernetes.io/instance=cluster-autoscaler",
		"-o", "jsonpath={.items[*].metadata.name}").CombinedOutput()
	if err != nil {
		// The namespace itself may already be gone — that IS absence.
		return nil
	}
	if names := strings.TrimSpace(string(out)); names != "" {
		return errors.Errorf("cluster-autoscaler deployment(s) still present: %s", names)
	}
	return nil
}

// deploymentName resolves the release's deployment through the Helm
// instance label (app.kubernetes.io/instance = the fixed release name),
// polling briefly because verification can race the deployment's own
// creation.
func (v *ClusterAutoscalerInstallVerifier) deploymentName(ctx context.Context, kubeconfig string) (string, error) {
	deadline := time.Now().Add(2 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "deployment", "-n", v.Namespace,
			"-l", "app.kubernetes.io/instance=cluster-autoscaler",
			"-o", "jsonpath={.items[0].metadata.name}").CombinedOutput()
		name := strings.TrimSpace(string(out))
		if err == nil && name != "" {
			return name, nil
		}
		last = fmt.Sprintf("out=%q err=%v", name, err)
		time.Sleep(5 * time.Second)
	}
	return "", errors.Errorf("no deployment carrying app.kubernetes.io/instance=cluster-autoscaler appeared in %q (last: %s)", v.Namespace, last)
}

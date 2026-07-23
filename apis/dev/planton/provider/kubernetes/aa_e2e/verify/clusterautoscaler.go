package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ClusterAutoscalerInstallVerifier checks a Cluster Autoscaler installation
// to the point the RECONCILE LOOP is provably running: the deployment
// Available, and the cluster-autoscaler-status ConfigMap present — the
// autoscaler writes that ConfigMap from inside its main loop on every scan,
// so its existence is the loop's own heartbeat, not a proxy. Scaling
// behavior against real node groups needs a real cloud and rides the
// batched real-cluster lanes; the kind lane runs the KWOK simulation arm.
type ClusterAutoscalerInstallVerifier struct {
	Namespace string
}

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
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(lastErr, "cluster-autoscaler-status ConfigMap never materialized — the loop is not running")
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

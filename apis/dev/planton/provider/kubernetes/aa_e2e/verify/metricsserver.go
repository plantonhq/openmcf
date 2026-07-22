package verify

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

// MetricsServerInstallVerifier checks a metrics-server installation all the
// way to METRIC FLOW: the Deployment is Available, the cluster-wide
// v1beta1.metrics.k8s.io APIService reports Available (the aggregation
// layer reaches the service), and `kubectl top nodes` actually returns
// values — the customer-grade contract, since an installed-but-not-scraping
// metrics-server is exactly the failure mode users hit (kubelet TLS
// rejection) and the reason HPAs mysteriously never scale.
type MetricsServerInstallVerifier struct {
	Namespace string
}

func (v *MetricsServerInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] metrics-server installation in namespace %q\n", v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for metrics-server", v.Namespace)
	}

	// The release name (and therefore the Deployment name) is fixed —
	// one installation per cluster.
	if err := kubectlWait(ctx, kubeconfig, "deployment", "metrics-server", v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrapf(err, "metrics-server deployment not available in namespace %q", v.Namespace)
	}

	if err := kubectlWait(ctx, kubeconfig, "apiservice", "v1beta1.metrics.k8s.io", "",
		"condition=Available", 2*time.Minute); err != nil {
		return errors.Wrap(err, "v1beta1.metrics.k8s.io APIService not available")
	}

	// Metric flow: kubectl top only succeeds once the first kubelet scrape
	// landed and the API serves real values.
	deadline := time.Now().Add(2 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "top", "nodes")
		if out, err := cmd.CombinedOutput(); err == nil {
			fmt.Printf("  [verify] kubectl top nodes returns values — metrics are flowing\n")
			return nil
		} else {
			lastErr = errors.Errorf("kubectl top nodes: %v: %s", err, string(out))
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Wrap(lastErr, "metrics never started flowing (kubectl top nodes kept failing)")
}

func (v *MetricsServerInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", "metrics-server", v.Namespace); err != nil {
		return err
	}
	// The APIService is cluster-scoped — it would orphan silently (and
	// break the aggregation layer with a dangling backend) if the release
	// left it behind.
	return KubectlResourceAbsent(ctx, kubeconfig, "apiservice", "v1beta1.metrics.k8s.io", "")
}

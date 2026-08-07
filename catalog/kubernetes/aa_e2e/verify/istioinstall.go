package verify

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

// IstioInstallVerifier checks an Istio control-plane installation to the
// customer-grade contract: istiod is Available (its webhooks and discovery
// service gate every mesh-config apply and every injection), the Istio CRDs
// are Established, and — in ambient mode — the istio-cni node agent and the
// ztunnel per-node proxy DaemonSets have all their pods ready (an ambient
// mesh whose node dataplane is not running silently leaves every enrolled
// pod unmeshed).
type IstioInstallVerifier struct {
	Namespace string
	// IstiodName is the istiod Deployment/Service name: "istiod" for the
	// default revision, "istiod-<revision>" for a named revision.
	IstiodName string
	// Ambient is true when the scenario installs the ambient data plane
	// (spec.dataplane_mode == "ambient").
	Ambient bool
}

func (v *IstioInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] istio control plane (%s) in namespace %q (ambient=%v)\n",
		v.IstiodName, v.Namespace, v.Ambient)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for istio", v.Namespace)
	}

	if err := kubectlWait(ctx, kubeconfig, "deployment", v.IstiodName, v.Namespace,
		"condition=Available", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "istiod deployment %q not available in namespace %q", v.IstiodName, v.Namespace)
	}

	// The base release's CRDs are the precondition every typed Istio kind
	// (DestinationRule, AuthorizationPolicy, ...) applies against.
	for _, crd := range []string{
		"destinationrules.networking.istio.io",
		"authorizationpolicies.security.istio.io",
		"telemetries.telemetry.istio.io",
	} {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", 2*time.Minute); err != nil {
			return errors.Wrapf(err, "istio CRD %q not established", crd)
		}
	}

	if v.Ambient {
		for _, ds := range []string{"istio-cni-node", "ztunnel"} {
			if err := kubectlDaemonSetReady(ctx, kubeconfig, ds, v.Namespace, 5*time.Minute); err != nil {
				return errors.Wrapf(err, "ambient dataplane DaemonSet %q not ready in namespace %q", ds, v.Namespace)
			}
		}
	}

	return nil
}

// kubectlDaemonSetReady polls `kubectl rollout status daemonset/<name>`
// until it reports success — the readiness contract for per-node agents
// (kubectl wait has no DaemonSet-level Ready condition to wait on).
func kubectlDaemonSetReady(ctx context.Context, kubeconfig, name, namespace string, timeout time.Duration) error {
	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"rollout", "status", "daemonset/"+name,
		"-n", namespace,
		fmt.Sprintf("--timeout=%s", timeout))
	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Errorf("rollout status daemonset/%s: %v: %s", name, err, string(out))
	}
	return nil
}

func (v *IstioInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.IstiodName, v.Namespace); err != nil {
		return err
	}
	// The CRDs are module-owned (applied outside the Helm release), so
	// destroy removes them — a leftover CRD is an orphan, not keep policy.
	if err := KubectlResourceAbsent(ctx, kubeconfig, "crd", "destinationrules.networking.istio.io", ""); err != nil {
		return err
	}
	if v.Ambient {
		for _, ds := range []string{"istio-cni-node", "ztunnel"} {
			if err := KubectlResourceAbsent(ctx, kubeconfig, "daemonset", ds, v.Namespace); err != nil {
				return err
			}
		}
	}
	return nil
}

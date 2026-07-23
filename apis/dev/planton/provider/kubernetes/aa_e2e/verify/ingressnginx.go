package verify

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// IngressNginxInstallVerifier checks an ingress-nginx installation: the
// controller Deployment is Available and the instance's IngressClass
// exists. The module pins the chart's fullname to metadata.name, so the
// controller Deployment is deterministically "<name>-controller" per
// instance (multiple controllers per cluster are first-class for this
// component).
//
// A cloud load-balancer address is asserted only when LbAddress is set (the
// aws-nlb scenario on the aws-eks profile): on clusters without a cloud LB
// controller (the kind lane) scenarios use node_port, so an Available
// controller owning its IngressClass is the kind-cluster contract, and the
// address itself is the real-cluster lane's proof — validating the spec's
// documented per-cloud annotation recipe as written.
type IngressNginxInstallVerifier struct {
	Namespace        string
	ComponentName    string
	IngressClassName string
	// LbAddress additionally asserts the controller Service received a real
	// cloud load-balancer address (see above).
	LbAddress bool
}

func (v *IngressNginxInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] ingress-nginx installation %q in namespace %q (class %q)\n",
		v.ComponentName, v.Namespace, v.IngressClassName)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for ingress-nginx %q", v.Namespace, v.ComponentName)
	}

	controllerDeployment := v.ComponentName + "-controller"
	if err := kubectlWait(ctx, kubeconfig, "deployment", controllerDeployment, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrapf(err, "controller deployment %q not available in namespace %q", controllerDeployment, v.Namespace)
	}

	// The IngressClass is the instance's routing identity — cluster-scoped.
	if err := KubectlResourceExists(ctx, kubeconfig, "ingressclass", v.IngressClassName, ""); err != nil {
		return errors.Wrapf(err, "ingressclass %q not found for ingress-nginx %q", v.IngressClassName, v.ComponentName)
	}

	if v.LbAddress {
		// The controller Service is "<fullname>-controller" (fullname
		// pinned to metadata.name — same derivation as the Deployment).
		address, err := waitForServiceLbAddress(ctx, kubeconfig, v.ComponentName+"-controller", v.Namespace, 6*time.Minute)
		if err != nil {
			return errors.Wrap(err, "controller service never received a cloud LB address — the documented annotation recipe failed")
		}
		fmt.Printf("  [verify] cloud LB provisioned for the controller: %s\n", address)
	}

	return nil
}

func (v *IngressNginxInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// Multiple instances can share a namespace, so the release's own
	// Deployment (not the namespace) is the absence signal — plus the
	// cluster-scoped IngressClass, which would otherwise orphan silently.
	controllerDeployment := v.ComponentName + "-controller"
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", controllerDeployment, v.Namespace); err != nil {
		return err
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "ingressclass", v.IngressClassName, "")
}

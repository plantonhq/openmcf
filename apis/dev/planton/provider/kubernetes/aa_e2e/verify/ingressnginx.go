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
// Deliberately NOT asserted here: a cloud load-balancer address. On
// clusters without a cloud LB controller (the kind lane) scenarios use
// node_port, and real LB provisioning rides the batched real-cluster
// lanes. An Available controller owning its IngressClass is the
// kind-cluster contract.
type IngressNginxInstallVerifier struct {
	Namespace        string
	ComponentName    string
	IngressClassName string
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

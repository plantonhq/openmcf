package component

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

// certManagerCertificateCRD is how cert-manager's presence is detected before
// asking it to issue a certificate that would otherwise never appear.
const certManagerCertificateCRD = "certificates.cert-manager.io"

// Ingress publishes Planton at a browser-reachable URL through the cluster's
// ingress controller: the console on "/" and the control plane's gRPC-Web
// endpoint on the API path prefix, one hostname for both (same origin).
//
// It has no dependencies so the Ingress is admitted -- and, when no hostname
// is configured, the public URL is derived from the controller's published
// address -- in parallel with the application's (slow) first boot. Backends
// referencing not-yet-ready Services are fine: the controller serves 503s
// until they come up.
type Ingress struct{ Base }

func (i *Ingress) Name() string                                { return "ingress" }
func (i *Ingress) Dependencies(_ *v1.PlantonPlatform) []string { return nil }

// IsEnabled reports whether external access is turned on. Without it the
// platform stays ClusterIP-only (kubectl port-forward access).
func (i *Ingress) IsEnabled(planton *v1.PlantonPlatform) bool {
	return isIngressEnabled(planton)
}

// Reconcile renders the front door the spec names -- an Ingress object for an
// Ingress controller, or an HTTPRoute attached to a Gateway API Gateway --
// from the one route table in resources.FrontDoorRoutes, publishes the public
// URL to status the moment it is known, and reports readiness only when the
// door actually serves the platform (a certificate still being issued, a
// route the Gateway has not accepted, are "not yet", not Ready).
func (i *Ingress) Reconcile(ctx context.Context, c client.Client, _ *runtime.Scheme, planton *v1.PlantonPlatform) (Result, error) {
	if planton.Spec.Ingress.GatewayRef != nil {
		return i.reconcileGatewayEdge(ctx, c, planton)
	}
	return i.reconcileIngressEdge(ctx, c, planton)
}

// preflightTLS verifies the certificate path can actually work before the
// Ingress advertises HTTPS: a BYO Secret must exist, and asking cert-manager
// for a certificate requires cert-manager to be installed.
func (i *Ingress) preflightTLS(ctx context.Context, c client.Client, planton *v1.PlantonPlatform) (string, error) {
	tls := planton.Spec.Ingress.TLS
	if tls == nil {
		return "", nil
	}

	if tls.SecretName != "" {
		var secret corev1.Secret
		err := c.Get(ctx, types.NamespacedName{Name: tls.SecretName, Namespace: planton.Namespace}, &secret)
		if apierrors.IsNotFound(err) {
			return fmt.Sprintf(
				"TLS Secret %q not found in namespace %q; create it (type kubernetes.io/tls) or switch to spec.ingress.tls.issuer",
				tls.SecretName, planton.Namespace), nil
		}
		if err != nil {
			return "", fmt.Errorf("checking TLS Secret %s: %w", tls.SecretName, err)
		}
	}

	if tls.Issuer != nil {
		installed, err := i.IsCRDInstalled(ctx, c, certManagerCertificateCRD)
		if err != nil {
			return "", fmt.Errorf("checking for cert-manager: %w", err)
		}
		if !installed {
			return "cert-manager is not installed (Certificate CRD missing); install cert-manager or bring a certificate via spec.ingress.tls.secretName", nil
		}
	}

	return "", nil
}

func isIngressEnabled(planton *v1.PlantonPlatform) bool {
	return planton.Spec.Ingress != nil && planton.Spec.Ingress.Enabled
}

// RBAC markers for the resources the Ingress component manages: it owns the
// platform's Ingress object, reads IngressClasses to validate the class and
// detect the serving controller, and reads cert-manager Certificates to gate
// readiness on actual issuance (read-only: cert-manager owns the objects,
// the operator only reports their truth).
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingressclasses,verbs=get;list;watch
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch

package component

import (
	"fmt"

	v1 "github.com/plantonhq/planton/operator/api/v1"
	"github.com/plantonhq/planton/operator/internal/resources"
)

// frontDoorURL resolves the browser-facing URL of the platform's ONE front
// door. Every URL-bearing surface derives from it -- Keycloak's advertised
// issuer and redirect URIs, the console's API endpoint and auth callbacks,
// status.consoleUrl -- so the two access modes stay the same architecture at
// different addresses:
//
//   - ingress enabled:  the public URL (explicit hostname, or the
//     auto-derived hostname once the ingress component resolves it)
//   - ingress disabled: the deterministic localhost URL the gateway
//     component serves over kubectl port-forward
//
// resolved is false only in the auto-hostname window where the ingress URL is
// not yet known; the gateway URL is always known.
func frontDoorURL(planton *v1.PlantonPlatform) (url string, resolved bool) {
	if !isIngressEnabled(planton) {
		return fmt.Sprintf("http://localhost:%d", gatewayLocalPort(planton)), true
	}
	spec := planton.Spec.Ingress
	if spec.Hostname != "" {
		// Explicit hostname: the URL is known statically -- render the
		// public endpoints immediately, independent of ingress health.
		return resources.PublicURL(spec.Hostname, spec.TLS != nil), true
	}
	// Auto-derived hostname, resolved by the Ingress component (which
	// reconciles earlier in the same pass) and published in status.
	if planton.Status.ConsoleURL != "" {
		return planton.Status.ConsoleURL, true
	}
	return "", false
}

// gatewayLocalPort returns the workstation port sign-in URLs are pinned to in
// gateway mode (spec.gateway.localPort, or the default).
func gatewayLocalPort(planton *v1.PlantonPlatform) int32 {
	if planton.Spec.Gateway != nil && planton.Spec.Gateway.LocalPort != nil {
		return *planton.Spec.Gateway.LocalPort
	}
	return resources.GatewayDefaultLocalPort
}

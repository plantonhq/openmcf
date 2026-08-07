package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrusttunnelroutev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrusttunnelroute/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	CloudflareProviderConfig       *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustTunnelRoute *cloudflarezerotrusttunnelroutev1alpha1.CloudflareZeroTrustTunnelRoute
}

func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrusttunnelroutev1alpha1.CloudflareZeroTrustTunnelRouteStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustTunnelRoute = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

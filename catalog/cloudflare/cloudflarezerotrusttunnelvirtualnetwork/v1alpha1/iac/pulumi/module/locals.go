package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrusttunnelvirtualnetworkv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrusttunnelvirtualnetwork/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	CloudflareProviderConfig                *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustTunnelVirtualNetwork *cloudflarezerotrusttunnelvirtualnetworkv1alpha1.CloudflareZeroTrustTunnelVirtualNetwork
}

func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrusttunnelvirtualnetworkv1alpha1.CloudflareZeroTrustTunnelVirtualNetworkStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustTunnelVirtualNetwork = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

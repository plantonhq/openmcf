package module

import (
	cloudflareprovider "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare"
	cloudflarecustomhostnamefallbackoriginv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare/cloudflarecustomhostnamefallbackorigin/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	CloudflareProviderConfig               *cloudflareprovider.CloudflareProviderConfig
	CloudflareCustomHostnameFallbackOrigin *cloudflarecustomhostnamefallbackoriginv1alpha1.CloudflareCustomHostnameFallbackOrigin
}

func initializeLocals(_ *pulumi.Context, stackInput *cloudflarecustomhostnamefallbackoriginv1alpha1.CloudflareCustomHostnameFallbackOriginStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareCustomHostnameFallbackOrigin = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

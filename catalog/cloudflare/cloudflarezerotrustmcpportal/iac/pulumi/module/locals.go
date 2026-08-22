package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustmcpportalv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustmcpportal/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig     *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustMcpPortal *cloudflarezerotrustmcpportalv1alpha1.CloudflareZeroTrustMcpPortal
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustmcpportalv1alpha1.CloudflareZeroTrustMcpPortalStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustMcpPortal = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustaccessservicetokenv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustaccessservicetoken/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig              *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustAccessServiceToken *cloudflarezerotrustaccessservicetokenv1alpha1.CloudflareZeroTrustAccessServiceToken
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustaccessservicetokenv1alpha1.CloudflareZeroTrustAccessServiceTokenStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustAccessServiceToken = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

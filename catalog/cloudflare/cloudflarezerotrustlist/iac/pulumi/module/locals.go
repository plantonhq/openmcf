package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustlistv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustlist/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustList  *cloudflarezerotrustlistv1alpha1.CloudflareZeroTrustList
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustlistv1alpha1.CloudflareZeroTrustListStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustList = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

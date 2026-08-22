package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustorganizationv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustorganization/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig        *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustOrganization *cloudflarezerotrustorganizationv1alpha1.CloudflareZeroTrustOrganization
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustorganizationv1alpha1.CloudflareZeroTrustOrganizationStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustOrganization = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

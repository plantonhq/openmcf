package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustaccessinfrastructuretargetv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustaccessinfrastructuretarget/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig                      *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustAccessInfrastructureTarget *cloudflarezerotrustaccessinfrastructuretargetv1alpha1.CloudflareZeroTrustAccessInfrastructureTarget
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustaccessinfrastructuretargetv1alpha1.CloudflareZeroTrustAccessInfrastructureTargetStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustAccessInfrastructureTarget = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

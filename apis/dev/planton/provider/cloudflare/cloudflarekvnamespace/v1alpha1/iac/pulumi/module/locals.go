package module

import (
	cloudflareprovider "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare"
	cloudflarekvnamespacev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare/cloudflarekvnamespace/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareKvNamespace    *cloudflarekvnamespacev1alpha1.CloudflareKvNamespace
}

// initializeLocals copies stack‑input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarekvnamespacev1alpha1.CloudflareKvNamespaceStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareKvNamespace = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

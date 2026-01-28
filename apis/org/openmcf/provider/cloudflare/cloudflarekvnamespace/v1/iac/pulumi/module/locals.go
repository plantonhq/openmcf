package module

import (
	cloudflareprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare"
	cloudflarekvnamespacev1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare/cloudflarekvnamespace/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareKvNamespace    *cloudflarekvnamespacev1.CloudflareKvNamespace
}

// initializeLocals copies stack‑input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarekvnamespacev1.CloudflareKvNamespaceStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareKvNamespace = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

package module

import (
	cloudflareprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare"
	cloudflarednszonev1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/cloudflare/cloudflarednszone/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles the bits we need everywhere else.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareDnsZone        *cloudflarednszonev1.CloudflareDnsZone
}

// initializeLocals copies fields from the stack‑input into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarednszonev1.CloudflareDnsZoneStackInput) *Locals {
	locals := &Locals{}

	locals.CloudflareDnsZone = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig

	return locals
}

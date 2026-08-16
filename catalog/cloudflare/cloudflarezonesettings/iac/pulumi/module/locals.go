package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezonesettingsv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezonesettings/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareZoneSettings   *cloudflarezonesettingsv1alpha1.CloudflareZoneSettings
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezonesettingsv1alpha1.CloudflareZoneSettingsStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZoneSettings = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

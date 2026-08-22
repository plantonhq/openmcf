package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarecachesettingsv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarecachesettings/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareCacheSettings  *cloudflarecachesettingsv1alpha1.CloudflareCacheSettings
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarecachesettingsv1alpha1.CloudflareCacheSettingsStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareCacheSettings = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

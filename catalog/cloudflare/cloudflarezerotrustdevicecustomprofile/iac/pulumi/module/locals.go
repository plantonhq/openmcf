package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustdevicecustomprofilev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustdevicecustomprofile/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig               *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustDeviceCustomProfile *cloudflarezerotrustdevicecustomprofilev1alpha1.CloudflareZeroTrustDeviceCustomProfile
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustdevicecustomprofilev1alpha1.CloudflareZeroTrustDeviceCustomProfileStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustDeviceCustomProfile = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

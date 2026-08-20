package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustdeviceposturerulev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustdeviceposturerule/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig             *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustDevicePostureRule *cloudflarezerotrustdeviceposturerulev1alpha1.CloudflareZeroTrustDevicePostureRule
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustdeviceposturerulev1alpha1.CloudflareZeroTrustDevicePostureRuleStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustDevicePostureRule = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

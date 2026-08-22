package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarebotmanagementv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarebotmanagement/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareBotManagement  *cloudflarebotmanagementv1alpha1.CloudflareBotManagement
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarebotmanagementv1alpha1.CloudflareBotManagementStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareBotManagement = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

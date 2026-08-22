package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflareaccountapitokenv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareaccountapitoken/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig  *cloudflareprovider.CloudflareProviderConfig
	CloudflareAccountApiToken *cloudflareaccountapitokenv1alpha1.CloudflareAccountApiToken
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareaccountapitokenv1alpha1.CloudflareAccountApiTokenStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareAccountApiToken = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

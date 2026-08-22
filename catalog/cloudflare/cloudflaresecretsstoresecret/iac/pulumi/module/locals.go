package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflaresecretsstoresecretv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflaresecretsstoresecret/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig     *cloudflareprovider.CloudflareProviderConfig
	CloudflareSecretsStoreSecret *cloudflaresecretsstoresecretv1alpha1.CloudflareSecretsStoreSecret
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflaresecretsstoresecretv1alpha1.CloudflareSecretsStoreSecretStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareSecretsStoreSecret = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

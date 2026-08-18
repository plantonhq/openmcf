package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflaresecretsstorev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflaresecretsstore/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareSecretsStore   *cloudflaresecretsstorev1alpha1.CloudflareSecretsStore
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflaresecretsstorev1alpha1.CloudflareSecretsStoreStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareSecretsStore = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

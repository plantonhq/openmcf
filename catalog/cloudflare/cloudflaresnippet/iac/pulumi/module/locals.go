package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflaresnippetv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflaresnippet/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareSnippet        *cloudflaresnippetv1alpha1.CloudflareSnippet
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflaresnippetv1alpha1.CloudflareSnippetStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareSnippet = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflareauthenticatedoriginpullsv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareauthenticatedoriginpulls/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig           *cloudflareprovider.CloudflareProviderConfig
	CloudflareAuthenticatedOriginPulls *cloudflareauthenticatedoriginpullsv1alpha1.CloudflareAuthenticatedOriginPulls
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareauthenticatedoriginpullsv1alpha1.CloudflareAuthenticatedOriginPullsStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareAuthenticatedOriginPulls = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

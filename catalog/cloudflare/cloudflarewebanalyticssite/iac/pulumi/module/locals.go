package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarewebanalyticssitev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarewebanalyticssite/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig   *cloudflareprovider.CloudflareProviderConfig
	CloudflareWebAnalyticsSite *cloudflarewebanalyticssitev1alpha1.CloudflareWebAnalyticsSite
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarewebanalyticssitev1alpha1.CloudflareWebAnalyticsSiteStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareWebAnalyticsSite = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

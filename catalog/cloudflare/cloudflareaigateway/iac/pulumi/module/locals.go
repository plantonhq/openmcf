package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflareaigatewayv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareaigateway/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareAiGateway      *cloudflareaigatewayv1alpha1.CloudflareAiGateway
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareaigatewayv1alpha1.CloudflareAiGatewayStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareAiGateway = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

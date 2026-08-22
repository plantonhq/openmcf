package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarezerotrustgatewaypolicyv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustgatewaypolicy/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig         *cloudflareprovider.CloudflareProviderConfig
	CloudflareZeroTrustGatewayPolicy *cloudflarezerotrustgatewaypolicyv1alpha1.CloudflareZeroTrustGatewayPolicy
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarezerotrustgatewaypolicyv1alpha1.CloudflareZeroTrustGatewayPolicyStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareZeroTrustGatewayPolicy = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

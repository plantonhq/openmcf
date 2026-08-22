package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflareipaccessrulev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareipaccessrule/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareIpAccessRule   *cloudflareipaccessrulev1alpha1.CloudflareIpAccessRule
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareipaccessrulev1alpha1.CloudflareIpAccessRuleStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareIpAccessRule = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

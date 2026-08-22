package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflareworkflowv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareworkflow/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareWorkflow       *cloudflareworkflowv1alpha1.CloudflareWorkflow
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareworkflowv1alpha1.CloudflareWorkflowStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareWorkflow = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

package module

import (
	cloudflareprovider "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare"
	cloudflarelistv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare/cloudflarelist/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareList           *cloudflarelistv1alpha1.CloudflareList
}

func initializeLocals(_ *pulumi.Context, stackInput *cloudflarelistv1alpha1.CloudflareListStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareList = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

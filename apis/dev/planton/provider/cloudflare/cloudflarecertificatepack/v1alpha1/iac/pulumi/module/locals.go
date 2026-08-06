package module

import (
	cloudflareprovider "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare"
	cloudflarecertificatepackv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare/cloudflarecertificatepack/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	CloudflareProviderConfig  *cloudflareprovider.CloudflareProviderConfig
	CloudflareCertificatePack *cloudflarecertificatepackv1alpha1.CloudflareCertificatePack
}

func initializeLocals(_ *pulumi.Context, stackInput *cloudflarecertificatepackv1alpha1.CloudflareCertificatePackStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareCertificatePack = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

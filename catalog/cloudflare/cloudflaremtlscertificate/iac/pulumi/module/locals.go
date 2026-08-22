package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflaremtlscertificatev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflaremtlscertificate/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig  *cloudflareprovider.CloudflareProviderConfig
	CloudflareMtlsCertificate *cloudflaremtlscertificatev1alpha1.CloudflareMtlsCertificate
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflaremtlscertificatev1alpha1.CloudflareMtlsCertificateStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareMtlsCertificate = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

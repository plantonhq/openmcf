package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarecustomsslcertificatev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarecustomsslcertificate/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig       *cloudflareprovider.CloudflareProviderConfig
	CloudflareCustomSslCertificate *cloudflarecustomsslcertificatev1alpha1.CloudflareCustomSslCertificate
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarecustomsslcertificatev1alpha1.CloudflareCustomSslCertificateStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareCustomSslCertificate = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflareauthenticatedoriginpullscertificatev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareauthenticatedoriginpullscertificate/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig                      *cloudflareprovider.CloudflareProviderConfig
	CloudflareAuthenticatedOriginPullsCertificate *cloudflareauthenticatedoriginpullscertificatev1alpha1.CloudflareAuthenticatedOriginPullsCertificate
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareauthenticatedoriginpullscertificatev1alpha1.CloudflareAuthenticatedOriginPullsCertificateStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareAuthenticatedOriginPullsCertificate = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

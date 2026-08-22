package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarewaitingroomv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarewaitingroom/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareWaitingRoom    *cloudflarewaitingroomv1alpha1.CloudflareWaitingRoom
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarewaitingroomv1alpha1.CloudflareWaitingRoomStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareWaitingRoom = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

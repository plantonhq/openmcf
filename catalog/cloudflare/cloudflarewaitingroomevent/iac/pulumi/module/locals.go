package module

import (
	cloudflareprovider "github.com/plantonhq/planton/catalog/cloudflare"
	cloudflarewaitingroomeventv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarewaitingroomevent/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references used across the module.
type Locals struct {
	CloudflareProviderConfig   *cloudflareprovider.CloudflareProviderConfig
	CloudflareWaitingRoomEvent *cloudflarewaitingroomeventv1alpha1.CloudflareWaitingRoomEvent
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflarewaitingroomeventv1alpha1.CloudflareWaitingRoomEventStackInput) *Locals {
	locals := &Locals{}
	locals.CloudflareWaitingRoomEvent = stackInput.Target
	locals.CloudflareProviderConfig = stackInput.ProviderConfig
	return locals
}

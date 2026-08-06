package module

import (
	gcpserviceaccountv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpserviceaccount/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds the resource structure for the GcpServiceAccount component
// as well as any auxiliary fields we might need in the Pulumi module.
type Locals struct {
	GcpServiceAccount *gcpserviceaccountv1alpha1.GcpServiceAccount
}

// initializeLocals creates and returns a Locals struct.
func initializeLocals(ctx *pulumi.Context, stackInput *gcpserviceaccountv1alpha1.GcpServiceAccountStackInput) *Locals {
	locals := &Locals{
		GcpServiceAccount: stackInput.Target,
	}
	return locals
}

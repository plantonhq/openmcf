package module

import (
	gcpprojectiammemberv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpprojectiammember/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpProjectIamMember *gcpprojectiammemberv1alpha1.GcpProjectIamMember
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpprojectiammemberv1alpha1.GcpProjectIamMemberStackInput) *Locals {
	return &Locals{
		GcpProjectIamMember: stackInput.Target,
	}
}

package module

import (
	gcpserviceaccountiammemberv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpserviceaccountiammember/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpServiceAccountIamMember *gcpserviceaccountiammemberv1.GcpServiceAccountIamMember
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpserviceaccountiammemberv1.GcpServiceAccountIamMemberStackInput) *Locals {
	return &Locals{
		GcpServiceAccountIamMember: stackInput.Target,
	}
}

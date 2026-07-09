package module

import (
	gcpkmskeyiammemberv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpkmskeyiammember/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpKmsKeyIamMember *gcpkmskeyiammemberv1.GcpKmsKeyIamMember
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpkmskeyiammemberv1.GcpKmsKeyIamMemberStackInput) *Locals {
	return &Locals{
		GcpKmsKeyIamMember: stackInput.Target,
	}
}

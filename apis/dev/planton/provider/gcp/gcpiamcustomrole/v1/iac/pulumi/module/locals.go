package module

import (
	gcpiamcustomrolev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpiamcustomrole/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpIamCustomRole *gcpiamcustomrolev1.GcpIamCustomRole
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpiamcustomrolev1.GcpIamCustomRoleStackInput) *Locals {
	return &Locals{
		GcpIamCustomRole: stackInput.Target,
	}
}

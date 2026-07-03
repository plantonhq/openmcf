package module

import (
	gcpworkloadidentitypoolv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpworkloadidentitypool/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpWorkloadIdentityPool *gcpworkloadidentitypoolv1.GcpWorkloadIdentityPool
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpworkloadidentitypoolv1.GcpWorkloadIdentityPoolStackInput) *Locals {
	return &Locals{
		GcpWorkloadIdentityPool: stackInput.Target,
	}
}

package module

import (
	gcpsubnetworkv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpsubnetwork/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs. Subnetworks accept no
// labels in GCP, so none are derived here.
type Locals struct {
	GcpSubnetwork *gcpsubnetworkv1.GcpSubnetwork
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpsubnetworkv1.GcpSubnetworkStackInput) *Locals {
	return &Locals{
		GcpSubnetwork: stackInput.Target,
	}
}

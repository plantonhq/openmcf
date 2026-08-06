package module

import (
	gcpregionnetworkendpointgroupv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpregionnetworkendpointgroup/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpRegionNetworkEndpointGroup *gcpregionnetworkendpointgroupv1alpha1.GcpRegionNetworkEndpointGroup

	// The cloud-side name defaults to metadata.name when the spec leaves
	// network_endpoint_group_name empty — the same naming basis every kind uses.
	NetworkEndpointGroupName string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpregionnetworkendpointgroupv1alpha1.GcpRegionNetworkEndpointGroupStackInput) *Locals {
	target := stackInput.Target

	networkEndpointGroupName := target.Spec.NetworkEndpointGroupName
	if networkEndpointGroupName == "" {
		networkEndpointGroupName = target.Metadata.Name
	}

	return &Locals{
		GcpRegionNetworkEndpointGroup: target,
		NetworkEndpointGroupName:      networkEndpointGroupName,
	}
}

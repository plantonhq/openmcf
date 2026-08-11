package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// vpc enables the Compute API and creates the VPC network.
//
// Optional spec fields follow the omit-when-empty convention shared with the
// Terraform module: an unset field is left off the request entirely so the
// GCP API applies its own default (mtu 1460, AFTER_CLASSIC_FIREWALL, ...),
// rather than sending a zero value the API might reject or diff against.
func vpc(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) (*compute.Network, error) {
	// The Compute API must be enabled before a network can be created in a
	// fresh project. disable_on_destroy=false leaves the API on at teardown
	// so other resources in the project keep working.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if locals.GcpVpcNetwork.Spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(locals.GcpVpcNetwork.Spec.ProjectId.GetValue())
	}
	createdComputeService, err := projects.NewService(ctx,
		"compute-api", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable compute api")
	}

	spec := locals.GcpVpcNetwork.Spec
	networkArgs := &compute.NetworkArgs{
		AutoCreateSubnetworks: pulumi.BoolPtr(spec.AutoCreateSubnetworks),
		Name:                  pulumi.String(spec.NetworkName),
	}
	if spec.ProjectId.GetValue() != "" {
		networkArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Sent explicitly on every apply (the attribute is Optional+Computed —
	// the same class as always_compare_med below): omitting it on a
	// GLOBAL→REGIONAL transition would silently keep GLOBAL on the live
	// network, and the live API rejects a create whose routingConfig
	// carries any BGP best-path field without routingMode ("Required field
	// 'resource.routingConfig.routingMode' not specified").
	networkArgs.RoutingMode = pulumi.String(spec.GetRoutingMode().String())

	if spec.Description != "" {
		networkArgs.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.Mtu != nil {
		networkArgs.Mtu = pulumi.IntPtr(int(*spec.Mtu))
	}
	if spec.EnableUlaInternalIpv6 {
		networkArgs.EnableUlaInternalIpv6 = pulumi.BoolPtr(true)
	}
	if spec.InternalIpv6Range != "" {
		networkArgs.InternalIpv6Range = pulumi.StringPtr(spec.InternalIpv6Range)
	}
	if spec.NetworkFirewallPolicyEnforcementOrder != "" {
		networkArgs.NetworkFirewallPolicyEnforcementOrder = pulumi.StringPtr(spec.NetworkFirewallPolicyEnforcementOrder)
	}
	if spec.NetworkProfile != "" {
		networkArgs.NetworkProfile = pulumi.StringPtr(spec.NetworkProfile)
	}
	if spec.DeleteDefaultRoutesOnCreate {
		networkArgs.DeleteDefaultRoutesOnCreate = pulumi.BoolPtr(true)
	}
	if spec.BgpBestPathSelection != nil {
		bgp := spec.BgpBestPathSelection
		if bgp.Mode != "" {
			networkArgs.BgpBestPathSelectionMode = pulumi.StringPtr(bgp.Mode)
		}
		// Sent explicitly (true or false) whenever the block is present: the
		// provider attribute is Optional+Computed, so omitting it on a
		// true→false transition would silently keep the old value on the
		// live network.
		networkArgs.BgpAlwaysCompareMed = pulumi.BoolPtr(bgp.AlwaysCompareMed)
		if bgp.InterRegionCost != "" {
			networkArgs.BgpInterRegionCost = pulumi.StringPtr(bgp.InterRegionCost)
		}
	}
	// Create-time resource-manager tags (org policy / IAM conditions).
	if len(spec.ResourceManagerTags) > 0 {
		networkArgs.Params = &compute.NetworkParamsArgs{
			ResourceManagerTags: pulumi.ToStringMap(spec.ResourceManagerTags),
		}
	}
	// Empty defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		networkArgs.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdNetwork, err := compute.NewNetwork(ctx,
		"vpc",
		networkArgs,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdComputeService}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vpc network")
	}

	ctx.Export(OpNetworkSelfLink, createdNetwork.SelfLink)
	ctx.Export(OpNetworkName, createdNetwork.Name)
	ctx.Export(OpNetworkId, createdNetwork.ID())
	ctx.Export(OpGatewayIpv4, createdNetwork.GatewayIpv4)
	ctx.Export(OpInternalIpv6Range, createdNetwork.InternalIpv6Range)

	return createdNetwork, nil
}

package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func globalAddress(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpGlobalAddress.Spec

	// Enable the Compute Engine API first so a fresh project works on the
	// first deploy. disable_on_destroy stays false: tearing down one address
	// must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"globaladdress-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.GlobalAddressArgs{
		Name:   pulumi.String(spec.AddressName),
		Labels: pulumi.ToStringMap(locals.GcpLabels),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Address type: EXTERNAL (default) or INTERNAL.
	args.AddressType = pulumi.StringPtr(spec.GetAddressType())

	// IP version: IPV4 (default) or IPV6.
	args.IpVersion = pulumi.StringPtr(spec.GetIpVersion())

	// Specific IP address to reserve. Omit to let GCP assign one.
	if spec.Address != "" {
		args.Address = pulumi.StringPtr(spec.Address)
	}

	// Description.
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	// Network — required for INTERNAL addresses.
	if spec.Network != nil && spec.Network.GetValue() != "" {
		args.Network = pulumi.StringPtr(spec.Network.GetValue())
	}

	// Purpose — VPC_PEERING or PRIVATE_SERVICE_CONNECT.
	if spec.Purpose != "" {
		args.Purpose = pulumi.StringPtr(spec.Purpose)
	}

	// Prefix length — CIDR range for VPC peering.
	if spec.PrefixLength != nil {
		args.PrefixLength = pulumi.IntPtr(int(*spec.PrefixLength))
	}

	createdAddress, err := compute.NewGlobalAddress(ctx, "global-address", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create global address")
	}

	ctx.Export(OpAddress, createdAddress.Address)
	ctx.Export(OpSelfLink, createdAddress.SelfLink)
	ctx.Export(OpCreationTimestamp, createdAddress.CreationTimestamp)

	return nil
}

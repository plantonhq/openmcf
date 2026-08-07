package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// regionalAddress provisions a regional Compute Engine address reservation —
// either a public static IP (EXTERNAL, for Cloud NAT, regional LBs, or VMs)
// or a private IP or range (INTERNAL, for GCE endpoints, internal LB VIPs,
// VPC peering, or IPsec interconnect).
//
// Every field except labels is immutable (ForceNew in the provider): any
// change destroys and recreates the reservation — and a recreated EXTERNAL
// address is a NEW IP, so DNS pointing at the old one breaks. Reserve once,
// reference everywhere.
func regionalAddress(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpAddress.Spec

	// Enable the Compute Engine API first so a fresh project works on the
	// first deploy. disable_on_destroy stays false: tearing down one address
	// must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"address-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.AddressArgs{
		Name:   pulumi.String(spec.AddressName),
		Region: pulumi.String(spec.Region),
		Labels: pulumi.ToStringMap(locals.GcpLabels),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	args.AddressType = pulumi.StringPtr(spec.GetAddressType())
	args.IpVersion = pulumi.StringPtr(spec.GetIpVersion())

	if spec.Address != "" {
		args.Address = pulumi.StringPtr(spec.Address)
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	if spec.Network != nil && spec.Network.GetValue() != "" {
		args.Network = pulumi.StringPtr(spec.Network.GetValue())
	}

	if spec.Subnetwork != nil && spec.Subnetwork.GetValue() != "" {
		args.Subnetwork = pulumi.StringPtr(spec.Subnetwork.GetValue())
	}

	if spec.Purpose != "" {
		args.Purpose = pulumi.StringPtr(spec.Purpose)
	}

	if spec.PrefixLength != nil {
		args.PrefixLength = pulumi.IntPtr(int(*spec.PrefixLength))
	}

	// EXTERNAL-only; spec CEL rejects network_tier on INTERNAL.
	if spec.NetworkTier != "" {
		args.NetworkTier = pulumi.StringPtr(spec.NetworkTier)
	}

	if spec.Ipv6EndpointType != "" {
		args.Ipv6EndpointType = pulumi.StringPtr(spec.Ipv6EndpointType)
	}

	createdAddress, err := compute.NewAddress(ctx, "regional-address", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create regional address")
	}

	ctx.Export(OpAddress, createdAddress.Address)
	ctx.Export(OpSelfLink, createdAddress.SelfLink)
	ctx.Export(OpName, createdAddress.Name)
	// Export the plain region NAME from the spec (matching the Terraform
	// module) — the provider's region attribute can carry a self-link, which
	// API callers and verification cannot use directly.
	ctx.Export(OpRegion, pulumi.String(spec.Region))

	return nil
}

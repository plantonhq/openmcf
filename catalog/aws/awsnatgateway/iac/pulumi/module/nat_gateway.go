package module

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func natGateway(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*ec2.NatGateway, error) {
	spec := locals.AwsNatGateway.Spec
	name := locals.AwsNatGateway.Metadata.Name

	args := &ec2.NatGatewayArgs{
		ConnectivityType: pulumi.String(spec.ConnectivityType),
		Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
			stringmaps.AddEntry(locals.AwsTags, "Name", name)),
	}

	// Placement: zonal gateways live in one subnet; regional gateways span a
	// VPC (availability_mode = "regional" + vpc_id). The spec's CEL rules
	// guarantee exactly one placement form arrives here; unset arms are
	// OMITTED, never sent as empty strings.
	if spec.SubnetId.GetValue() != "" {
		args.SubnetId = pulumi.StringPtr(spec.SubnetId.GetValue())
	}
	if spec.AvailabilityMode != "" {
		args.AvailabilityMode = pulumi.StringPtr(spec.AvailabilityMode)
	}
	if spec.VpcId.GetValue() != "" {
		args.VpcId = pulumi.StringPtr(spec.VpcId.GetValue())
	}

	// Regional gateways: explicit per-AZ layout (manual mode). An empty list
	// means AWS picks the zones and manages addresses itself (auto mode);
	// switching a live gateway between auto and manual replaces it (provider
	// behavior, via its internal auto-mode marker attribute).
	if len(spec.AvailabilityZoneAddresses) > 0 {
		azAddresses := make(ec2.NatGatewayAvailabilityZoneAddressArray, 0, len(spec.AvailabilityZoneAddresses))
		for _, azAddress := range spec.AvailabilityZoneAddresses {
			azArgs := ec2.NatGatewayAvailabilityZoneAddressArgs{}
			if azAddress.AvailabilityZone != "" {
				azArgs.AvailabilityZone = pulumi.StringPtr(azAddress.AvailabilityZone)
			}
			if azAddress.AvailabilityZoneId != "" {
				azArgs.AvailabilityZoneId = pulumi.StringPtr(azAddress.AvailabilityZoneId)
			}
			if len(azAddress.AllocationIds) > 0 {
				allocationIds := make(pulumi.StringArray, 0, len(azAddress.AllocationIds))
				for _, allocationId := range azAddress.AllocationIds {
					allocationIds = append(allocationIds, pulumi.String(allocationId.GetValue()))
				}
				azArgs.AllocationIds = allocationIds
			}
			azAddresses = append(azAddresses, azArgs)
		}
		args.AvailabilityZoneAddresses = azAddresses
	}

	// Public gateways carry an Elastic IP; private gateways carry private IPs.
	// The spec's CEL rules guarantee these never mix, so we just forward whatever
	// is set and leave the rest for AWS to default.
	if spec.AllocationId != nil && spec.AllocationId.GetValue() != "" {
		args.AllocationId = pulumi.String(spec.AllocationId.GetValue())
	}
	if len(spec.SecondaryAllocationIds) > 0 {
		secondary := make(pulumi.StringArray, 0, len(spec.SecondaryAllocationIds))
		for _, a := range spec.SecondaryAllocationIds {
			secondary = append(secondary, pulumi.String(a.GetValue()))
		}
		args.SecondaryAllocationIds = secondary
	}
	if spec.PrivateIp != "" {
		args.PrivateIp = pulumi.String(spec.PrivateIp)
	}
	if len(spec.SecondaryPrivateIpAddresses) > 0 {
		secondary := make(pulumi.StringArray, 0, len(spec.SecondaryPrivateIpAddresses))
		for _, ip := range spec.SecondaryPrivateIpAddresses {
			secondary = append(secondary, pulumi.String(ip))
		}
		args.SecondaryPrivateIpAddresses = secondary
	}
	if spec.SecondaryPrivateIpAddressCount > 0 {
		args.SecondaryPrivateIpAddressCount = pulumi.Int(int(spec.SecondaryPrivateIpAddressCount))
	}

	createdNatGateway, err := ec2.NewNatGateway(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nat gateway")
	}

	ctx.Export(OpNatGatewayId, createdNatGateway.ID())
	ctx.Export(OpPublicIp, createdNatGateway.PublicIp)
	ctx.Export(OpPrivateIp, createdNatGateway.PrivateIp)
	ctx.Export(OpNetworkInterfaceId, createdNatGateway.NetworkInterfaceId)
	ctx.Export(OpSubnetId, pulumi.String(spec.SubnetId.GetValue()))
	ctx.Export(OpRegion, pulumi.String(spec.Region))

	return createdNatGateway, nil
}

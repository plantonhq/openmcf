package module

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func subnet(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*ec2.Subnet, error) {
	spec := locals.AwsSubnet.Spec
	name := locals.AwsSubnet.Metadata.Name

	subnetArgs := &ec2.SubnetArgs{
		VpcId:                                   pulumi.String(spec.VpcId.GetValue()),
		MapPublicIpOnLaunch:                     pulumi.BoolPtr(spec.MapPublicIpOnLaunch),
		AssignIpv6AddressOnCreation:             pulumi.BoolPtr(spec.AssignIpv6AddressOnCreation),
		EnableDns64:                             pulumi.BoolPtr(spec.EnableDns64),
		EnableResourceNameDnsARecordOnLaunch:    pulumi.BoolPtr(spec.EnableResourceNameDnsARecordOnLaunch),
		EnableResourceNameDnsAaaaRecordOnLaunch: pulumi.BoolPtr(spec.EnableResourceNameDnsAaaaRecordOnLaunch),
		Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
			stringmaps.AddEntry(locals.AwsTags, "Name", name)),
	}

	// Exactly one placement form is present (spec CEL): the zone name or the
	// account-stable zone id.
	if spec.AvailabilityZone != "" {
		subnetArgs.AvailabilityZone = pulumi.StringPtr(spec.AvailabilityZone)
	}
	if spec.AvailabilityZoneId != "" {
		subnetArgs.AvailabilityZoneId = pulumi.StringPtr(spec.AvailabilityZoneId)
	}

	// IPv4 addressing: an explicit CIDR, an IPAM allocation, or neither
	// (ipv6-native subnets). The spec CEL guarantees the arms never mix, and
	// unset arms are OMITTED (never sent as empty strings) so AWS applies its
	// own semantics.
	if spec.CidrBlock != "" {
		subnetArgs.CidrBlock = pulumi.StringPtr(spec.CidrBlock)
	}
	if spec.Ipv4IpamPoolId.GetValue() != "" {
		subnetArgs.Ipv4IpamPoolId = pulumi.StringPtr(spec.Ipv4IpamPoolId.GetValue())
	}
	if spec.Ipv4NetmaskLength != nil {
		subnetArgs.Ipv4NetmaskLength = pulumi.IntPtr(int(spec.GetIpv4NetmaskLength()))
	}

	// IPv6 addressing: an explicit CIDR or an IPAM allocation (spec CEL keeps
	// them exclusive); ipv6_native drops IPv4 entirely (IPv6-only subnet).
	if spec.Ipv6CidrBlock != "" {
		subnetArgs.Ipv6CidrBlock = pulumi.StringPtr(spec.Ipv6CidrBlock)
	}
	if spec.Ipv6IpamPoolId.GetValue() != "" {
		subnetArgs.Ipv6IpamPoolId = pulumi.StringPtr(spec.Ipv6IpamPoolId.GetValue())
	}
	if spec.Ipv6NetmaskLength != nil {
		subnetArgs.Ipv6NetmaskLength = pulumi.IntPtr(int(spec.GetIpv6NetmaskLength()))
	}
	if spec.Ipv6Native {
		subnetArgs.Ipv6Native = pulumi.BoolPtr(true)
	}

	if spec.GetPrivateDnsHostnameTypeOnLaunch() != "" {
		subnetArgs.PrivateDnsHostnameTypeOnLaunch = pulumi.StringPtr(spec.GetPrivateDnsHostnameTypeOnLaunch())
	}

	createdSubnet, err := ec2.NewSubnet(ctx, name, subnetArgs, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnet")
	}

	ctx.Export(OpSubnetId, createdSubnet.ID())
	ctx.Export(OpSubnetArn, createdSubnet.Arn)
	ctx.Export(OpAvailabilityZone, createdSubnet.AvailabilityZone)
	ctx.Export(OpCidrBlock, createdSubnet.CidrBlock)
	ctx.Export(OpRegion, pulumi.String(spec.Region))

	return createdSubnet, nil
}

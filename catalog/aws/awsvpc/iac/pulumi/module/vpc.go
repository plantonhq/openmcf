package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func vpc(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsVpc.Spec
	name := locals.AwsVpc.Metadata.Name

	vpcArgs := &ec2.VpcArgs{
		Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
			stringmaps.AddEntry(locals.AwsTags, "Name", name)),
		EnableDnsHostnames:               pulumi.BoolPtr(spec.EnableDnsHostnames),
		EnableNetworkAddressUsageMetrics: pulumi.BoolPtr(spec.EnableNetworkAddressUsageMetrics),
	}

	// The IPAM pool ids are StringValueOrRef spec fields, resolved to plain
	// strings before the module runs (no IPAM pool catalog kind exists yet,
	// so today they carry literal pool ids).
	if spec.CidrBlock != "" {
		vpcArgs.CidrBlock = pulumi.StringPtr(spec.CidrBlock)
	}
	if spec.Ipv4IpamPoolId.GetValue() != "" {
		vpcArgs.Ipv4IpamPoolId = pulumi.StringPtr(spec.Ipv4IpamPoolId.GetValue())
	}
	if spec.Ipv4NetmaskLength != 0 {
		vpcArgs.Ipv4NetmaskLength = pulumi.IntPtr(int(spec.Ipv4NetmaskLength))
	}
	if spec.InstanceTenancy != "" {
		vpcArgs.InstanceTenancy = pulumi.StringPtr(spec.InstanceTenancy)
	}
	// enable_dns_support is proto3 optional: honor an explicit value, otherwise
	// leave the argument unset so AWS applies its default (DNS support on).
	if spec.EnableDnsSupport != nil {
		vpcArgs.EnableDnsSupport = pulumi.BoolPtr(spec.GetEnableDnsSupport())
	}

	if spec.AssignGeneratedIpv6CidrBlock {
		vpcArgs.AssignGeneratedIpv6CidrBlock = pulumi.BoolPtr(true)
	}
	if spec.Ipv6CidrBlock != "" {
		vpcArgs.Ipv6CidrBlock = pulumi.StringPtr(spec.Ipv6CidrBlock)
	}
	if spec.Ipv6CidrBlockNetworkBorderGroup != "" {
		vpcArgs.Ipv6CidrBlockNetworkBorderGroup = pulumi.StringPtr(spec.Ipv6CidrBlockNetworkBorderGroup)
	}
	if spec.Ipv6IpamPoolId.GetValue() != "" {
		vpcArgs.Ipv6IpamPoolId = pulumi.StringPtr(spec.Ipv6IpamPoolId.GetValue())
	}
	if spec.Ipv6NetmaskLength != 0 {
		vpcArgs.Ipv6NetmaskLength = pulumi.IntPtr(int(spec.Ipv6NetmaskLength))
	}

	createdVpc, err := ec2.NewVpc(ctx, name, vpcArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create vpc")
	}

	// Associate each secondary IPv4 CIDR as its own resource so it can be added
	// or removed without recreating the VPC. An entry names an explicit CIDR,
	// an IPAM-sized allocation, or a pool-pinned block (spec CEL enforces the
	// shapes). The association-id output map is keyed IDENTICALLY to the
	// Terraform module's for_each (literal CIDR, else "ipam-<index>") -- the
	// import recipes resolve per-instance IDs through these keys, so the two
	// engines must agree on them.
	secondaryIpv4AssociationIds := pulumi.StringMap{}
	for i, entry := range spec.SecondaryIpv4Cidrs {
		assocArgs := &ec2.VpcIpv4CidrBlockAssociationArgs{
			VpcId: createdVpc.ID(),
		}
		key := fmt.Sprintf("ipam-%d", i)
		if entry.CidrBlock != "" {
			key = entry.CidrBlock
			assocArgs.CidrBlock = pulumi.StringPtr(entry.CidrBlock)
		}
		if entry.IpamPoolId.GetValue() != "" {
			assocArgs.Ipv4IpamPoolId = pulumi.StringPtr(entry.IpamPoolId.GetValue())
		}
		if entry.NetmaskLength != 0 {
			assocArgs.Ipv4NetmaskLength = pulumi.IntPtr(int(entry.NetmaskLength))
		}
		createdAssociation, err := ec2.NewVpcIpv4CidrBlockAssociation(ctx,
			fmt.Sprintf("%s-secondary-%d", name, i),
			assocArgs, pulumi.Provider(provider), pulumi.Parent(createdVpc))
		if err != nil {
			return errors.Wrapf(err, "failed to associate secondary ipv4 cidr %d", i)
		}
		secondaryIpv4AssociationIds[key] = createdAssociation.ID().ToStringOutput()
	}

	// Associate each secondary IPv6 CIDR as its own resource. Exactly one
	// source per entry (spec CEL): an Amazon-provided block, a BYOIP public
	// pool, or an IPAM pool. Output-map keys mirror the Terraform for_each
	// (pinned CIDR, else "ipv6-<index>") for the same import-recipe reason.
	secondaryIpv6AssociationIds := pulumi.StringMap{}
	for i, entry := range spec.SecondaryIpv6Cidrs {
		assocArgs := &ec2.VpcIpv6CidrBlockAssociationArgs{
			VpcId: createdVpc.ID(),
		}
		if entry.AssignGenerated {
			assocArgs.AssignGeneratedIpv6CidrBlock = pulumi.BoolPtr(true)
		}
		if entry.Ipv6Pool != "" {
			assocArgs.Ipv6Pool = pulumi.StringPtr(entry.Ipv6Pool)
		}
		if entry.IpamPoolId.GetValue() != "" {
			assocArgs.Ipv6IpamPoolId = pulumi.StringPtr(entry.IpamPoolId.GetValue())
		}
		key := fmt.Sprintf("ipv6-%d", i)
		if entry.CidrBlock != "" {
			key = entry.CidrBlock
			assocArgs.Ipv6CidrBlock = pulumi.StringPtr(entry.CidrBlock)
		}
		if entry.NetmaskLength != 0 {
			assocArgs.Ipv6NetmaskLength = pulumi.IntPtr(int(entry.NetmaskLength))
		}
		createdAssociation, err := ec2.NewVpcIpv6CidrBlockAssociation(ctx,
			fmt.Sprintf("%s-secondary-ipv6-%d", name, i),
			assocArgs, pulumi.Provider(provider), pulumi.Parent(createdVpc))
		if err != nil {
			return errors.Wrapf(err, "failed to associate secondary ipv6 cidr %d", i)
		}
		secondaryIpv6AssociationIds[key] = createdAssociation.ID().ToStringOutput()
	}

	// VPC Encryption Control: AWS's VPC-wide monitor/enforce switch for
	// encryption in transit. Rendered only when configured; exclusions are
	// sent enable/disable per service and only apply in enforce mode (spec
	// CEL keeps monitor-mode exclusions out).
	// Empty when encryption control is not configured (mirrors the Terraform
	// module's try(...[0].id, "") output shape).
	encryptionControlId := pulumi.String("").ToStringOutput()
	if spec.EncryptionControl != nil {
		ec := spec.EncryptionControl
		exclusion := func(excluded bool) pulumi.StringPtrInput {
			if excluded {
				return pulumi.StringPtr("enable")
			}
			return pulumi.StringPtr("disable")
		}
		createdEncryptionControl, err := ec2.NewVpcEncryptionControl(ctx, name+"-encryption-control",
			&ec2.VpcEncryptionControlArgs{
				VpcId:                              createdVpc.ID(),
				Mode:                               pulumi.String(ec.Mode),
				InternetGatewayExclusion:           exclusion(ec.ExcludeInternetGateway),
				EgressOnlyInternetGatewayExclusion: exclusion(ec.ExcludeEgressOnlyInternetGateway),
				NatGatewayExclusion:                exclusion(ec.ExcludeNatGateway),
				VirtualPrivateGatewayExclusion:     exclusion(ec.ExcludeVirtualPrivateGateway),
				VpcPeeringExclusion:                exclusion(ec.ExcludeVpcPeering),
				VpcLatticeExclusion:                exclusion(ec.ExcludeVpcLattice),
				LambdaExclusion:                    exclusion(ec.ExcludeLambda),
				ElasticFileSystemExclusion:         exclusion(ec.ExcludeElasticFileSystem),
				Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
					stringmaps.AddEntry(locals.AwsTags, "Name", name)),
			}, pulumi.Provider(provider), pulumi.Parent(createdVpc))
		if err != nil {
			return errors.Wrap(err, "failed to create vpc encryption control")
		}
		encryptionControlId = createdEncryptionControl.ID().ToStringOutput()
	}

	ctx.Export(OpVpcId, createdVpc.ID())
	ctx.Export(OpVpcArn, createdVpc.Arn)
	ctx.Export(OpCidrBlock, createdVpc.CidrBlock)
	ctx.Export(OpIpv6CidrBlock, createdVpc.Ipv6CidrBlock)
	ctx.Export(OpOwnerId, createdVpc.OwnerId)
	ctx.Export(OpMainRouteTableId, createdVpc.MainRouteTableId)
	ctx.Export(OpDefaultSecurityGroupId, createdVpc.DefaultSecurityGroupId)
	ctx.Export(OpDefaultNetworkAclId, createdVpc.DefaultNetworkAclId)
	ctx.Export(OpDefaultRouteTableId, createdVpc.DefaultRouteTableId)
	ctx.Export(OpRegion, pulumi.String(spec.Region))
	ctx.Export(OpSecondaryIpv4CidrAssociationIds, secondaryIpv4AssociationIds)
	ctx.Export(OpSecondaryIpv6CidrAssociationIds, secondaryIpv6AssociationIds)
	ctx.Export(OpEncryptionControlId, encryptionControlId)

	return nil
}

package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/efs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// MountTargetResults holds the per-subnet mount target outputs keyed by the
// resolved subnet ID (the same keys the Terraform module's for_each produces).
type MountTargetResults struct {
	MountTargetIds           pulumi.StringMap
	MountTargetIps           pulumi.StringMap
	MountTargetIpv6Addresses pulumi.StringMap
	MountTargetDnsNames      pulumi.StringMap
}

func mountTargets(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, fs *efs.FileSystem) (*MountTargetResults, error) {
	spec := locals.AwsElasticFileSystem.Spec

	// Security groups apply uniformly to every mount target (they gate the
	// same NFS clients regardless of AZ). When empty, AWS attaches the VPC's
	// default security group.
	var securityGroups pulumi.StringArray
	for _, sg := range spec.SecurityGroupIds {
		if sg.GetValue() != "" {
			securityGroups = append(securityGroups, pulumi.String(sg.GetValue()))
		}
	}

	results := &MountTargetResults{
		MountTargetIds:           pulumi.StringMap{},
		MountTargetIps:           pulumi.StringMap{},
		MountTargetIpv6Addresses: pulumi.StringMap{},
		MountTargetDnsNames:      pulumi.StringMap{},
	}

	// AWS allows at most one mount target per Availability Zone, and returns
	// the SAME mount-target ID to parallel same-AZ create calls — so the spec
	// requires one subnet per AZ and the provider serializes creation per AZ.
	for i, mtSpec := range spec.MountTargets {
		subnetId := mtSpec.SubnetId.GetValue()
		if subnetId == "" {
			continue
		}

		args := &efs.MountTargetArgs{
			FileSystemId: fs.ID(),
			SubnetId:     pulumi.String(subnetId),
		}

		if len(securityGroups) > 0 {
			args.SecurityGroups = securityGroups
		}

		// Static IPv4 address from the subnet's CIDR. ForceNew.
		if mtSpec.IpAddress != "" {
			args.IpAddress = pulumi.StringPtr(mtSpec.IpAddress)
		}

		// Address family (IPV4_ONLY / IPV6_ONLY / DUAL_STACK). Empty keeps the
		// AWS default (IPV4_ONLY). ForceNew.
		if mtSpec.IpAddressType != "" {
			args.IpAddressType = pulumi.StringPtr(mtSpec.IpAddressType)
		}

		// Static IPv6 address (requires an IPv6-capable address type;
		// CEL-enforced). ForceNew.
		if mtSpec.Ipv6Address != "" {
			args.Ipv6Address = pulumi.StringPtr(mtSpec.Ipv6Address)
		}

		mt, err := efs.NewMountTarget(ctx, fmt.Sprintf("mount-target-%d", i), args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create mount target for subnet %s", subnetId)
		}

		results.MountTargetIds[subnetId] = mt.ID().ToStringOutput()
		results.MountTargetIps[subnetId] = mt.IpAddress
		results.MountTargetIpv6Addresses[subnetId] = mt.Ipv6Address
		results.MountTargetDnsNames[subnetId] = mt.MountTargetDnsName
	}

	return results, nil
}

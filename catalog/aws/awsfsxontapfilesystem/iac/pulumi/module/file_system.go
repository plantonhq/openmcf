package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/fsx"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func fileSystem(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*fsx.OntapFileSystem, error) {
	spec := locals.AwsFsxOntapFileSystem.Spec
	name := locals.AwsFsxOntapFileSystem.Metadata.Name

	subnetIds := make(pulumi.StringArray, 0, len(spec.SubnetIds))
	for _, s := range spec.SubnetIds {
		subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
	}

	args := &fsx.OntapFileSystemArgs{
		SubnetIds:       subnetIds,
		StorageCapacity: pulumi.Int(int(spec.StorageCapacityGib)),
		Tags:            pulumi.ToStringMap(locals.AwsTags),
	}

	// Throughput is sized through exactly one arm (spec-enforced XOR):
	// whole-file-system ThroughputCapacity (the first-generation sizing) or
	// per-HA-pair ThroughputCapacityPerHaPair (required for scale-out and the
	// second generation's tiers). The unset arm stays nil so the provider
	// omits it — sending both is an AWS error.
	if spec.ThroughputCapacity != nil {
		args.ThroughputCapacity = pulumi.IntPtr(int(spec.GetThroughputCapacity()))
	}
	if spec.ThroughputCapacityPerHaPair != nil {
		args.ThroughputCapacityPerHaPair = pulumi.IntPtr(int(spec.GetThroughputCapacityPerHaPair()))
	}

	// Deployment type (spec default SINGLE_AZ_2, materialized before the
	// module runs). ForceNew.
	if spec.GetDeploymentType() != "" {
		args.DeploymentType = pulumi.String(spec.GetDeploymentType())
	}

	// Storage type — ONTAP is SSD-only (spec-enforced); sent explicitly so
	// the preview states the storage class. ForceNew.
	if spec.GetStorageType() != "" {
		args.StorageType = pulumi.StringPtr(spec.GetStorageType())
	}

	// HA pairs: sent whenever set (not just >1) so an explicit single-pair
	// choice reaches AWS and both engines plan identically. AWS requires the
	// per-HA-pair throughput arm to be re-sent whenever HA pairs change —
	// both values flowing from the spec keeps that invariant automatically.
	if spec.HaPairs != nil {
		args.HaPairs = pulumi.IntPtr(int(spec.GetHaPairs()))
	}

	// Security groups (ForceNew). Omitted entirely when empty — AWS then
	// attaches the VPC default SG, and an empty set is a different plan.
	if len(spec.SecurityGroupIds) > 0 {
		sgIds := make(pulumi.StringArray, 0, len(spec.SecurityGroupIds))
		for _, sg := range spec.SecurityGroupIds {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
		args.SecurityGroupIds = sgIds
	}

	// The provider requires PreferredSubnetId for every deployment type,
	// while the decision only exists for multi-AZ — the spec requires it
	// there and forbids it for single-AZ, where it derives deterministically
	// from the only subnet (the Terraform module derives the same way).
	if spec.PreferredSubnetId != nil && spec.PreferredSubnetId.GetValue() != "" {
		args.PreferredSubnetId = pulumi.String(spec.PreferredSubnetId.GetValue())
	} else {
		args.PreferredSubnetId = pulumi.String(spec.SubnetIds[0].GetValue())
	}

	// Endpoint IP address range (multi-AZ only).
	if spec.EndpointIpAddressRange != "" {
		args.EndpointIpAddressRange = pulumi.StringPtr(spec.EndpointIpAddressRange)
	}

	// Route table IDs (multi-AZ only).
	if len(spec.RouteTableIds) > 0 {
		rtIds := make(pulumi.StringArray, 0, len(spec.RouteTableIds))
		for _, rt := range spec.RouteTableIds {
			rtIds = append(rtIds, pulumi.String(rt.GetValue()))
		}
		args.RouteTableIds = rtIds
	}

	// Customer-managed KMS key for encryption at rest.
	if spec.KmsKeyId != nil && spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	// ONTAP admin password (sensitive).
	if spec.FsxAdminPassword != "" {
		args.FsxAdminPassword = pulumi.StringPtr(spec.FsxAdminPassword)
	}

	// Disk IOPS configuration.
	if spec.DiskIopsConfiguration != nil {
		iopsArgs := &fsx.OntapFileSystemDiskIopsConfigurationArgs{}
		if spec.DiskIopsConfiguration.GetMode() != "" {
			iopsArgs.Mode = pulumi.StringPtr(spec.DiskIopsConfiguration.GetMode())
		}
		if spec.DiskIopsConfiguration.Iops > 0 {
			iopsArgs.Iops = pulumi.IntPtr(int(spec.DiskIopsConfiguration.Iops))
		}
		args.DiskIopsConfiguration = iopsArgs
	}

	// Automatic backup retention: zero is a real value ("no automatic
	// backups"), so the resolved value is sent whenever present rather than
	// only when positive.
	if spec.AutomaticBackupRetentionDays != nil {
		args.AutomaticBackupRetentionDays = pulumi.IntPtr(int(spec.GetAutomaticBackupRetentionDays()))
	}

	// Daily automatic backup start time (HH:MM format).
	if spec.DailyAutomaticBackupStartTime != "" {
		args.DailyAutomaticBackupStartTime = pulumi.StringPtr(spec.DailyAutomaticBackupStartTime)
	}

	// Weekly maintenance start time (d:HH:MM format).
	if spec.WeeklyMaintenanceStartTime != "" {
		args.WeeklyMaintenanceStartTime = pulumi.StringPtr(spec.WeeklyMaintenanceStartTime)
	}

	createdFs, err := fsx.NewOntapFileSystem(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fsx ontap file system")
	}

	return createdFs, nil
}

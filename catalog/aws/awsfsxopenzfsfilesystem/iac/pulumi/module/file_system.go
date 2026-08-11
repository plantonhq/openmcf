package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/fsx"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// fileSystem creates the FSx for OpenZFS file system. FSx file systems have no
// cloud name argument — the console name is the Name tag, which locals pins to
// metadata.name on the same basis as the Terraform module.
func fileSystem(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*fsx.OpenZfsFileSystem, error) {
	spec := locals.AwsFsxOpenzfsFileSystem.Spec
	name := locals.AwsFsxOpenzfsFileSystem.Metadata.Name

	// Subnet IDs: one for the single-AZ types, two for MULTI_AZ_1
	// (CEL-enforced at validation).
	subnetIds := make(pulumi.StringArray, 0, len(spec.SubnetIds))
	for _, s := range spec.SubnetIds {
		subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
	}

	args := &fsx.OpenZfsFileSystemArgs{
		SubnetIds:          subnetIds,
		ThroughputCapacity: pulumi.Int(int(spec.ThroughputCapacity)),
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}

	// Storage capacity is presence-honest: unset means the capacity comes
	// from a backup restore or the elastic INTELLIGENT_TIERING class, and the
	// argument must be omitted entirely (AWS rejects it for both shapes).
	if spec.StorageCapacityGib != nil {
		args.StorageCapacity = pulumi.IntPtr(int(spec.GetStorageCapacityGib()))
	}

	// Deployment type — always sent; the resolved spec default is SINGLE_AZ_2.
	if spec.GetDeploymentType() != "" {
		args.DeploymentType = pulumi.String(spec.GetDeploymentType())
	}

	// Storage class (SSD / INTELLIGENT_TIERING). ForceNew.
	if spec.GetStorageType() != "" {
		args.StorageType = pulumi.StringPtr(spec.GetStorageType())
	}

	// Provisioned SSD read cache for INTELLIGENT_TIERING.
	if spec.ReadCacheConfiguration != nil {
		cacheArgs := &fsx.OpenZfsFileSystemReadCacheConfigurationArgs{
			SizingMode: pulumi.StringPtr(spec.ReadCacheConfiguration.SizingMode),
		}
		if spec.ReadCacheConfiguration.SizeGib != nil {
			cacheArgs.Size = pulumi.IntPtr(int(spec.ReadCacheConfiguration.GetSizeGib()))
		}
		args.ReadCacheConfiguration = cacheArgs
	}

	// Security groups. ForceNew; empty lets AWS attach the VPC default SG.
	if len(spec.SecurityGroupIds) > 0 {
		sgIds := make(pulumi.StringArray, 0, len(spec.SecurityGroupIds))
		for _, sg := range spec.SecurityGroupIds {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
		args.SecurityGroupIds = sgIds
	}

	// MULTI_AZ_1 networking: the active file server's subnet, the floating
	// endpoint CIDR, and the route tables AWS manages routes in.
	if spec.PreferredSubnetId != nil && spec.PreferredSubnetId.GetValue() != "" {
		args.PreferredSubnetId = pulumi.StringPtr(spec.PreferredSubnetId.GetValue())
	}
	if spec.EndpointIpAddressRange != "" {
		args.EndpointIpAddressRange = pulumi.StringPtr(spec.EndpointIpAddressRange)
	}
	if len(spec.RouteTableIds) > 0 {
		rtIds := make(pulumi.StringArray, 0, len(spec.RouteTableIds))
		for _, rt := range spec.RouteTableIds {
			rtIds = append(rtIds, pulumi.String(rt.GetValue()))
		}
		args.RouteTableIds = rtIds
	}

	// Customer-managed KMS key for encryption at rest. ForceNew.
	if spec.KmsKeyId != nil && spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	// Restore-from-backup create shape. ForceNew.
	if spec.BackupId != "" {
		args.BackupId = pulumi.StringPtr(spec.BackupId)
	}

	// Disk IOPS configuration.
	if spec.DiskIopsConfiguration != nil {
		iopsArgs := &fsx.OpenZfsFileSystemDiskIopsConfigurationArgs{}
		if spec.DiskIopsConfiguration.GetMode() != "" {
			iopsArgs.Mode = pulumi.StringPtr(spec.DiskIopsConfiguration.GetMode())
		}
		if spec.DiskIopsConfiguration.Iops != nil {
			iopsArgs.Iops = pulumi.IntPtr(int(spec.DiskIopsConfiguration.GetIops()))
		}
		args.DiskIopsConfiguration = iopsArgs
	}

	// Root volume configuration — the file system's default NFS mount target.
	if spec.RootVolumeConfiguration != nil {
		rootVolArgs := &fsx.OpenZfsFileSystemRootVolumeConfigurationArgs{}

		if spec.RootVolumeConfiguration.GetDataCompressionType() != "" {
			rootVolArgs.DataCompressionType = pulumi.StringPtr(spec.RootVolumeConfiguration.GetDataCompressionType())
		}

		// Always sent (matching the Terraform module): the provider attribute
		// is Optional+Computed, so an omitted read_only would let out-of-band
		// flips to read-only be adopted silently instead of reverted. An
		// explicit resolved value keeps drift correction on both engines.
		rootVolArgs.ReadOnly = pulumi.BoolPtr(spec.RootVolumeConfiguration.ReadOnly)

		if spec.RootVolumeConfiguration.GetRecordSizeKib() > 0 {
			rootVolArgs.RecordSizeKib = pulumi.IntPtr(int(spec.RootVolumeConfiguration.GetRecordSizeKib()))
		}

		// ForceNew inside an otherwise in-place block: flipping this replaces
		// the whole file system.
		if spec.RootVolumeConfiguration.CopyTagsToSnapshots {
			rootVolArgs.CopyTagsToSnapshots = pulumi.BoolPtr(true)
		}

		if spec.RootVolumeConfiguration.NfsExports != nil &&
			len(spec.RootVolumeConfiguration.NfsExports.ClientConfigurations) > 0 {

			clientConfigs := make(fsx.OpenZfsFileSystemRootVolumeConfigurationNfsExportsClientConfigurationArray, 0,
				len(spec.RootVolumeConfiguration.NfsExports.ClientConfigurations))

			for _, cc := range spec.RootVolumeConfiguration.NfsExports.ClientConfigurations {
				opts := make(pulumi.StringArray, 0, len(cc.Options))
				for _, o := range cc.Options {
					opts = append(opts, pulumi.String(o))
				}
				clientConfigs = append(clientConfigs, &fsx.OpenZfsFileSystemRootVolumeConfigurationNfsExportsClientConfigurationArgs{
					Clients: pulumi.String(cc.Clients),
					Options: opts,
				})
			}

			rootVolArgs.NfsExports = &fsx.OpenZfsFileSystemRootVolumeConfigurationNfsExportsArgs{
				ClientConfigurations: clientConfigs,
			}
		}

		if len(spec.RootVolumeConfiguration.UserAndGroupQuotas) > 0 {
			quotas := make(fsx.OpenZfsFileSystemRootVolumeConfigurationUserAndGroupQuotaArray, 0,
				len(spec.RootVolumeConfiguration.UserAndGroupQuotas))

			for _, q := range spec.RootVolumeConfiguration.UserAndGroupQuotas {
				quotas = append(quotas, &fsx.OpenZfsFileSystemRootVolumeConfigurationUserAndGroupQuotaArgs{
					Id:                      pulumi.Int(int(q.Id)),
					StorageCapacityQuotaGib: pulumi.Int(int(q.StorageCapacityQuotaGib)),
					Type:                    pulumi.String(q.Type),
				})
			}

			rootVolArgs.UserAndGroupQuotas = quotas
		}

		args.RootVolumeConfiguration = rootVolArgs
	}

	// Automatic backups. Zero is a real value ("no automatic backups"), so
	// the resolved default is always sent explicitly.
	if spec.AutomaticBackupRetentionDays != nil {
		args.AutomaticBackupRetentionDays = pulumi.IntPtr(int(spec.GetAutomaticBackupRetentionDays()))
	}
	if spec.DailyAutomaticBackupStartTime != "" {
		args.DailyAutomaticBackupStartTime = pulumi.StringPtr(spec.DailyAutomaticBackupStartTime)
	}
	if spec.CopyTagsToBackups {
		args.CopyTagsToBackups = pulumi.BoolPtr(true)
	}
	if spec.CopyTagsToVolumes {
		args.CopyTagsToVolumes = pulumi.BoolPtr(true)
	}

	// The final-backup decision is presence-honest: an explicit false ("take
	// a final backup on delete") must reach AWS, so the resolved value is
	// always sent — only-send-true would silently revert an explicit false to
	// the provider default.
	if spec.SkipFinalBackup != nil {
		args.SkipFinalBackup = pulumi.BoolPtr(spec.GetSkipFinalBackup())
	}
	if len(spec.FinalBackupTags) > 0 {
		args.FinalBackupTags = pulumi.ToStringMap(spec.FinalBackupTags)
	}

	// Cascading-delete opt-in: without DELETE_CHILD_VOLUMES_AND_SNAPSHOTS,
	// deletion fails while child volumes or snapshots exist.
	if len(spec.DeleteOptions) > 0 {
		args.DeleteOptions = pulumi.ToStringArray(spec.DeleteOptions)
	}

	// Weekly maintenance window ("d:HH:MM").
	if spec.WeeklyMaintenanceStartTime != "" {
		args.WeeklyMaintenanceStartTime = pulumi.StringPtr(spec.WeeklyMaintenanceStartTime)
	}

	createdFs, err := fsx.NewOpenZfsFileSystem(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fsx openzfs file system")
	}

	return createdFs, nil
}

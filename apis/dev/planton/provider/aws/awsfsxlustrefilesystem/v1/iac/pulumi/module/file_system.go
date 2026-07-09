package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/fsx"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// fileSystem creates the FSx for Lustre file system. FSx file systems have no
// cloud name argument — the console name is the Name tag, which locals pins to
// metadata.name on the same basis as the Terraform module.
func fileSystem(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*fsx.LustreFileSystem, error) {
	spec := locals.AwsFsxLustreFileSystem.Spec
	name := locals.AwsFsxLustreFileSystem.Metadata.Name

	args := &fsx.LustreFileSystemArgs{
		// Lustre is single-AZ: the SDK models subnet_ids as a single string
		// (the underlying list is capped at exactly one subnet).
		SubnetIds: pulumi.String(spec.SubnetId.GetValue()),
		Tags:      pulumi.ToStringMap(locals.AwsTags),
	}

	// Storage capacity is presence-honest: unset means "capacity comes from
	// somewhere else" — a backup restore or the elastic INTELLIGENT_TIERING
	// class — and the argument must be omitted entirely.
	if spec.StorageCapacityGib != nil {
		args.StorageCapacity = pulumi.IntPtr(int(spec.GetStorageCapacityGib()))
	}

	// Deployment type (spec default SCRATCH_2; the provider's own default is
	// the legacy SCRATCH_1, so always send the resolved value).
	if spec.GetDeploymentType() != "" {
		args.DeploymentType = pulumi.StringPtr(spec.GetDeploymentType())
	}

	// Storage class (SSD / HDD / INTELLIGENT_TIERING). ForceNew.
	if spec.GetStorageType() != "" {
		args.StorageType = pulumi.StringPtr(spec.GetStorageType())
	}

	// Per-TiB throughput for provisioned-capacity PERSISTENT deployments.
	if spec.PerUnitStorageThroughput != nil {
		args.PerUnitStorageThroughput = pulumi.IntPtr(int(spec.GetPerUnitStorageThroughput()))
	}

	// Absolute throughput for the INTELLIGENT_TIERING storage class.
	if spec.ThroughputCapacity != nil {
		args.ThroughputCapacity = pulumi.IntPtr(int(spec.GetThroughputCapacity()))
	}

	// LZ4 compression toggle (in-place; new writes only).
	if spec.GetDataCompressionType() != "" {
		args.DataCompressionType = pulumi.StringPtr(spec.GetDataCompressionType())
	}

	// Lustre version pin ("x.y"). Upgrades apply in place; downgrades replace.
	if spec.FileSystemTypeVersion != "" {
		args.FileSystemTypeVersion = pulumi.StringPtr(spec.FileSystemTypeVersion)
	}

	// EFA/GPUDirect Storage. ForceNew and only meaningful at creation, so the
	// flag is sent only when enabled — AWS computes the default (false).
	if spec.EfaEnabled {
		args.EfaEnabled = pulumi.BoolPtr(true)
	}

	// HDD read cache decision (required by AWS for HDD storage; CEL-enforced).
	if spec.DriveCacheType != "" {
		args.DriveCacheType = pulumi.StringPtr(spec.DriveCacheType)
	}

	// Provisioned SSD read cache for INTELLIGENT_TIERING.
	if spec.DataReadCacheConfiguration != nil {
		cacheArgs := &fsx.LustreFileSystemDataReadCacheConfigurationArgs{
			SizingMode: pulumi.String(spec.DataReadCacheConfiguration.SizingMode),
		}
		if spec.DataReadCacheConfiguration.SizeGib != nil {
			cacheArgs.Size = pulumi.IntPtr(int(spec.DataReadCacheConfiguration.GetSizeGib()))
		}
		args.DataReadCacheConfiguration = cacheArgs
	}

	// Security groups. ForceNew; empty lets AWS attach the VPC default SG.
	if len(spec.SecurityGroupIds) > 0 {
		sgIds := make(pulumi.StringArray, 0, len(spec.SecurityGroupIds))
		for _, sg := range spec.SecurityGroupIds {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
		args.SecurityGroupIds = sgIds
	}

	// Customer-managed KMS key for encryption at rest. ForceNew.
	if spec.KmsKeyId != nil && spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	// Restore-from-backup create shape. ForceNew.
	if spec.BackupId != "" {
		args.BackupId = pulumi.StringPtr(spec.BackupId)
	}

	// Legacy S3 data repository link (SCRATCH_1/SCRATCH_2/PERSISTENT_1 only —
	// PERSISTENT_2 uses data repository associations). All four knobs ForceNew.
	if spec.ImportPath != "" {
		args.ImportPath = pulumi.StringPtr(spec.ImportPath)
	}
	if spec.ExportPath != "" {
		args.ExportPath = pulumi.StringPtr(spec.ExportPath)
	}
	if spec.AutoImportPolicy != "" {
		args.AutoImportPolicy = pulumi.StringPtr(spec.AutoImportPolicy)
	}
	if spec.ImportedFileChunkSize != nil {
		args.ImportedFileChunkSize = pulumi.IntPtr(int(spec.GetImportedFileChunkSize()))
	}

	// POSIX root squash. Updatable in place.
	if spec.RootSquashConfiguration != nil {
		squashArgs := &fsx.LustreFileSystemRootSquashConfigurationArgs{}
		if spec.RootSquashConfiguration.RootSquash != "" {
			squashArgs.RootSquash = pulumi.StringPtr(spec.RootSquashConfiguration.RootSquash)
		}
		if len(spec.RootSquashConfiguration.NoSquashNids) > 0 {
			squashArgs.NoSquashNids = pulumi.ToStringArray(spec.RootSquashConfiguration.NoSquashNids)
		}
		args.RootSquashConfiguration = squashArgs
	}

	// CloudWatch logging for data repository events.
	if spec.LogConfiguration != nil {
		logArgs := &fsx.LustreFileSystemLogConfigurationArgs{}
		if spec.LogConfiguration.Destination != nil && spec.LogConfiguration.Destination.GetValue() != "" {
			logArgs.Destination = pulumi.StringPtr(spec.LogConfiguration.Destination.GetValue())
		}
		if spec.LogConfiguration.GetLevel() != "" {
			logArgs.Level = pulumi.StringPtr(spec.LogConfiguration.GetLevel())
		}
		args.LogConfiguration = logArgs
	}

	// Metadata IOPS configuration (PERSISTENT_2 only; CEL-enforced).
	if spec.MetadataConfiguration != nil {
		metaArgs := &fsx.LustreFileSystemMetadataConfigurationArgs{}
		if spec.MetadataConfiguration.GetMode() != "" {
			metaArgs.Mode = pulumi.StringPtr(spec.MetadataConfiguration.GetMode())
		}
		if spec.MetadataConfiguration.Iops != nil {
			metaArgs.Iops = pulumi.IntPtr(int(spec.MetadataConfiguration.GetIops()))
		}
		args.MetadataConfiguration = metaArgs
	}

	// Automatic backups (PERSISTENT only). Zero is a real value ("no automatic
	// backups"), so the resolved default is always sent explicitly.
	if spec.AutomaticBackupRetentionDays != nil {
		args.AutomaticBackupRetentionDays = pulumi.IntPtr(int(spec.GetAutomaticBackupRetentionDays()))
	}
	if spec.DailyAutomaticBackupStartTime != "" {
		args.DailyAutomaticBackupStartTime = pulumi.StringPtr(spec.DailyAutomaticBackupStartTime)
	}
	if spec.CopyTagsToBackups {
		args.CopyTagsToBackups = pulumi.BoolPtr(true)
	}

	// The final-backup decision is presence-honest: an explicit false ("take a
	// final backup on delete") must reach AWS, so the resolved value is always
	// sent — only-send-true would silently revert an explicit false to the
	// provider default.
	if spec.SkipFinalBackup != nil {
		args.SkipFinalBackup = pulumi.BoolPtr(spec.GetSkipFinalBackup())
	}
	if len(spec.FinalBackupTags) > 0 {
		args.FinalBackupTags = pulumi.ToStringMap(spec.FinalBackupTags)
	}

	// Weekly maintenance window ("d:HH:MM").
	if spec.WeeklyMaintenanceStartTime != "" {
		args.WeeklyMaintenanceStartTime = pulumi.StringPtr(spec.WeeklyMaintenanceStartTime)
	}

	createdFs, err := fsx.NewLustreFileSystem(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fsx lustre file system")
	}

	return createdFs, nil
}

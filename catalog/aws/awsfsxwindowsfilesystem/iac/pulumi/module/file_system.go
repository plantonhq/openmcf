package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/fsx"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func fileSystem(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*fsx.WindowsFileSystem, error) {
	spec := locals.AwsFsxWindowsFileSystem.Spec
	name := locals.AwsFsxWindowsFileSystem.Metadata.Name

	// Subnet IDs (required, at least 1).
	subnetIds := make(pulumi.StringArray, 0, len(spec.SubnetIds))
	for _, s := range spec.SubnetIds {
		subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
	}

	args := &fsx.WindowsFileSystemArgs{
		SubnetIds:          subnetIds,
		ThroughputCapacity: pulumi.Int(int(spec.ThroughputCapacity)),
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}

	// Storage capacity is presence-honest: unset means the capacity comes
	// from a backup restore and the argument must be omitted entirely.
	if spec.StorageCapacityGib != nil {
		args.StorageCapacity = pulumi.IntPtr(int(spec.GetStorageCapacityGib()))
	}

	// Deployment type (optional).
	if spec.GetDeploymentType() != "" {
		args.DeploymentType = pulumi.StringPtr(spec.GetDeploymentType())
	}

	// Storage type (optional).
	if spec.GetStorageType() != "" {
		args.StorageType = pulumi.StringPtr(spec.GetStorageType())
	}

	// Security groups.
	if len(spec.SecurityGroupIds) > 0 {
		sgIds := make(pulumi.StringArray, 0, len(spec.SecurityGroupIds))
		for _, sg := range spec.SecurityGroupIds {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
		args.SecurityGroupIds = sgIds
	}

	// Preferred subnet (MULTI_AZ_1 only).
	if spec.PreferredSubnetId != nil && spec.PreferredSubnetId.GetValue() != "" {
		args.PreferredSubnetId = pulumi.StringPtr(spec.PreferredSubnetId.GetValue())
	}

	// Customer-managed KMS key for encryption at rest.
	if spec.KmsKeyId != nil && spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	// Restore-from-backup create shape. ForceNew.
	if spec.BackupId != "" {
		args.BackupId = pulumi.StringPtr(spec.BackupId)
	}

	// Active Directory ID (AWS Managed Microsoft AD).
	if spec.ActiveDirectoryId != nil && spec.ActiveDirectoryId.GetValue() != "" {
		args.ActiveDirectoryId = pulumi.StringPtr(spec.ActiveDirectoryId.GetValue())
	}

	// Self-managed Active Directory: the domain join uses EITHER the Secrets
	// Manager service-account secret (recommended — credentials never touch
	// the manifest) OR direct username/password; the spec's CEL rules keep the
	// two arms exclusive, so exactly one is wired here.
	if spec.SelfManagedActiveDirectory != nil {
		smad := spec.SelfManagedActiveDirectory

		dnsIps := make(pulumi.StringArray, 0, len(smad.DnsIps))
		for _, ip := range smad.DnsIps {
			dnsIps = append(dnsIps, pulumi.String(ip))
		}

		smadArgs := &fsx.WindowsFileSystemSelfManagedActiveDirectoryArgs{
			DomainName: pulumi.String(smad.DomainName),
			DnsIps:     dnsIps,
		}

		// Each credential field is guarded independently (matching the
		// Terraform module) rather than branched either/or: the secret ARN is
		// a StringValueOrRef resolved at deploy time, so a reference that
		// resolves empty must not force empty username/password strings into
		// the provider's length validation.
		if smad.DomainJoinServiceAccountSecretArn != nil && smad.DomainJoinServiceAccountSecretArn.GetValue() != "" {
			smadArgs.DomainJoinServiceAccountSecret = pulumi.StringPtr(smad.DomainJoinServiceAccountSecretArn.GetValue())
		}
		if smad.Username != "" {
			smadArgs.Username = pulumi.StringPtr(smad.Username)
		}
		if smad.Password != "" {
			smadArgs.Password = pulumi.StringPtr(smad.Password)
		}

		// File system administrators group.
		if smad.GetFileSystemAdministratorsGroup() != "" {
			smadArgs.FileSystemAdministratorsGroup = pulumi.StringPtr(smad.GetFileSystemAdministratorsGroup())
		}

		// Organizational unit distinguished name.
		if smad.OrganizationalUnitDistinguishedName != "" {
			smadArgs.OrganizationalUnitDistinguishedName = pulumi.StringPtr(smad.OrganizationalUnitDistinguishedName)
		}

		args.SelfManagedActiveDirectory = smadArgs
	}

	// DNS aliases.
	if len(spec.Aliases) > 0 {
		aliases := make(pulumi.StringArray, 0, len(spec.Aliases))
		for _, a := range spec.Aliases {
			aliases = append(aliases, pulumi.String(a))
		}
		args.Aliases = aliases
	}

	// Audit log configuration.
	if spec.AuditLogConfiguration != nil {
		auditArgs := &fsx.WindowsFileSystemAuditLogConfigurationArgs{}

		if spec.AuditLogConfiguration.GetFileAccessAuditLogLevel() != "" {
			auditArgs.FileAccessAuditLogLevel = pulumi.StringPtr(spec.AuditLogConfiguration.GetFileAccessAuditLogLevel())
		}

		if spec.AuditLogConfiguration.GetFileShareAccessAuditLogLevel() != "" {
			auditArgs.FileShareAccessAuditLogLevel = pulumi.StringPtr(spec.AuditLogConfiguration.GetFileShareAccessAuditLogLevel())
		}

		if spec.AuditLogConfiguration.AuditLogDestination != nil && spec.AuditLogConfiguration.AuditLogDestination.GetValue() != "" {
			auditArgs.AuditLogDestination = pulumi.StringPtr(spec.AuditLogConfiguration.AuditLogDestination.GetValue())
		}

		args.AuditLogConfiguration = auditArgs
	}

	// Disk IOPS configuration.
	if spec.DiskIopsConfiguration != nil {
		iopsArgs := &fsx.WindowsFileSystemDiskIopsConfigurationArgs{}
		if spec.DiskIopsConfiguration.GetMode() != "" {
			iopsArgs.Mode = pulumi.StringPtr(spec.DiskIopsConfiguration.GetMode())
		}
		if spec.DiskIopsConfiguration.Iops != nil {
			iopsArgs.Iops = pulumi.IntPtr(int(spec.DiskIopsConfiguration.GetIops()))
		}
		args.DiskIopsConfiguration = iopsArgs
	}

	// Automatic backups. Zero is a real value ("no automatic backups") and
	// the spec default (7) diverges from nothing — the resolved value is
	// always sent explicitly.
	if spec.AutomaticBackupRetentionDays != nil {
		args.AutomaticBackupRetentionDays = pulumi.IntPtr(int(spec.GetAutomaticBackupRetentionDays()))
	}

	// Daily automatic backup start time (HH:MM format).
	if spec.DailyAutomaticBackupStartTime != "" {
		args.DailyAutomaticBackupStartTime = pulumi.StringPtr(spec.DailyAutomaticBackupStartTime)
	}

	// Copy tags to backups.
	if spec.CopyTagsToBackups {
		args.CopyTagsToBackups = pulumi.BoolPtr(true)
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

	// Weekly maintenance start time (d:HH:MM format).
	if spec.WeeklyMaintenanceStartTime != "" {
		args.WeeklyMaintenanceStartTime = pulumi.StringPtr(spec.WeeklyMaintenanceStartTime)
	}

	createdFs, err := fsx.NewWindowsFileSystem(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fsx windows file system")
	}

	return createdFs, nil
}

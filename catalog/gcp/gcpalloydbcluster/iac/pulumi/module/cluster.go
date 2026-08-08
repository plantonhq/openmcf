package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/alloydb"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func cluster(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) (*alloydb.Cluster, error) {
	spec := locals.GcpAlloydbCluster.Spec

	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("alloydb.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"alloydb-alloydb.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable alloydb.googleapis.com api")
	}

	args := &alloydb.ClusterArgs{
		ClusterId: pulumi.String(spec.ClusterName),
		Location:  pulumi.String(spec.Location),
		Labels:    pulumi.ToStringMap(locals.GcpLabels),
		// Explicitly false: the provider's client-side deletion_protection
		// flag defaults TRUE and blocks destroy. Both engines disable it so
		// destroy semantics stay identical; deletion safety belongs to the
		// platform's lifecycle layer, not a provider-local flag.
		DeletionProtection: pulumi.BoolPtr(false),
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	if spec.GetSkipAwaitMajorVersionUpgrade() {
		args.SkipAwaitMajorVersionUpgrade = pulumi.BoolPtr(true)
	}

	if spec.Network.GetValue() != "" {
		networkConfig := &alloydb.ClusterNetworkConfigArgs{
			Network: pulumi.StringPtr(spec.Network.GetValue()),
		}
		if spec.AllocatedIpRange != "" {
			networkConfig.AllocatedIpRange = pulumi.StringPtr(spec.AllocatedIpRange)
		}
		args.NetworkConfig = networkConfig
	}

	if spec.PscConfig != nil && spec.PscConfig.PscEnabled {
		args.PscConfig = &alloydb.ClusterPscConfigArgs{
			PscEnabled: pulumi.BoolPtr(true),
		}
	}

	if spec.GetClusterType() != "" {
		args.ClusterType = pulumi.StringPtr(spec.GetClusterType())
	}

	if spec.SecondaryConfig != nil && spec.SecondaryConfig.PrimaryClusterName != "" {
		args.SecondaryConfig = &alloydb.ClusterSecondaryConfigArgs{
			PrimaryClusterName: pulumi.String(spec.SecondaryConfig.PrimaryClusterName),
		}
	}

	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(spec.Annotations)
	}

	if spec.SubscriptionType != "" {
		args.SubscriptionType = pulumi.StringPtr(spec.SubscriptionType)
	}

	if spec.DatabaseVersion != "" {
		args.DatabaseVersion = pulumi.StringPtr(spec.DatabaseVersion)
	}

	if spec.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(spec.DisplayName)
	}

	if spec.InitialUser != nil {
		initialUserArgs := &alloydb.ClusterInitialUserArgs{
			Password: pulumi.String(spec.InitialUser.Password),
		}
		if spec.InitialUser.User != "" {
			initialUserArgs.User = pulumi.StringPtr(spec.InitialUser.User)
		}
		args.InitialUser = initialUserArgs
	}

	if spec.AutomatedBackupPolicy != nil {
		backupArgs := &alloydb.ClusterAutomatedBackupPolicyArgs{}

		if spec.AutomatedBackupPolicy.Enabled {
			backupArgs.Enabled = pulumi.BoolPtr(true)
		}
		if spec.AutomatedBackupPolicy.BackupWindow != "" {
			backupArgs.BackupWindow = pulumi.StringPtr(spec.AutomatedBackupPolicy.BackupWindow)
		}
		if spec.AutomatedBackupPolicy.Location != "" {
			backupArgs.Location = pulumi.StringPtr(spec.AutomatedBackupPolicy.Location)
		}

		if spec.AutomatedBackupPolicy.QuantityBasedRetentionCount > 0 {
			backupArgs.QuantityBasedRetention = &alloydb.ClusterAutomatedBackupPolicyQuantityBasedRetentionArgs{
				Count: pulumi.IntPtr(int(spec.AutomatedBackupPolicy.QuantityBasedRetentionCount)),
			}
		}
		if spec.AutomatedBackupPolicy.TimeBasedRetentionPeriod != "" {
			backupArgs.TimeBasedRetention = &alloydb.ClusterAutomatedBackupPolicyTimeBasedRetentionArgs{
				RetentionPeriod: pulumi.StringPtr(spec.AutomatedBackupPolicy.TimeBasedRetentionPeriod),
			}
		}

		if spec.AutomatedBackupPolicy.WeeklySchedule != nil {
			schedule := spec.AutomatedBackupPolicy.WeeklySchedule
			scheduleArgs := &alloydb.ClusterAutomatedBackupPolicyWeeklyScheduleArgs{}

			if len(schedule.DaysOfWeek) > 0 {
				scheduleArgs.DaysOfWeeks = pulumi.ToStringArray(schedule.DaysOfWeek)
			}

			scheduleArgs.StartTimes = alloydb.ClusterAutomatedBackupPolicyWeeklyScheduleStartTimeArray{
				&alloydb.ClusterAutomatedBackupPolicyWeeklyScheduleStartTimeArgs{
					Hours: pulumi.IntPtr(int(schedule.StartHour)),
				},
			}

			backupArgs.WeeklySchedule = scheduleArgs
		}

		if spec.AutomatedBackupPolicy.EncryptionKmsKeyName != nil && spec.AutomatedBackupPolicy.EncryptionKmsKeyName.GetValue() != "" {
			backupArgs.EncryptionConfig = &alloydb.ClusterAutomatedBackupPolicyEncryptionConfigArgs{
				KmsKeyName: pulumi.StringPtr(spec.AutomatedBackupPolicy.EncryptionKmsKeyName.GetValue()),
			}
		}

		args.AutomatedBackupPolicy = backupArgs
	}

	if spec.ContinuousBackupConfig != nil {
		continuousArgs := &alloydb.ClusterContinuousBackupConfigArgs{
			Enabled: pulumi.BoolPtr(spec.ContinuousBackupConfig.Enabled),
		}

		if spec.ContinuousBackupConfig.RecoveryWindowDays > 0 {
			continuousArgs.RecoveryWindowDays = pulumi.IntPtr(int(spec.ContinuousBackupConfig.RecoveryWindowDays))
		}

		if spec.ContinuousBackupConfig.EncryptionKmsKeyName != nil && spec.ContinuousBackupConfig.EncryptionKmsKeyName.GetValue() != "" {
			continuousArgs.EncryptionConfig = &alloydb.ClusterContinuousBackupConfigEncryptionConfigArgs{
				KmsKeyName: pulumi.StringPtr(spec.ContinuousBackupConfig.EncryptionKmsKeyName.GetValue()),
			}
		}

		args.ContinuousBackupConfig = continuousArgs
	}

	if spec.KmsKeyName != nil && spec.KmsKeyName.GetValue() != "" {
		args.EncryptionConfig = &alloydb.ClusterEncryptionConfigArgs{
			KmsKeyName: pulumi.StringPtr(spec.KmsKeyName.GetValue()),
		}
	}

	if spec.MaintenanceWindow != nil {
		args.MaintenanceUpdatePolicy = &alloydb.ClusterMaintenanceUpdatePolicyArgs{
			MaintenanceWindows: alloydb.ClusterMaintenanceUpdatePolicyMaintenanceWindowArray{
				&alloydb.ClusterMaintenanceUpdatePolicyMaintenanceWindowArgs{
					Day: pulumi.String(spec.MaintenanceWindow.Day),
					StartTime: &alloydb.ClusterMaintenanceUpdatePolicyMaintenanceWindowStartTimeArgs{
						Hours: pulumi.Int(int(spec.MaintenanceWindow.StartHour)),
					},
				},
			},
		}
	}

	createdCluster, err := alloydb.NewCluster(ctx, "alloydb-cluster", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create alloydb cluster")
	}

	ctx.Export(OpClusterId, createdCluster.Name)
	ctx.Export(OpClusterName, pulumi.String(spec.ClusterName))
	ctx.Export(OpDatabaseVersion, createdCluster.DatabaseVersion)
	ctx.Export(OpState, createdCluster.State)

	return createdCluster, nil
}

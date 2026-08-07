package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/spanner"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// spannerBackupSchedule provisions the Spanner backup schedule — backups of
// one database on a cron cadence, each retained for retention_duration. A
// database commonly carries a daily incremental schedule AND a weekly full
// schedule side by side. name, instance, database, and backup type are
// immutable; cron, retention, and encryption update in place.
func spannerBackupSchedule(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpSpannerBackupSchedule.Spec

	// Enable the Spanner API so a fresh project can host backup schedules.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("spanner.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"spanner-spanner.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable spanner.googleapis.com api")
	}

	args := &spanner.BackupScheduleArgs{
		Name:     pulumi.StringPtr(locals.ScheduleName),
		Instance: pulumi.String(spec.Instance.GetValue()),
		Database: pulumi.String(spec.Database.GetValue()),

		// Applies to backups created AFTER a change; existing backups keep
		// the retention they were created with.
		RetentionDuration: pulumi.String(spec.RetentionDuration),

		Spec: &spanner.BackupScheduleSpecArgs{
			CronSpec: &spanner.BackupScheduleSpecCronSpecArgs{
				// Evaluated in UTC. Spanner accepts a bounded set of
				// frequencies: every 12 hours, daily, weekly, or monthly.
				Text: pulumi.StringPtr(spec.Cron),
			},
		},
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Exactly one backup kind, expressed by the provider as a pair of empty
	// marker blocks. INCREMENTAL chains store only changes since the
	// previous backup (cheaper storage, same restore semantics) and require
	// the instance to be ENTERPRISE or ENTERPRISE_PLUS edition. The spec
	// default is FULL.
	if spec.GetBackupType() == "INCREMENTAL" {
		args.IncrementalBackupSpec = &spanner.BackupScheduleIncrementalBackupSpecArgs{}
	} else {
		args.FullBackupSpec = &spanner.BackupScheduleFullBackupSpecArgs{}
	}

	// If omitted, backups use USE_DATABASE_ENCRYPTION (inherit the
	// database's posture). CMEK requires exactly one key shape — enforced
	// pre-deploy by spec validation.
	if spec.EncryptionConfig != nil {
		encryptionArgs := &spanner.BackupScheduleEncryptionConfigArgs{
			EncryptionType: pulumi.String(spec.EncryptionConfig.EncryptionType),
		}
		if spec.EncryptionConfig.KmsKeyName.GetValue() != "" {
			encryptionArgs.KmsKeyName = pulumi.StringPtr(spec.EncryptionConfig.KmsKeyName.GetValue())
		}
		if len(spec.EncryptionConfig.KmsKeyNames) > 0 {
			keyNames := make([]string, 0, len(spec.EncryptionConfig.KmsKeyNames))
			for _, keyName := range spec.EncryptionConfig.KmsKeyNames {
				keyNames = append(keyNames, keyName.GetValue())
			}
			encryptionArgs.KmsKeyNames = pulumi.ToStringArray(keyNames)
		}
		args.EncryptionConfig = encryptionArgs
	}

	createdSchedule, err := spanner.NewBackupSchedule(ctx, "spanner-backup-schedule", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create spanner backup schedule")
	}

	// schedule_id is built from the created resource's resolved attributes
	// so the output is correct under the ambient-project fallback (the spec
	// project may be empty).
	ctx.Export(OpScheduleId, pulumi.Sprintf(
		"projects/%s/instances/%s/databases/%s/backupSchedules/%s",
		createdSchedule.Project,
		createdSchedule.Instance,
		createdSchedule.Database,
		createdSchedule.Name,
	))
	ctx.Export(OpScheduleName, createdSchedule.Name)

	return nil
}

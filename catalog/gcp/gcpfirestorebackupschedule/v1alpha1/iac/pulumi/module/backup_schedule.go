package module

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// firestoreBackupSchedule provisions a Firestore backup schedule — periodic
// managed backups with retention, distinct from point-in-time recovery.
// A database supports one daily and one weekly schedule; the
// daily-plus-weekly pattern is two of these resources on the same database.
// Recurrence (daily or weekly day) is immutable; retention updates in place.
// Backups already taken outlive the schedule — deleting this resource stops
// future backups but never deletes existing ones.
func firestoreBackupSchedule(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpFirestoreBackupSchedule.Spec

	// Enable the Firestore API — backup schedules are managed through the
	// Firestore Admin API. disable_on_destroy stays false: tearing down one
	// schedule must never disable the API for everything else in the project.
	firestoreApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("firestore.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		firestoreApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdFirestoreApi, err := projects.NewService(ctx,
		"fsbs-firestore.googleapis.com", firestoreApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable firestore.googleapis.com api")
	}

	args := &firestore.BackupScheduleArgs{
		Database:  pulumi.String(spec.Database.GetValue()),
		Retention: pulumi.String(spec.Retention),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Exactly one recurrence shape — expressed by the provider as a pair of
	// marker blocks (daily empty struct, or weekly with a day).
	if spec.Daily {
		args.DailyRecurrence = &firestore.BackupScheduleDailyRecurrenceArgs{}
	} else if spec.WeeklyRecurrence != nil {
		args.WeeklyRecurrence = &firestore.BackupScheduleWeeklyRecurrenceArgs{
			Day: pulumi.StringPtr(spec.WeeklyRecurrence.Day),
		}
	}

	createdSchedule, err := firestore.NewBackupSchedule(ctx, "firestore-backup-schedule", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdFirestoreApi}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create firestore backup schedule")
	}

	// schedule_id is the last path segment of the server-assigned name —
	// what Admin API calls address the schedule by.
	ctx.Export(OpScheduleId, createdSchedule.Name.ApplyT(func(name string) string {
		parts := strings.Split(name, "/")
		return parts[len(parts)-1]
	}))
	ctx.Export(OpDatabase, pulumi.String(spec.Database.GetValue()))

	return nil
}

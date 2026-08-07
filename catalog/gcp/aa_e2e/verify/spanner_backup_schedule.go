package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// spannerBackupScheduleVerifier probes a Cloud Spanner backup schedule via
// the spanner admin API. Posture assertions confirm the schedule carries a
// cron cadence, a retention duration, and exactly one backup-kind spec —
// proof GCP accepted the schedule as configured.
type spannerBackupScheduleVerifier struct{}

func (v *spannerBackupScheduleVerifier) IDOutputKey() string { return "schedule_id" }

func (v *spannerBackupScheduleVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	scheduleID := outputs["schedule_id"]
	if scheduleID == "" {
		return errors.New("schedule_id output missing after deploy")
	}

	schedule, err := svc.Spanner.Projects.Instances.Databases.BackupSchedules.Get(scheduleID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "spanner backup schedule %s not found after deploy", scheduleID)
	}

	if schedule.Spec == nil || schedule.Spec.CronSpec == nil || schedule.Spec.CronSpec.Text == "" {
		return errors.Errorf("spanner backup schedule %s has no cron cadence after deploy", scheduleID)
	}
	if schedule.RetentionDuration == "" {
		return errors.Errorf("spanner backup schedule %s has no retention duration after deploy", scheduleID)
	}
	hasFull := schedule.FullBackupSpec != nil
	hasIncremental := schedule.IncrementalBackupSpec != nil
	if hasFull == hasIncremental {
		return errors.Errorf("spanner backup schedule %s must carry exactly one backup kind (full=%t incremental=%t)",
			scheduleID, hasFull, hasIncremental)
	}
	return nil
}

func (v *spannerBackupScheduleVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	scheduleID := outputs["schedule_id"]
	if scheduleID == "" {
		return nil
	}

	_, err := svc.Spanner.Projects.Instances.Databases.BackupSchedules.Get(scheduleID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("spanner backup schedule %s still exists after destroy", scheduleID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	// When the whole chain (database or instance included) was destroyed, the
	// API answers for the missing parent instead of the schedule.
	if apiErr != nil && apiErr.Code == 400 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing spanner backup schedule %s after destroy", scheduleID)
}

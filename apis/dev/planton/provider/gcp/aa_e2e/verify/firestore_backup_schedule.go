package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// firestoreBackupScheduleVerifier probes a Firestore backup schedule via
// the Firestore Admin API. Schedules are addressed by
// projects/{p}/databases/{db}/backupSchedules/{id} — the verifier
// assembles the path from the schedule_id and database outputs. The
// posture assertion confirms the live schedule carries a retention and
// exactly one recurrence shape.
type firestoreBackupScheduleVerifier struct{}

func (v *firestoreBackupScheduleVerifier) IDOutputKey() string { return "schedule_id" }

func (v *firestoreBackupScheduleVerifier) schedulePath(svc *Services, outputs map[string]string) (string, error) {
	scheduleId := outputs["schedule_id"]
	database := outputs["database"]
	if scheduleId == "" || database == "" {
		return "", errors.New("schedule_id or database output missing")
	}
	return fmt.Sprintf("projects/%s/databases/%s/backupSchedules/%s", svc.Project, database, scheduleId), nil
}

func (v *firestoreBackupScheduleVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.schedulePath(svc, outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	schedule, err := svc.Firestore.Projects.Databases.BackupSchedules.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "firestore backup schedule %s not found after deploy", path)
	}
	if schedule.Retention == "" {
		return errors.Errorf("firestore backup schedule %s carries no retention", path)
	}
	daily := schedule.DailyRecurrence != nil
	weekly := schedule.WeeklyRecurrence != nil
	if daily == weekly {
		return errors.Errorf("firestore backup schedule %s recurrence posture is wrong: daily=%v weekly=%v (want exactly one)",
			path, daily, weekly)
	}
	return nil
}

func (v *firestoreBackupScheduleVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.schedulePath(svc, outputs)
	if err != nil {
		return nil
	}

	_, err = svc.Firestore.Projects.Databases.BackupSchedules.Get(path).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && (apiErr.Code == 404 || apiErr.Code == 403) {
			// 404: schedule gone. 403/NOT_FOUND wrapping: the chain's
			// database is destroyed with (or before) the schedule, and
			// probing a schedule under a deleted database can surface as
			// a permission-shaped error rather than a clean 404.
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing firestore backup schedule %s after destroy", path)
	}
	return errors.Errorf("firestore backup schedule %s still exists after destroy", path)
}

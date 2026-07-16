package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudSchedulerJobVerifier probes a Cloud Scheduler job via the Cloud
// Scheduler API. The job_id output is the fully qualified resource path
// (projects/{p}/locations/{l}/jobs/{name}) the Jobs.Get call consumes
// directly. Posture assertions confirm the job is ENABLED (created
// unpaused), carries a schedule, and has exactly one target attached —
// the live proof of the exactly-one-target contract.
type cloudSchedulerJobVerifier struct{}

func (v *cloudSchedulerJobVerifier) IDOutputKey() string { return "job_id" }

func (v *cloudSchedulerJobVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	jobID := outputs["job_id"]
	if jobID == "" {
		return errors.New("job_id output missing after deploy")
	}

	job, err := svc.CloudScheduler.Projects.Locations.Jobs.Get(jobID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "cloud scheduler job %s not found after deploy", jobID)
	}

	// Scenarios create unpaused jobs; anything but ENABLED means the deploy
	// produced a job that will never fire (PAUSED/DISABLED/UPDATE_FAILED).
	if job.State != "ENABLED" {
		return errors.Errorf("cloud scheduler job %s in state %s after deploy (expected ENABLED)", jobID, job.State)
	}

	if job.Schedule == "" {
		return errors.Errorf("cloud scheduler job %s has no schedule after deploy", jobID)
	}

	// Exactly one target must be attached (the spec's CEL rule mirrored
	// against live state).
	targets := 0
	if job.HttpTarget != nil {
		targets++
	}
	if job.PubsubTarget != nil {
		targets++
	}
	if job.AppEngineHttpTarget != nil {
		targets++
	}
	if targets != 1 {
		return errors.Errorf("cloud scheduler job %s has %d targets after deploy (expected exactly 1)", jobID, targets)
	}
	return nil
}

func (v *cloudSchedulerJobVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	jobID := outputs["job_id"]
	if jobID == "" {
		return nil
	}

	_, err := svc.CloudScheduler.Projects.Locations.Jobs.Get(jobID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("cloud scheduler job %s still exists after destroy", jobID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing cloud scheduler job %s after destroy", jobID)
}

package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudRunJobVerifier probes a Cloud Run job via the Run Admin API v2.
// Posture assertions confirm the job reconciled successfully (the Ready
// terminal condition) and that a task template exists.
type cloudRunJobVerifier struct{}

func (v *cloudRunJobVerifier) IDOutputKey() string { return "job_name" }

func (v *cloudRunJobVerifier) jobPath(svc *Services, outputs map[string]string) (string, error) {
	name := outputs["job_name"]
	location := outputs["location"]
	if name == "" || location == "" {
		return "", errors.New("job_name or location output missing")
	}
	return fmt.Sprintf("projects/%s/locations/%s/jobs/%s", svc.Project, location, name), nil
}

func (v *cloudRunJobVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.jobPath(svc, outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	job, err := svc.Run.Projects.Locations.Jobs.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "cloud run job %s not found after deploy", path)
	}

	if job.TerminalCondition == nil || job.TerminalCondition.State != "CONDITION_SUCCEEDED" {
		state := "<none>"
		if job.TerminalCondition != nil {
			state = job.TerminalCondition.State
		}
		return errors.Errorf("cloud run job %s terminal condition is %q, want CONDITION_SUCCEEDED", path, state)
	}
	if job.Template == nil || job.Template.Template == nil {
		return errors.Errorf("cloud run job %s has no task template after deploy", path)
	}
	return nil
}

func (v *cloudRunJobVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.jobPath(svc, outputs)
	if err != nil {
		return nil
	}

	_, err = svc.Run.Projects.Locations.Jobs.Get(path).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("cloud run job %s still exists after destroy", path)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing cloud run job %s after destroy", path)
}

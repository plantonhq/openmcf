package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// plantonRunnerVerifier probes a Planton runner appliance via the Run
// Admin API v2. The appliance is a single-instance Cloud Run service, so
// the posture assertions confirm the service reconciled successfully (the
// Ready terminal condition — Cloud Run gates creation on first-revision
// readiness, so a runner that could not join its control plane never gets
// here) and that the module's scaling pin held (exactly one always-on
// instance; the singleton law — a second instance joining under the same
// runner name would revoke the first's key).
type plantonRunnerVerifier struct{}

func (v *plantonRunnerVerifier) IDOutputKey() string { return "service_short_name" }

// servicePath builds the projects/{p}/locations/{region}/services/{name}
// resource path the run API addresses services by.
func (v *plantonRunnerVerifier) servicePath(svc *Services, outputs map[string]string) (string, error) {
	name := outputs["service_short_name"]
	region := outputs["region"]
	if name == "" || region == "" {
		return "", errors.New("service_short_name or region output missing")
	}
	project := outputs["project_id"]
	if project == "" {
		project = svc.Project
	}
	return fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, name), nil
}

func (v *plantonRunnerVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.servicePath(svc, outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	service, err := svc.Run.Projects.Locations.Services.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "planton runner service %s not found after deploy", path)
	}

	// The Ready terminal condition is the API's own "reconciled and
	// serving" signal — for the runner it additionally proves the join
	// succeeded, because the container exits on a refused join and the
	// first revision would never turn ready.
	if service.TerminalCondition == nil || service.TerminalCondition.State != "CONDITION_SUCCEEDED" {
		state := "<none>"
		if service.TerminalCondition != nil {
			state = service.TerminalCondition.State
		}
		return errors.Errorf("planton runner service %s terminal condition is %q, want CONDITION_SUCCEEDED — a refused join (bad or unreachable token) keeps the first revision from turning ready", path, state)
	}

	// The singleton law, straight off the live template.
	if service.Template == nil || service.Template.Scaling == nil ||
		service.Template.Scaling.MinInstanceCount != 1 || service.Template.Scaling.MaxInstanceCount != 1 {
		return errors.Errorf("planton runner service %s must pin scaling to exactly one instance (a second instance joining under the same runner name would revoke the first's key)", path)
	}

	return nil
}

func (v *plantonRunnerVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.servicePath(svc, outputs)
	if err != nil {
		// Without a path there is nothing left to probe; treat as gone.
		return nil
	}

	_, err = svc.Run.Projects.Locations.Services.Get(path).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("planton runner service %s still exists after destroy", path)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing planton runner service %s after destroy", path)
}

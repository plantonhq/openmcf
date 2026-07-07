package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// composerEnvironmentVerifier probes a Cloud Composer environment via
// the Composer API. The environment_id output is the fully qualified
// resource name (projects/{p}/locations/{r}/environments/{e}) the
// environments.get call is addressed with. Posture assertions confirm
// the environment reached RUNNING and reports the Airflow web UI and
// DAG bucket — proof the whole managed stack (GKE cluster, database,
// web server, bucket) actually assembled.
type composerEnvironmentVerifier struct{}

func (v *composerEnvironmentVerifier) IDOutputKey() string { return "environment_id" }

func (v *composerEnvironmentVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	environmentID := outputs["environment_id"]
	if environmentID == "" {
		return errors.New("environment_id output missing after deploy")
	}

	env, err := svc.Composer.Projects.Locations.Environments.Get(environmentID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "composer environment %s not found after deploy", environmentID)
	}

	if env.State != "RUNNING" {
		return errors.Errorf("composer environment %s is %s after deploy, expected RUNNING", environmentID, env.State)
	}
	if env.Config == nil || env.Config.AirflowUri == "" {
		return errors.Errorf("composer environment %s reports no Airflow web UI after deploy", environmentID)
	}
	if env.Config.DagGcsPrefix == "" {
		return errors.Errorf("composer environment %s reports no DAG bucket after deploy", environmentID)
	}
	return nil
}

func (v *composerEnvironmentVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	environmentID := outputs["environment_id"]
	if environmentID == "" {
		return nil
	}

	env, err := svc.Composer.Projects.Locations.Environments.Get(environmentID).Context(ctx).Do()
	if err == nil {
		// An environment mid-teardown reports DELETING — that is the
		// destroyed path, not an orphan.
		if env.State == "DELETING" {
			return nil
		}
		return errors.Errorf("composer environment %s still exists after destroy", environmentID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing composer environment %s after destroy", environmentID)
}

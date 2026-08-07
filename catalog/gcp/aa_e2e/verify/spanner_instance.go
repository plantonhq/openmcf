package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// spannerInstanceVerifier probes a Cloud Spanner instance via the spanner
// admin API. Posture assertions confirm the instance reached READY, its
// config matches the config output, and — when autoscaling is off — GCP
// allocated real compute capacity.
type spannerInstanceVerifier struct{}

func (v *spannerInstanceVerifier) IDOutputKey() string { return "instance_id" }

func (v *spannerInstanceVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	instanceID := outputs["instance_id"]
	if instanceID == "" {
		return errors.New("instance_id output missing after deploy")
	}

	instance, err := svc.Spanner.Projects.Instances.Get(instanceID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "spanner instance %s not found after deploy", instanceID)
	}

	if instance.State != "READY" {
		return errors.Errorf("spanner instance %s state is %q, want READY", instanceID, instance.State)
	}

	// The config output is the plain config name (e.g. regional-us-central1);
	// the API returns the fully qualified projects/{p}/instanceConfigs/{name}.
	if wantConfig := outputs["config"]; wantConfig != "" && pathTail(instance.Config) != wantConfig {
		return errors.Errorf("spanner instance %s config mismatch: output %q, live %q",
			instanceID, wantConfig, pathTail(instance.Config))
	}

	// Whatever capacity mode the scenario chose (nodes, processing units, or
	// autoscaling), a READY instance must have real allocated capacity.
	if instance.NodeCount == 0 && instance.ProcessingUnits == 0 {
		return errors.Errorf("spanner instance %s has zero allocated capacity after deploy", instanceID)
	}
	return nil
}

func (v *spannerInstanceVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	instanceID := outputs["instance_id"]
	if instanceID == "" {
		return nil
	}

	_, err := svc.Spanner.Projects.Instances.Get(instanceID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing spanner instance %s after destroy", instanceID)
	}
	return errors.Errorf("spanner instance %s still exists after destroy", instanceID)
}

// pathTail returns the last segment of a slash-separated resource path.
func pathTail(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

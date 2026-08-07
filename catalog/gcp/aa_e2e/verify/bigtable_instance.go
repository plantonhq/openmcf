package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// bigtableInstanceVerifier probes a Bigtable instance via the Bigtable
// Admin API. Posture assertions confirm the instance is READY and that at
// least one cluster is serving — an instance without clusters cannot
// store or serve data.
type bigtableInstanceVerifier struct{}

func (v *bigtableInstanceVerifier) IDOutputKey() string { return "instance_id" }

func (v *bigtableInstanceVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	instanceId := outputs["instance_id"]
	if instanceId == "" {
		return errors.New("instance_id output missing after deploy")
	}

	instance, err := svc.BigtableAdmin.Projects.Instances.Get(instanceId).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "bigtable instance %s not found after deploy", instanceId)
	}
	if instance.State != "READY" {
		return errors.Errorf("bigtable instance %s state is %q, want READY", instanceId, instance.State)
	}

	clusters, err := svc.BigtableAdmin.Projects.Instances.Clusters.List(instanceId).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "failed to list clusters for bigtable instance %s", instanceId)
	}
	if len(clusters.Clusters) == 0 {
		return errors.Errorf("bigtable instance %s has no clusters — it cannot store or serve data", instanceId)
	}
	for _, c := range clusters.Clusters {
		if c.State != "READY" {
			return errors.Errorf("bigtable instance %s cluster %s state is %q, want READY", instanceId, c.Name, c.State)
		}
	}
	return nil
}

func (v *bigtableInstanceVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	instanceId := outputs["instance_id"]
	if instanceId == "" {
		return nil
	}

	_, err := svc.BigtableAdmin.Projects.Instances.Get(instanceId).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing bigtable instance %s after destroy", instanceId)
	}
	return errors.Errorf("bigtable instance %s still exists after destroy", instanceId)
}

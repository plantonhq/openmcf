package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// alloydbInstanceVerifier probes an AlloyDB instance via the alloydb API.
// Posture assertions confirm the instance reconciled to READY.
type alloydbInstanceVerifier struct{}

func (v *alloydbInstanceVerifier) IDOutputKey() string { return "instance_name" }

func (v *alloydbInstanceVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["instance_name"]
	if name == "" {
		return errors.New("instance_name output missing after deploy")
	}

	inst, err := svc.AlloyDB.Projects.Locations.Clusters.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "alloydb instance %s not found after deploy", name)
	}
	if inst.State != "READY" {
		return errors.Errorf("alloydb instance %s state is %q, want READY", name, inst.State)
	}
	return nil
}

func (v *alloydbInstanceVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["instance_name"]
	if name == "" {
		return nil
	}

	_, err := svc.AlloyDB.Projects.Locations.Clusters.Instances.Get(name).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("alloydb instance %s still exists after destroy", name)
	}
	if isAlloyDBNotFound(err) {
		return nil
	}
	if clusterID := alloyDBClusterFromChildPath(name, "/instances/"); clusterID != "" {
		if _, clusterErr := svc.AlloyDB.Projects.Locations.Clusters.Get(clusterID).Context(ctx).Do(); clusterErr != nil && isAlloyDBNotFound(clusterErr) {
			return nil
		}
	}
	return errors.Wrapf(err, "unexpected error probing alloydb instance %s after destroy", name)
}

func alloyDBClusterFromChildPath(resourcePath, childMarker string) string {
	if idx := strings.Index(resourcePath, childMarker); idx >= 0 {
		return resourcePath[:idx]
	}
	return ""
}

func isAlloyDBNotFound(err error) bool {
	var apiErr *googleapi.Error
	return errors.As(err, &apiErr) && apiErr.Code == 404
}

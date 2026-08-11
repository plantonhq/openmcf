package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// computeMigVerifier probes a managed instance group by the group
// manager's name and location outputs. One kind maps to two API
// collections (zonal and regional instance group managers), so the
// verifier routes on the location's shape — a zone's final segment is a
// single letter ("us-central1-a"); a region's is not ("us-central1") —
// rather than requiring the caller to know the scope. Outputs-driven and
// scenario-agnostic: every scenario exposes the same output keys.
type computeMigVerifier struct{}

func (v *computeMigVerifier) IDOutputKey() string { return "self_link" }

// locationIsZone reports whether the location output names a zone
// (region + one-letter zone suffix) rather than a region.
func locationIsZone(location string) bool {
	segments := strings.Split(location, "-")
	return len(segments) > 0 && len(segments[len(segments)-1]) == 1
}

func (v *computeMigVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["mig_name"]
	location := outputs["location"]

	if locationIsZone(location) {
		manager, err := svc.Compute.InstanceGroupManagers.Get(svc.Project, location, name).Context(ctx).Do()
		if err != nil {
			return errors.Wrapf(err, "zonal instance group manager %s/%s not found after deploy", location, name)
		}
		if manager.InstanceGroup == "" {
			return errors.Errorf("zonal instance group manager %s reports no instance group", name)
		}
		return nil
	}
	manager, err := svc.Compute.RegionInstanceGroupManagers.Get(svc.Project, location, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "regional instance group manager %s/%s not found after deploy", location, name)
	}
	if manager.InstanceGroup == "" {
		return errors.Errorf("regional instance group manager %s reports no instance group", name)
	}
	return nil
}

func (v *computeMigVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["mig_name"]
	location := outputs["location"]

	var err error
	if locationIsZone(location) {
		_, err = svc.Compute.InstanceGroupManagers.Get(svc.Project, location, name).Context(ctx).Do()
	} else {
		_, err = svc.Compute.RegionInstanceGroupManagers.Get(svc.Project, location, name).Context(ctx).Do()
	}
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing instance group manager %s after destroy", name)
	}
	return errors.Errorf("instance group manager %s still exists after destroy", strings.TrimSpace(name))
}

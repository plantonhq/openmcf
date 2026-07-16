package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// gkeClusterVerifier probes a GKE cluster via the container API. The
// cluster_id output is already the fully qualified resource path
// (projects/{p}/locations/{l}/clusters/{n}) the API addresses clusters by.
// Posture assertions confirm the control plane is RUNNING, the endpoint and
// name outputs match live state, and — when the outputs claim a Workload
// Identity pool — that the live cluster actually carries it.
type gkeClusterVerifier struct{}

func (v *gkeClusterVerifier) IDOutputKey() string { return "cluster_id" }

func (v *gkeClusterVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path := outputs["cluster_id"]
	if path == "" {
		return errors.New("cluster_id output missing after deploy")
	}

	cluster, err := svc.Container.Projects.Locations.Clusters.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "gke cluster %s not found after deploy", path)
	}
	if cluster.Status != "RUNNING" {
		return errors.Errorf("gke cluster %s status is %q, want RUNNING", path, cluster.Status)
	}
	if name := outputs["name"]; name != "" && cluster.Name != name {
		return errors.Errorf("gke cluster %s name mismatch: output %q, live %q", path, name, cluster.Name)
	}
	if endpoint := outputs["endpoint"]; endpoint != "" && cluster.Endpoint != endpoint {
		return errors.Errorf("gke cluster %s endpoint mismatch: output %q, live %q", path, endpoint, cluster.Endpoint)
	}
	// Assert output → live state: a claimed Workload Identity pool must exist
	// on the cluster. (The reverse — pool live but output empty — is legal
	// only for Autopilot-managed defaults, so it is not asserted.)
	if pool := outputs["workload_identity_pool"]; pool != "" {
		if cluster.WorkloadIdentityConfig == nil || cluster.WorkloadIdentityConfig.WorkloadPool != pool {
			live := ""
			if cluster.WorkloadIdentityConfig != nil {
				live = cluster.WorkloadIdentityConfig.WorkloadPool
			}
			return errors.Errorf("gke cluster %s workload identity pool mismatch: output %q, live %q", path, pool, live)
		}
	}
	return nil
}

func (v *gkeClusterVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path := outputs["cluster_id"]
	if path == "" {
		return nil
	}

	cluster, err := svc.Container.Projects.Locations.Clusters.Get(path).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing gke cluster %s after destroy", path)
	}
	return errors.Errorf("gke cluster %s still exists after destroy (status %s)", path, cluster.Status)
}

package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// gkeNodePoolVerifier probes a GKE node pool via the container API. The
// node_pool_id output is already the fully qualified resource path
// (projects/{p}/locations/{l}/clusters/{c}/nodePools/{n}) the API addresses
// pools by. Posture assertions confirm the pool is RUNNING, the name output
// matches live state, the pool is backed by real instance groups, and —
// when the outputs claim autoscaling bounds — that the live pool actually
// carries an enabled autoscaler.
type gkeNodePoolVerifier struct{}

func (v *gkeNodePoolVerifier) IDOutputKey() string { return "node_pool_id" }

func (v *gkeNodePoolVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path := outputs["node_pool_id"]
	if path == "" {
		return errors.New("node_pool_id output missing after deploy")
	}

	nodePool, err := svc.Container.Projects.Locations.Clusters.NodePools.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "gke node pool %s not found after deploy", path)
	}
	if nodePool.Status != "RUNNING" {
		return errors.Errorf("gke node pool %s status is %q, want RUNNING", path, nodePool.Status)
	}
	if name := outputs["node_pool_name"]; name != "" && nodePool.Name != name {
		return errors.Errorf("gke node pool %s name mismatch: output %q, live %q", path, name, nodePool.Name)
	}
	if len(nodePool.InstanceGroupUrls) == 0 {
		return errors.Errorf("gke node pool %s has no backing instance groups", path)
	}
	// Assert output → live state: outputs claiming an autoscaling maximum
	// must correspond to a live, enabled autoscaler. (min_nodes/max_nodes
	// mirror the fixed node_count for unmanaged pools, where max equals the
	// count and no autoscaler exists — so only a max>count posture without
	// an autoscaler would be dishonest; the scenario manifests keep the two
	// modes distinct, and this assertion runs only when autoscaling is the
	// claimed mode via the min!=max signal.)
	if minOut, maxOut := outputs["min_nodes"], outputs["max_nodes"]; minOut != "" && maxOut != "" && minOut != maxOut {
		if nodePool.Autoscaling == nil || !nodePool.Autoscaling.Enabled {
			return errors.Errorf("gke node pool %s outputs claim autoscaling bounds [%s..%s] but the live pool has no enabled autoscaler", path, minOut, maxOut)
		}
	}
	return nil
}

func (v *gkeNodePoolVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path := outputs["node_pool_id"]
	if path == "" {
		return nil
	}

	nodePool, err := svc.Container.Projects.Locations.Clusters.NodePools.Get(path).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && (apiErr.Code == 404 ||
			// Destroying the pool's own cluster (a torn-down prerequisite)
			// removes the pool with it; the API answers 400/404 for pools
			// addressed under a deleted cluster depending on timing.
			apiErr.Code == 400) {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing gke node pool %s after destroy", path)
	}
	return errors.Errorf("gke node pool %s still exists after destroy (status %s)", path, nodePool.Status)
}

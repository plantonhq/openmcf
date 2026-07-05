package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// alloydbClusterVerifier probes an AlloyDB cluster via the alloydb API.
// Posture assertions confirm the cluster reconciled to READY and that the
// cluster_name output matches the live resource path tail.
type alloydbClusterVerifier struct{}

func (v *alloydbClusterVerifier) IDOutputKey() string { return "cluster_id" }

func (v *alloydbClusterVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	clusterID := outputs["cluster_id"]
	if clusterID == "" {
		return errors.New("cluster_id output missing after deploy")
	}

	cluster, err := svc.AlloyDB.Projects.Locations.Clusters.Get(clusterID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "alloydb cluster %s not found after deploy", clusterID)
	}
	if cluster.State != "READY" {
		return errors.Errorf("alloydb cluster %s state is %q, want READY", clusterID, cluster.State)
	}

	if wantName := outputs["cluster_name"]; wantName != "" {
		gotName := clusterPathTail(clusterID)
		if gotName != wantName {
			return errors.Errorf("alloydb cluster %s cluster_name mismatch: output %q, live path tail %q", clusterID, wantName, gotName)
		}
	}
	return nil
}

func (v *alloydbClusterVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	clusterID := outputs["cluster_id"]
	if clusterID == "" {
		return nil
	}

	_, err := svc.AlloyDB.Projects.Locations.Clusters.Get(clusterID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing alloydb cluster %s after destroy", clusterID)
	}
	return errors.Errorf("alloydb cluster %s still exists after destroy", clusterID)
}

// clusterPathTail returns the cluster_id segment from a fully qualified
// projects/.../locations/.../clusters/{id} path.
func clusterPathTail(clusterID string) string {
	const marker = "/clusters/"
	if idx := strings.LastIndex(clusterID, marker); idx >= 0 {
		return clusterID[idx+len(marker):]
	}
	return clusterID
}

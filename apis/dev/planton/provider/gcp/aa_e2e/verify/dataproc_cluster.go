package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// dataprocClusterVerifier probes a Dataproc cluster via the Dataproc
// API. The cluster_id output is the fully qualified resource name
// (projects/{p}/regions/{r}/clusters/{c}) — the compound identifier the
// clusters.get call is addressed with. Posture assertions confirm the
// cluster reached RUNNING and reports a staging bucket — proof the
// control plane fully assembled it, not just accepted the request.
type dataprocClusterVerifier struct{}

func (v *dataprocClusterVerifier) IDOutputKey() string { return "cluster_id" }

// parseClusterID splits projects/{p}/regions/{r}/clusters/{c}.
func parseClusterID(clusterID string) (project, region, name string, err error) {
	parts := strings.Split(clusterID, "/")
	if len(parts) != 6 || parts[0] != "projects" || parts[2] != "regions" || parts[4] != "clusters" {
		return "", "", "", errors.Errorf("cluster_id %q is not projects/{p}/regions/{r}/clusters/{c}", clusterID)
	}
	return parts[1], parts[3], parts[5], nil
}

func (v *dataprocClusterVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	clusterID := outputs["cluster_id"]
	if clusterID == "" {
		return errors.New("cluster_id output missing after deploy")
	}
	project, region, name, err := parseClusterID(clusterID)
	if err != nil {
		return err
	}

	cluster, err := svc.Dataproc.Projects.Regions.Clusters.Get(project, region, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "dataproc cluster %s not found after deploy", clusterID)
	}

	if cluster.Status == nil || cluster.Status.State != "RUNNING" {
		state := "<nil>"
		if cluster.Status != nil {
			state = cluster.Status.State
		}
		return errors.Errorf("dataproc cluster %s is %s after deploy, expected RUNNING", clusterID, state)
	}
	if outputs["cluster_name"] != "" && cluster.ClusterName != outputs["cluster_name"] {
		return errors.Errorf("dataproc cluster name mismatch: cloud %q vs output %q",
			cluster.ClusterName, outputs["cluster_name"])
	}
	// Both arms report the staging bucket in use (user-supplied or
	// auto-created) — its presence proves the config was materialized.
	if outputs["staging_bucket"] == "" {
		return errors.Errorf("dataproc cluster %s reported no staging_bucket output after deploy", clusterID)
	}
	return nil
}

func (v *dataprocClusterVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	clusterID := outputs["cluster_id"]
	if clusterID == "" {
		return nil
	}
	project, region, name, err := parseClusterID(clusterID)
	if err != nil {
		return err
	}

	cluster, err := svc.Dataproc.Projects.Regions.Clusters.Get(project, region, name).Context(ctx).Do()
	if err == nil {
		// A cluster mid-teardown reports DELETING — that is the destroyed
		// path, not an orphan.
		if cluster.Status != nil && cluster.Status.State == "DELETING" {
			return nil
		}
		return errors.Errorf("dataproc cluster %s still exists after destroy", clusterID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing dataproc cluster %s after destroy", clusterID)
}

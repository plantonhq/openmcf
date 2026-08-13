package verify

import (
	"context"

	"github.com/digitalocean/godo"
	pkgerrors "github.com/pkg/errors"
)

// kubernetesClusterVerifier verifies a DigitalOceanKubernetesCluster via
// GET /v2/kubernetes/clusters/{id}.
type kubernetesClusterVerifier struct{}

func (*kubernetesClusterVerifier) IDOutputKey() string { return "cluster_id" }

func (*kubernetesClusterVerifier) VerifyExists(ctx context.Context, client *godo.Client, id string) error {
	exists, err := kubernetesClusterExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceankubernetescluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("digitaloceankubernetescluster %q not found after deploy", id)
	}
	return nil
}

func (*kubernetesClusterVerifier) VerifyAbsent(ctx context.Context, client *godo.Client, id string) error {
	exists, err := kubernetesClusterExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceankubernetescluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("digitaloceankubernetescluster %q still exists after destroy", id)
	}
	return nil
}

func kubernetesClusterExists(ctx context.Context, client *godo.Client, id string) (bool, error) {
	_, _, err := client.Kubernetes.Get(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// kubernetesNodePoolVerifier verifies a DigitalOceanKubernetesNodePool. The
// API addresses node pools as /v2/kubernetes/clusters/{cluster_id}/node_pools/{id},
// but the kind's stack outputs carry only node_pool_id (no cluster id) -- a
// recorded gap for the kind's depth wave -- so the verifier resolves the
// owning cluster by scanning the account's clusters, exactly as the upstream
// provider's own node-pool importer does.
type kubernetesNodePoolVerifier struct{}

func (*kubernetesNodePoolVerifier) IDOutputKey() string { return "node_pool_id" }

func (*kubernetesNodePoolVerifier) VerifyExists(ctx context.Context, client *godo.Client, id string) error {
	exists, err := kubernetesNodePoolExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceankubernetesnodepool verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("digitaloceankubernetesnodepool %q not found after deploy", id)
	}
	return nil
}

func (*kubernetesNodePoolVerifier) VerifyAbsent(ctx context.Context, client *godo.Client, id string) error {
	exists, err := kubernetesNodePoolExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceankubernetesnodepool verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("digitaloceankubernetesnodepool %q still exists after destroy", id)
	}
	return nil
}

func kubernetesNodePoolExists(ctx context.Context, client *godo.Client, poolID string) (bool, error) {
	opts := &godo.ListOptions{PerPage: 200}
	for {
		clusters, resp, err := client.Kubernetes.List(ctx, opts)
		if err != nil {
			return false, pkgerrors.Wrap(err, "listing kubernetes clusters to locate the node pool")
		}
		for _, cluster := range clusters {
			for _, pool := range cluster.NodePools {
				if pool.ID == poolID {
					return true, nil
				}
			}
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			return false, nil
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return false, pkgerrors.Wrap(err, "reading kubernetes cluster list pagination")
		}
		opts.Page = page + 1
	}
}

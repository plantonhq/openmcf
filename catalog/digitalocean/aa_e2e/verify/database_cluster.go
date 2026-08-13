package verify

import (
	"context"

	"github.com/digitalocean/godo"
	pkgerrors "github.com/pkg/errors"
)

// databaseClusterVerifier verifies a DigitalOceanDatabaseCluster via
// GET /v2/databases/{id}.
type databaseClusterVerifier struct{}

func (*databaseClusterVerifier) IDOutputKey() string { return "cluster_id" }

func (*databaseClusterVerifier) VerifyExists(ctx context.Context, client *godo.Client, id string) error {
	exists, err := databaseClusterExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceandatabasecluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("digitaloceandatabasecluster %q not found after deploy", id)
	}
	return nil
}

func (*databaseClusterVerifier) VerifyAbsent(ctx context.Context, client *godo.Client, id string) error {
	exists, err := databaseClusterExists(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "digitaloceandatabasecluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("digitaloceandatabasecluster %q still exists after destroy", id)
	}
	return nil
}

func databaseClusterExists(ctx context.Context, client *godo.Client, id string) (bool, error) {
	_, _, err := client.Databases.Get(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

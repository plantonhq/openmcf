package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mongoClusterAPIVersion is the Microsoft.DocumentDB/mongoClusters API
// line the pinned provider vendors (mongocluster/2025-09-01).
const mongoClusterAPIVersion = "2025-09-01"

// mongoClusterVerifier verifies an AzureMongoCluster via the generic
// ARM resources GetByID (see armResourceExists), keyed on the
// cluster's ARM ID.
type mongoClusterVerifier struct{}

func (*mongoClusterVerifier) IDOutputKey() string {
	return "mongo_cluster_id"
}

func (*mongoClusterVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mongoClusterAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremongocluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremongocluster %q not found after deploy", id)
	}
	return nil
}

func (*mongoClusterVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mongoClusterAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremongocluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremongocluster %q still exists after destroy", id)
	}
	return nil
}

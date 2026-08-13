package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mongoClusterUserVerifier verifies an AzureMongoClusterUser via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// grant's cluster-scoped ARM ID ({cluster_id}/users/{object_id}).
// Rides the mongoClusters API pin.
type mongoClusterUserVerifier struct{}

func (*mongoClusterUserVerifier) IDOutputKey() string {
	return "mongo_cluster_user_id"
}

func (*mongoClusterUserVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mongoClusterAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremongoclusteruser verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremongoclusteruser %q not found after deploy", id)
	}
	return nil
}

func (*mongoClusterUserVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mongoClusterAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremongoclusteruser verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremongoclusteruser %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbMongoCollectionVerifier verifies an AzureCosmosdbMongoCollection
// via the generic ARM resources GetByID, keyed on the collection's full
// ARM ID (.../mongodbDatabases/{db}/collections/{name} is a first-class
// ARM read path).
type cosmosdbMongoCollectionVerifier struct{}

// IDOutputKey is the collection's full ARM ID.
func (*cosmosdbMongoCollectionVerifier) IDOutputKey() string {
	return "mongo_collection_id"
}

func (*cosmosdbMongoCollectionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbmongocollection verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbmongocollection %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbMongoCollectionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbmongocollection verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbmongocollection %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbMongoDatabaseVerifier verifies an AzureCosmosdbMongoDatabase
// via the generic ARM resources GetByID, keyed on the database's full
// ARM ID (.../mongodbDatabases/{name} is a first-class ARM read path).
type cosmosdbMongoDatabaseVerifier struct{}

// IDOutputKey is the database's full ARM ID.
func (*cosmosdbMongoDatabaseVerifier) IDOutputKey() string {
	return "mongo_database_id"
}

func (*cosmosdbMongoDatabaseVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbmongodatabase verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbmongodatabase %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbMongoDatabaseVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbmongodatabase verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbmongodatabase %q still exists after destroy", id)
	}
	return nil
}

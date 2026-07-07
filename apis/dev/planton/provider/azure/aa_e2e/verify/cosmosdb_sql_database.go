package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbSqlDatabaseVerifier verifies an AzureCosmosdbSqlDatabase via
// the generic ARM resources GetByID, keyed on the database's full ARM ID
// (Cosmos SQL databases have a first-class ARM read proxy under
// databaseAccounts/{account}/sqlDatabases -- no data-plane client
// needed).
type cosmosdbSqlDatabaseVerifier struct{}

// IDOutputKey is the database's full ARM ID.
func (*cosmosdbSqlDatabaseVerifier) IDOutputKey() string {
	return "sql_database_id"
}

func (*cosmosdbSqlDatabaseVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqldatabase verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbsqldatabase %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbSqlDatabaseVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqldatabase verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbsqldatabase %q still exists after destroy", id)
	}
	return nil
}

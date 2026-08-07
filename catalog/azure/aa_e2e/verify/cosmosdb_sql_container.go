package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbSqlContainerVerifier verifies an AzureCosmosdbSqlContainer via
// the generic ARM resources GetByID, keyed on the container's full ARM
// ID (.../sqlDatabases/{db}/containers/{name} is a first-class ARM read
// path).
type cosmosdbSqlContainerVerifier struct{}

// IDOutputKey is the container's full ARM ID.
func (*cosmosdbSqlContainerVerifier) IDOutputKey() string {
	return "sql_container_id"
}

func (*cosmosdbSqlContainerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqlcontainer verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbsqlcontainer %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbSqlContainerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqlcontainer verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbsqlcontainer %q still exists after destroy", id)
	}
	return nil
}

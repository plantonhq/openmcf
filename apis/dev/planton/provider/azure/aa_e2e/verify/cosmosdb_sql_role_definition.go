package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbSqlRoleDefinitionVerifier verifies an
// AzureCosmosdbSqlRoleDefinition via the generic ARM resources GetByID,
// keyed on the definition's full ARM ID (Cosmos SQL role definitions
// are management-plane objects under
// databaseAccounts/{account}/sqlRoleDefinitions -- no data-plane client
// needed).
type cosmosdbSqlRoleDefinitionVerifier struct{}

// IDOutputKey is the role definition's full ARM ID.
func (*cosmosdbSqlRoleDefinitionVerifier) IDOutputKey() string {
	return "role_definition_id"
}

func (*cosmosdbSqlRoleDefinitionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqlroledefinition verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbsqlroledefinition %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbSqlRoleDefinitionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqlroledefinition verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbsqlroledefinition %q still exists after destroy", id)
	}
	return nil
}

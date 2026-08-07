package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbSqlRoleAssignmentVerifier verifies an
// AzureCosmosdbSqlRoleAssignment via the generic ARM resources GetByID,
// keyed on the assignment's full ARM ID (Cosmos SQL role assignments
// are management-plane objects under
// databaseAccounts/{account}/sqlRoleAssignments -- no data-plane client
// needed).
type cosmosdbSqlRoleAssignmentVerifier struct{}

// IDOutputKey is the role assignment's full ARM ID.
func (*cosmosdbSqlRoleAssignmentVerifier) IDOutputKey() string {
	return "role_assignment_id"
}

func (*cosmosdbSqlRoleAssignmentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqlroleassignment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbsqlroleassignment %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbSqlRoleAssignmentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbsqlroleassignment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbsqlroleassignment %q still exists after destroy", id)
	}
	return nil
}

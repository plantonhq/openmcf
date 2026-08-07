package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cosmosdbAPIVersion pins the Microsoft.DocumentDB management API version
// used by the generic ARM existence probes for the account and its
// database/container children.
const cosmosdbAPIVersion = "2024-08-15"

// cosmosdbAccountVerifier verifies an AzureCosmosdbAccount via the
// generic ARM resources GetByID, keyed on the account's full ARM ID.
type cosmosdbAccountVerifier struct{}

// IDOutputKey is the account's full ARM ID.
func (*cosmosdbAccountVerifier) IDOutputKey() string {
	return "cosmosdb_account_id"
}

func (*cosmosdbAccountVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbaccount verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecosmosdbaccount %q not found after deploy", id)
	}
	return nil
}

func (*cosmosdbAccountVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cosmosdbAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecosmosdbaccount verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecosmosdbaccount %q still exists after destroy", id)
	}
	return nil
}

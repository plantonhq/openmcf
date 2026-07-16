package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// keyVaultAPIVersion is the stable Microsoft.KeyVault API version the generic
// existence probe is pinned to.
const keyVaultAPIVersion = "2023-07-01"

// keyVaultVerifier verifies an AzureKeyVault via the generic ARM resources
// GetByID, keyed on the vault's full ARM ID. Soft delete makes this probe
// meaningful for absence too: a soft-deleted (or purged) vault 404s the
// live-resource GET even while its name is still reserved in the
// deleted-vaults list.
type keyVaultVerifier struct{}

// IDOutputKey is the vault's full ARM ID.
func (*keyVaultVerifier) IDOutputKey() string {
	return "key_vault_id"
}

func (*keyVaultVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, keyVaultAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvault verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurekeyvault %q not found after deploy", id)
	}
	return nil
}

func (*keyVaultVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, keyVaultAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvault verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurekeyvault %q still exists after destroy", id)
	}
	return nil
}

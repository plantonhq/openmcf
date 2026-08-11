package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// keyVaultSecretAPIVersion is the stable Microsoft.KeyVault API version for
// the vaults/secrets ARM proxy.
const keyVaultSecretAPIVersion = "2023-07-01"

// keyVaultSecretVerifier verifies an AzureKeyVaultSecret through its
// read-only ARM proxy (Microsoft.KeyVault/vaults/secrets) via the generic
// GetByID probe -- secrets are data-plane objects, but ARM exposes a
// control-plane read for them (the key sibling's exact pattern), which keeps
// this verifier on the established zero-extra-permission path and never
// touches the secret's VALUE. Keyed on the secret's VERSIONLESS ARM resource
// id so the probe is version-agnostic.
type keyVaultSecretVerifier struct{}

// IDOutputKey is the secret's versionless ARM resource ID.
func (*keyVaultSecretVerifier) IDOutputKey() string {
	return "resource_versionless_id"
}

func (*keyVaultSecretVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, keyVaultSecretAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvaultsecret verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurekeyvaultsecret %q not found after deploy", id)
	}
	return nil
}

func (*keyVaultSecretVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, keyVaultSecretAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvaultsecret verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurekeyvaultsecret %q still exists after destroy", id)
	}
	return nil
}

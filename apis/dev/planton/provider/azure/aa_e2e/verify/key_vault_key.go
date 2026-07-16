package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// keyVaultKeyAPIVersion is the stable Microsoft.KeyVault API version for the
// vaults/keys ARM proxy.
const keyVaultKeyAPIVersion = "2023-07-01"

// keyVaultKeyVerifier verifies an AzureKeyVaultKey through its read-only ARM
// proxy (Microsoft.KeyVault/vaults/keys) via the generic GetByID probe --
// keys are data-plane objects, but ARM exposes a control-plane read for them,
// which keeps this verifier on the established zero-extra-permission pattern
// (certificates have no such proxy and verify through the data plane
// instead). Keyed on the key's VERSIONLESS ARM resource id so the probe is
// rotation-agnostic.
type keyVaultKeyVerifier struct{}

// IDOutputKey is the key's versionless ARM resource ID.
func (*keyVaultKeyVerifier) IDOutputKey() string {
	return "resource_versionless_id"
}

func (*keyVaultKeyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, keyVaultKeyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvaultkey verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurekeyvaultkey %q not found after deploy", id)
	}
	return nil
}

func (*keyVaultKeyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, keyVaultKeyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvaultkey verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurekeyvaultkey %q still exists after destroy", id)
	}
	return nil
}

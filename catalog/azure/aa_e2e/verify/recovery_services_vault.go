package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// recoveryServicesAPIVersion pins the Microsoft.RecoveryServices GA
// line the vault verifier reads at -- the same line the pinned azurerm
// provider vendors (recoveryservices/2025-08-01), so one source of
// truth checks both engines.
const recoveryServicesAPIVersion = "2025-08-01"

// recoveryServicesVaultVerifier verifies an AzureRecoveryServicesVault
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the vault's full ARM ID. Existence is the honest bar for a smoke
// lane: the vault is a container object -- its protections are what
// dependent kinds' lanes prove. Absence-after-destroy carries no
// soft-delete caveat for the vault OBJECT itself (protected items
// soft-delete; an empty vault deletes outright).
type recoveryServicesVaultVerifier struct{}

// IDOutputKey is the vault's full ARM ID.
func (*recoveryServicesVaultVerifier) IDOutputKey() string {
	return "recovery_services_vault_id"
}

func (*recoveryServicesVaultVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerecoveryservicesvault verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurerecoveryservicesvault %q not found after deploy", id)
	}
	return nil
}

func (*recoveryServicesVaultVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerecoveryservicesvault verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurerecoveryservicesvault %q still exists after destroy", id)
	}
	return nil
}

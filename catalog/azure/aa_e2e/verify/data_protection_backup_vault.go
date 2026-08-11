package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataProtectionAPIVersion pins the Microsoft.DataProtection GA line
// the family's verifiers read at -- the same line the pinned azurerm
// provider vendors (dataprotection/2025-07-01), so one source of truth
// checks both engines. The backup vault, backup policy, and resource
// guard verifiers all ride this pin (one service, one API line).
const dataProtectionAPIVersion = "2025-07-01"

// dataProtectionBackupVaultVerifier verifies an
// AzureDataProtectionBackupVault via the generic ARM resources GetByID
// (see armResourceExists), keyed on the vault's full ARM ID. Existence
// is the honest bar for a smoke lane: the vault is a container object
// -- its policies and instances are what dependent kinds' lanes prove.
// Absence-after-destroy notes: Azure's delete returns before the vault
// is fully gone and the provider polls it out (issue-documented), so
// by the time destroy completes, absence is genuinely absence.
type dataProtectionBackupVaultVerifier struct{}

// IDOutputKey is the vault's full ARM ID.
func (*dataProtectionBackupVaultVerifier) IDOutputKey() string {
	return "backup_vault_id"
}

func (*dataProtectionBackupVaultVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionbackupvault verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredataprotectionbackupvault %q not found after deploy", id)
	}
	return nil
}

func (*dataProtectionBackupVaultVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionbackupvault verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredataprotectionbackupvault %q still exists after destroy", id)
	}
	return nil
}

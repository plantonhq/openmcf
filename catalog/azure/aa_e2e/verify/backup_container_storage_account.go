package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// recoveryServicesProtectionContainerAPIVersion pins the
// Microsoft.RecoveryServices backup GA line the container verifier
// reads at -- the same line the pinned azurerm provider vendors for
// protection containers (recoveryservicesbackup/2023-02-01).
const recoveryServicesProtectionContainerAPIVersion = "2023-02-01"

// backupContainerStorageAccountVerifier verifies an
// AzureBackupContainerStorageAccount via the generic ARM resources
// GetByID (see armResourceExists), keyed on the registration's full
// ARM ID (.../backupFabrics/Azure/protectionContainers/
// StorageContainer;storage;{sa-rg};{sa-name}). Existence is the honest
// bar: registration is a free binding -- the protected-share lane is
// where a registered container is exercised. Unregister completes
// synchronously through the provider's LRO polling, and nothing
// soft-holds a container registration itself (soft delete rides the
// protected ITEMS -- their lane's caveat, not this one's).
type backupContainerStorageAccountVerifier struct{}

// IDOutputKey is the registration's full ARM ID.
func (*backupContainerStorageAccountVerifier) IDOutputKey() string {
	return "backup_container_id"
}

func (*backupContainerStorageAccountVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectionContainerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupcontainerstorageaccount verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebackupcontainerstorageaccount %q not found after deploy", id)
	}
	return nil
}

func (*backupContainerStorageAccountVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectionContainerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupcontainerstorageaccount verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurebackupcontainerstorageaccount %q still exists after destroy", id)
	}
	return nil
}

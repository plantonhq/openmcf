package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// recoveryServicesProtectedItemAPIVersion pins the
// Microsoft.RecoveryServices backup GA line the protected-item
// verifier reads at -- the same line the pinned azurerm provider
// vendors for protected items (recoveryservicesbackup/2023-02-01).
const recoveryServicesProtectedItemAPIVersion = "2023-02-01"

// backupProtectedVmVerifier verifies an AzureBackupProtectedVm via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// protected item's full ARM ID (.../backupFabrics/Azure/
// protectionContainers/.../protectedItems/VM;iaasvmcontainerv2;...).
// Existence is the honest bar: creation only REGISTERS protection (the
// first backup runs on the policy's schedule), so the registration
// object is what a smoke lane can verify.
//
// ABSENCE-AFTER-DESTROY CAVEAT: deleting protection SOFT-DELETES the
// item for 14 days -- ARM's read on a soft-deleted item still answers.
// The module's destroy (DeleteThenPoll) completes the delete before
// returning, and the provider treats a scheduled-for-deferred-delete
// item as gone; if this lane's absence check ever flakes on the ghost,
// the proof session records the tolerance here and in the queue
// watch-list (az backup CLI shows soft-deleted items).
type backupProtectedVmVerifier struct{}

// IDOutputKey is the protected item's full ARM ID.
func (*backupProtectedVmVerifier) IDOutputKey() string {
	return "backup_protected_vm_id"
}

func (*backupProtectedVmVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectedItemAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedvm verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebackupprotectedvm %q not found after deploy", id)
	}
	return nil
}

func (*backupProtectedVmVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectedItemAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedvm verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurebackupprotectedvm %q still exists after destroy", id)
	}
	return nil
}

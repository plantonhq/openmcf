package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// backupProtectedFileShareVerifier verifies an
// AzureBackupProtectedFileShare via the generic ARM resources GetByID
// (see armResourceExists), keyed on the protected item's full ARM ID
// (.../protectionContainers/StorageContainer;storage;{sa-rg};
// {sa-name}/protectedItems/AzureFileShare;{system-name}) and read at
// the same protected-items GA line as the VM sibling
// (recoveryServicesProtectedItemAPIVersion -- the line the pinned
// azurerm provider vendors). Existence is the honest bar: creation
// only REGISTERS protection (the first backup runs on the policy's
// schedule), so the registration object is what a smoke lane can
// verify.
//
// ABSENCE-AFTER-DESTROY CAVEAT: deleting protection SOFT-DELETES the
// item for 14 days -- ARM's read on a soft-deleted item still answers.
// The module's destroy (DeleteThenPoll) completes the delete before
// returning, and the provider treats a scheduled-for-deferred-delete
// item as gone; if this lane's absence check ever flakes on the ghost,
// the proof session records the tolerance here and in the queue
// watch-list (az backup CLI shows soft-deleted items). The ghost can
// also delay the fixture registration's unregister at teardown -- the
// queue watch-list leads with it.
type backupProtectedFileShareVerifier struct{}

// IDOutputKey is the protected item's full ARM ID.
func (*backupProtectedFileShareVerifier) IDOutputKey() string {
	return "backup_protected_file_share_id"
}

func (*backupProtectedFileShareVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectedItemAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedfileshare verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebackupprotectedfileshare %q not found after deploy", id)
	}
	return nil
}

func (*backupProtectedFileShareVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesProtectedItemAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedfileshare verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurebackupprotectedfileshare %q still exists after destroy", id)
	}
	return nil
}

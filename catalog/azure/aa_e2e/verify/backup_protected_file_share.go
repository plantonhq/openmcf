package verify

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
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
// ABSENCE-AFTER-DESTROY (measured live, proof session 056 -- the
// same class the VM sibling measured in 055): the protected-item
// delete is asynchronous BEYOND the provider's poller -- ARM reads
// keep answering 200 after a destroy the engine already reported
// successful (a smoke item has ZERO recovery points, so a landed
// delete removes it outright; the 14-day soft-delete ghost class
// only applies to items WITH recovery points). So absence is
// verified with a bounded poll: 404 is absence, a 200 whose
// properties.isScheduledForDeferredDelete is true is ALSO absence
// (the bar the azurerm provider itself uses for ghosts), and an item
// still ACTIVE at the deadline fails honestly -- the VM sibling
// measured a run where the engine's DeleteThenPoll succeeded while
// Azure ran no delete job at all, and this poll must not absorb that
// drop. Measured this session: the item's brief 200 cleared well
// before the fixture registration's unregister, which succeeded --
// the ghost never held the teardown.
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
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackupprotectedfileshare verify-absent failed for %q", id)
	}
	deadline := time.Now().Add(recoveryServicesProtectedItemAbsencePollTimeout)
	for {
		resp, err := client.GetByID(ctx, id, recoveryServicesProtectedItemAPIVersion, nil)
		if err != nil {
			var respErr *azcore.ResponseError
			if errors.As(err, &respErr) && respErr.StatusCode == 404 {
				return nil
			}
			return pkgerrors.Wrapf(err, "azurebackupprotectedfileshare verify-absent failed for %q", id)
		}
		// The read answered: a deferred-delete ghost is absence; an
		// active item may just be the async delete's read lag -- keep
		// polling until the deadline before failing.
		if properties, ok := resp.Properties.(map[string]interface{}); ok {
			if deferred, _ := properties["isScheduledForDeferredDelete"].(bool); deferred {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return pkgerrors.Errorf("azurebackupprotectedfileshare %q still exists (active, not soft-delete-ghosted) %s after destroy -- the engine-reported delete never landed", id, recoveryServicesProtectedItemAbsencePollTimeout)
		}
		select {
		case <-ctx.Done():
			return pkgerrors.Wrapf(ctx.Err(), "azurebackupprotectedfileshare verify-absent cancelled for %q", id)
		case <-time.After(recoveryServicesProtectedItemAbsencePollInterval):
		}
	}
}

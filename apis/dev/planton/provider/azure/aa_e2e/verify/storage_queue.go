package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageQueueVerifier verifies an AzureStorageQueue via the generic ARM
// resources GetByID, keyed on the queue's full ARM ID
// (.../storageAccounts/{account}/queueServices/default/queues/{name}).
// Queues have a first-class ARM read proxy, so no data-plane client is
// needed -- and ARM reads are not subject to the account's data-plane
// firewall, keeping verification independent of network rules.
type storageQueueVerifier struct{}

// IDOutputKey is the queue's full ARM ID.
func (*storageQueueVerifier) IDOutputKey() string {
	return "queue_id"
}

func (*storageQueueVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragequeue verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestoragequeue %q not found after deploy", id)
	}
	return nil
}

func (*storageQueueVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragequeue verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestoragequeue %q still exists after destroy", id)
	}
	return nil
}

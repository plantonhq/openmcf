package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageShareVerifier verifies an AzureStorageShare via the generic ARM
// resources GetByID, keyed on the share's management ARM ID
// (.../storageAccounts/{account}/fileServices/default/shares/{name}).
// Shares have a first-class ARM read proxy, so no data-plane client is
// needed -- and ARM reads are not subject to the account's data-plane
// firewall, keeping verification independent of network rules. (The
// share's rbac_scope_id output uses a different, RBAC-only segment and
// is NOT GETtable -- verification keys on share_id.)
type storageShareVerifier struct{}

// IDOutputKey is the share's management ARM ID.
func (*storageShareVerifier) IDOutputKey() string {
	return "share_id"
}

func (*storageShareVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageshare verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestorageshare %q not found after deploy", id)
	}
	return nil
}

func (*storageShareVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageshare verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestorageshare %q still exists after destroy", id)
	}
	return nil
}

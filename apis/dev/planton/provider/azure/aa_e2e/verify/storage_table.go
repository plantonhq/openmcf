package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageTableVerifier verifies an AzureStorageTable via the generic ARM
// resources GetByID, keyed on the table's resource-manager ID
// (.../storageAccounts/{account}/tableServices/default/tables/{name}).
// The Microsoft.Storage management API exposes a first-class table read
// proxy even though the PROVIDERS drive table creation through the data
// plane -- so verification stays management-plane, independent of the
// account's data-plane firewall and shared-key posture.
type storageTableVerifier struct{}

// IDOutputKey is the table's resource-manager ID.
func (*storageTableVerifier) IDOutputKey() string {
	return "table_id"
}

func (*storageTableVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragetable verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestoragetable %q not found after deploy", id)
	}
	return nil
}

func (*storageTableVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragetable verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestoragetable %q still exists after destroy", id)
	}
	return nil
}

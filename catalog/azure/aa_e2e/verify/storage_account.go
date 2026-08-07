package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageAPIVersion is the stable Microsoft.Storage API version the generic
// existence probes pin (shared by the account and container verifiers).
const storageAPIVersion = "2023-05-01"

// storageAccountVerifier verifies an AzureStorageAccount via the generic ARM
// resources GetByID, keyed on the account's full ARM ID.
type storageAccountVerifier struct{}

// IDOutputKey is the account's full ARM ID.
func (*storageAccountVerifier) IDOutputKey() string {
	return "storage_account_id"
}

func (*storageAccountVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageaccount verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestorageaccount %q not found after deploy", id)
	}
	return nil
}

func (*storageAccountVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageaccount verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestorageaccount %q still exists after destroy", id)
	}
	return nil
}

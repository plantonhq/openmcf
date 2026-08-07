package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageDataLakeGen2FilesystemVerifier verifies an
// AzureStorageDataLakeGen2Filesystem via the generic ARM resources
// GetByID, keyed on the filesystem's container-proxy ARM ID
// (.../storageAccounts/{account}/blobServices/default/containers/{name}).
// ADLS filesystems surface in ARM as blob containers -- the dual dfs/blob
// endpoints front the same namespace -- so the container proxy is a
// first-class management-plane read: no data-plane client is needed, and
// ARM reads are not subject to the account's data-plane firewall.
type storageDataLakeGen2FilesystemVerifier struct{}

// IDOutputKey is the filesystem's container-proxy ARM ID.
func (*storageDataLakeGen2FilesystemVerifier) IDOutputKey() string {
	return "filesystem_id"
}

func (*storageDataLakeGen2FilesystemVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragedatalakegen2filesystem verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestoragedatalakegen2filesystem %q not found after deploy", id)
	}
	return nil
}

func (*storageDataLakeGen2FilesystemVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragedatalakegen2filesystem verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestoragedatalakegen2filesystem %q still exists after destroy", id)
	}
	return nil
}

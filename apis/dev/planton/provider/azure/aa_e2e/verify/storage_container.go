package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageContainerVerifier verifies an AzureStorageContainer via the generic
// ARM resources GetByID, keyed on the container's full ARM ID
// (.../storageAccounts/{account}/blobServices/default/containers/{name}).
// Containers have a first-class ARM read proxy, so no data-plane client is
// needed -- and ARM reads are not subject to the account's data-plane
// firewall, keeping verification independent of network rules.
type storageContainerVerifier struct{}

// IDOutputKey is the container's full ARM ID.
func (*storageContainerVerifier) IDOutputKey() string {
	return "container_id"
}

func (*storageContainerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragecontainer verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestoragecontainer %q not found after deploy", id)
	}
	return nil
}

func (*storageContainerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragecontainer verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestoragecontainer %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// managedDiskAPIVersion is the stable Microsoft.Compute disks API version
// the generic existence probe is pinned to.
const managedDiskAPIVersion = "2024-03-02"

// managedDiskVerifier verifies an AzureManagedDisk via the generic ARM
// resources GetByID (see armResourceExists), keyed on the disk's full ARM
// ID.
type managedDiskVerifier struct{}

// IDOutputKey is the disk's full ARM ID.
func (*managedDiskVerifier) IDOutputKey() string {
	return "disk_id"
}

func (*managedDiskVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, managedDiskAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanageddisk verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremanageddisk %q not found after deploy", id)
	}
	return nil
}

func (*managedDiskVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, managedDiskAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanageddisk verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremanageddisk %q still exists after destroy", id)
	}
	return nil
}

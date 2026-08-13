package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// diskSnapshotAPIVersion is the stable Microsoft.Compute snapshots API
// version the generic existence probe is pinned to (the provider's own
// snapshots SDK pin at v5.0.0).
const diskSnapshotAPIVersion = "2022-03-02"

// diskSnapshotVerifier verifies an AzureDiskSnapshot via the generic ARM
// resources GetByID (see armResourceExists), keyed on the snapshot's
// full ARM ID.
type diskSnapshotVerifier struct{}

// IDOutputKey is the snapshot's full ARM ID.
func (*diskSnapshotVerifier) IDOutputKey() string {
	return "snapshot_id"
}

func (*diskSnapshotVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, diskSnapshotAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredisksnapshot verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredisksnapshot %q not found after deploy", id)
	}
	return nil
}

func (*diskSnapshotVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, diskSnapshotAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredisksnapshot verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredisksnapshot %q still exists after destroy", id)
	}
	return nil
}

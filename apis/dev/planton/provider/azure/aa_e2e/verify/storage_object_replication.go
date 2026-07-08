package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageObjectReplicationVerifier verifies an
// AzureStorageObjectReplication via the generic ARM resources GetByID,
// keyed on the DESTINATION-side policy ARM ID
// (.../storageAccounts/{dst}/objectReplicationPolicies/{policyId}) --
// the authoritative copy Azure materializes first and assigns rule IDs
// on. Destroy removes the policy from BOTH accounts; probing the
// destination side is sufficient because the source mirror cannot exist
// without it.
type storageObjectReplicationVerifier struct{}

// IDOutputKey is the destination-side policy ARM ID.
func (*storageObjectReplicationVerifier) IDOutputKey() string {
	return "destination_object_replication_id"
}

func (*storageObjectReplicationVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageobjectreplication verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestorageobjectreplication %q not found after deploy", id)
	}
	return nil
}

func (*storageObjectReplicationVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageobjectreplication verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestorageobjectreplication %q still exists after destroy", id)
	}
	return nil
}

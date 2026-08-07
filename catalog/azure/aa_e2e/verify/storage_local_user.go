package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// storageLocalUserVerifier verifies an AzureStorageLocalUser via the
// generic ARM resources GetByID, keyed on the user's full ARM ID
// (.../storageAccounts/{account}/localUsers/{name}). Local users are
// pure management-plane objects (their credentials serve the SFTP data
// plane, but the identity itself is ARM-resident), so a plain
// existence probe is the right shape.
type storageLocalUserVerifier struct{}

// IDOutputKey is the local user's full ARM ID.
func (*storageLocalUserVerifier) IDOutputKey() string {
	return "local_user_id"
}

func (*storageLocalUserVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragelocaluser verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestoragelocaluser %q not found after deploy", id)
	}
	return nil
}

func (*storageLocalUserVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, storageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestoragelocaluser verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurestoragelocaluser %q still exists after destroy", id)
	}
	return nil
}

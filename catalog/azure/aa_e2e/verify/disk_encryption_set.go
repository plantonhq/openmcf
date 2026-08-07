package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// diskEncryptionSetAPIVersion is the stable Microsoft.Compute API version
// the generic existence probe is pinned to.
const diskEncryptionSetAPIVersion = "2023-04-02"

// diskEncryptionSetVerifier verifies an AzureDiskEncryptionSet via the
// generic ARM resources GetByID (see armResourceExists), keyed on the set's
// full ARM ID.
type diskEncryptionSetVerifier struct{}

// IDOutputKey is the set's full ARM ID.
func (*diskEncryptionSetVerifier) IDOutputKey() string {
	return "disk_encryption_set_id"
}

func (*diskEncryptionSetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, diskEncryptionSetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurediskencryptionset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurediskencryptionset %q not found after deploy", id)
	}
	return nil
}

func (*diskEncryptionSetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, diskEncryptionSetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurediskencryptionset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurediskencryptionset %q still exists after destroy", id)
	}
	return nil
}

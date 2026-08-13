package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataProtectionBackupInstanceVerifier verifies an
// AzureDataProtectionBackupInstance via the generic ARM resources
// GetByID (see armResourceExists), keyed on the instance's full ARM ID
// (.../backupVaults/{vault}/backupInstances/{name}). The SAME ID shape
// serves all six variant resources, so one verifier covers the union.
// Absence after destroy is immediate when the vault's soft delete is
// Off (the fixture vault's required posture); on a soft-delete-On
// vault the deleted instance lingers as a soft-deleted item that this
// ARM read no longer returns -- the ORPHAN SWEEP, not this verifier,
// owns that class (`az dataprotection backup-instance
// list-deleted`-style checks).
type dataProtectionBackupInstanceVerifier struct{}

// IDOutputKey is the instance's full ARM ID.
func (*dataProtectionBackupInstanceVerifier) IDOutputKey() string {
	return "backup_instance_id"
}

func (*dataProtectionBackupInstanceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionbackupinstance verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredataprotectionbackupinstance %q not found after deploy", id)
	}
	return nil
}

func (*dataProtectionBackupInstanceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionbackupinstance verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredataprotectionbackupinstance %q still exists after destroy", id)
	}
	return nil
}

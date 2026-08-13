package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataProtectionBackupPolicyVerifier verifies an
// AzureDataProtectionBackupPolicy via the generic ARM resources
// GetByID (see armResourceExists), keyed on the policy's full ARM ID
// (.../backupVaults/{vault}/backupPolicies/{name}). The SAME ID shape
// serves all six variant resources, so one verifier covers the union.
// Policies are pure configuration objects with no soft-delete
// semantics -- absence after destroy is immediate.
type dataProtectionBackupPolicyVerifier struct{}

// IDOutputKey is the policy's full ARM ID.
func (*dataProtectionBackupPolicyVerifier) IDOutputKey() string {
	return "backup_policy_id"
}

func (*dataProtectionBackupPolicyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionbackuppolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredataprotectionbackuppolicy %q not found after deploy", id)
	}
	return nil
}

func (*dataProtectionBackupPolicyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionbackuppolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredataprotectionbackuppolicy %q still exists after destroy", id)
	}
	return nil
}

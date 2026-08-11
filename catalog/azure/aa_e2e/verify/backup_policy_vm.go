package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// recoveryServicesBackupPolicyAPIVersion pins the
// Microsoft.RecoveryServices backup GA line the policy verifier reads
// at -- the same line the pinned azurerm provider vendors for
// protection policies (recoveryservicesbackup/2024-10-01).
const recoveryServicesBackupPolicyAPIVersion = "2024-10-01"

// backupPolicyVmVerifier verifies an AzureBackupPolicyVm via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// policy's full ARM ID (.../vaults/{vault}/backupPolicies/{name}).
// Existence is the honest bar: the policy is a free configuration
// object -- its schedule fires on the service side, and the protected
// VM's lane is where a policy is exercised against a machine. Policies
// delete outright (nothing soft-holds a policy name).
type backupPolicyVmVerifier struct{}

// IDOutputKey is the policy's full ARM ID.
func (*backupPolicyVmVerifier) IDOutputKey() string {
	return "backup_policy_id"
}

func (*backupPolicyVmVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesBackupPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackuppolicyvm verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebackuppolicyvm %q not found after deploy", id)
	}
	return nil
}

func (*backupPolicyVmVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesBackupPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackuppolicyvm verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurebackuppolicyvm %q still exists after destroy", id)
	}
	return nil
}

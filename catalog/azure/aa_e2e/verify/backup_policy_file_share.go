package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// backupPolicyFileShareVerifier verifies an AzureBackupPolicyFileShare
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the policy's full ARM ID (.../vaults/{vault}/backupPolicies/
// {name}) and read at the same protection-policies GA line as the VM
// policy verifier (recoveryServicesBackupPolicyAPIVersion -- the line
// the pinned azurerm provider vendors). Existence is the honest bar:
// the policy is a free configuration object -- its schedule fires on
// the service side, and the protected-share lane is where a policy is
// exercised against a share. Policies delete outright (nothing
// soft-holds a policy name).
type backupPolicyFileShareVerifier struct{}

// IDOutputKey is the policy's full ARM ID.
func (*backupPolicyFileShareVerifier) IDOutputKey() string {
	return "backup_policy_id"
}

func (*backupPolicyFileShareVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesBackupPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackuppolicyfileshare verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebackuppolicyfileshare %q not found after deploy", id)
	}
	return nil
}

func (*backupPolicyFileShareVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, recoveryServicesBackupPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebackuppolicyfileshare verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurebackuppolicyfileshare %q still exists after destroy", id)
	}
	return nil
}

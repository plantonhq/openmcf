package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// managedRedisAccessPolicyAssignmentVerifier verifies an
// AzureManagedRedisAccessPolicyAssignment via the generic ARM resources
// GetByID, keyed on the assignment's full ARM ID
// (.../redisEnterprise/{cluster}/databases/default/accessPolicyAssignments/{objectId}).
type managedRedisAccessPolicyAssignmentVerifier struct{}

// IDOutputKey is the assignment's full ARM ID.
func (*managedRedisAccessPolicyAssignmentVerifier) IDOutputKey() string {
	return "access_policy_assignment_id"
}

func (*managedRedisAccessPolicyAssignmentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, managedRedisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanagedredisaccesspolicyassignment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremanagedredisaccesspolicyassignment %q not found after deploy", id)
	}
	return nil
}

func (*managedRedisAccessPolicyAssignmentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, managedRedisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanagedredisaccesspolicyassignment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremanagedredisaccesspolicyassignment %q still exists after destroy", id)
	}
	return nil
}

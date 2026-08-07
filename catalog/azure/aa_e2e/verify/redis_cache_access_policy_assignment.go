package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// redisCacheAccessPolicyAssignmentVerifier verifies an
// AzureRedisCacheAccessPolicyAssignment via the generic ARM resources
// GetByID, keyed on the assignment's full ARM ID
// (.../redis/{cache}/accessPolicyAssignments/{name}).
type redisCacheAccessPolicyAssignmentVerifier struct{}

// IDOutputKey is the assignment's full ARM ID.
func (*redisCacheAccessPolicyAssignmentVerifier) IDOutputKey() string {
	return "access_policy_assignment_id"
}

func (*redisCacheAccessPolicyAssignmentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerediscacheaccesspolicyassignment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurerediscacheaccesspolicyassignment %q not found after deploy", id)
	}
	return nil
}

func (*redisCacheAccessPolicyAssignmentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerediscacheaccesspolicyassignment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurerediscacheaccesspolicyassignment %q still exists after destroy", id)
	}
	return nil
}

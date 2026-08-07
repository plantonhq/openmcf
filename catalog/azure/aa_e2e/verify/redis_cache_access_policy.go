package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// redisCacheAccessPolicyVerifier verifies an AzureRedisCacheAccessPolicy
// via the generic ARM resources GetByID, keyed on the policy's full ARM
// ID (.../redis/{cache}/accessPolicies/{name}).
type redisCacheAccessPolicyVerifier struct{}

// IDOutputKey is the access policy's full ARM ID.
func (*redisCacheAccessPolicyVerifier) IDOutputKey() string {
	return "access_policy_id"
}

func (*redisCacheAccessPolicyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerediscacheaccesspolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurerediscacheaccesspolicy %q not found after deploy", id)
	}
	return nil
}

func (*redisCacheAccessPolicyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerediscacheaccesspolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurerediscacheaccesspolicy %q still exists after destroy", id)
	}
	return nil
}

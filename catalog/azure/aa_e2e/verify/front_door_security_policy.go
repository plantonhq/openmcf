package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorSecurityPolicyVerifier verifies an
// AzureFrontDoorSecurityPolicy via the generic ARM resources GetByID,
// keyed on the security policy's full ARM ID (a Microsoft.Cdn profile
// child, so the family API version applies).
type frontDoorSecurityPolicyVerifier struct{}

// IDOutputKey is the security policy's full ARM ID.
func (*frontDoorSecurityPolicyVerifier) IDOutputKey() string {
	return "security_policy_id"
}

func (*frontDoorSecurityPolicyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorsecuritypolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorsecuritypolicy %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorSecurityPolicyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorsecuritypolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorsecuritypolicy %q still exists after destroy", id)
	}
	return nil
}

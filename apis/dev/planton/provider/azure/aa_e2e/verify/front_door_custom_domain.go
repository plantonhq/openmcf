package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorCustomDomainVerifier verifies an AzureFrontDoorCustomDomain
// via the generic ARM resources GetByID, keyed on the domain's full ARM
// ID. A domain is GETtable while still in the pending-validation state
// (creation never blocks on DNS proof), so existence is the right
// signal regardless of validation progress.
type frontDoorCustomDomainVerifier struct{}

// IDOutputKey is the custom domain's full ARM ID.
func (*frontDoorCustomDomainVerifier) IDOutputKey() string {
	return "custom_domain_id"
}

func (*frontDoorCustomDomainVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorcustomdomain verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorcustomdomain %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorCustomDomainVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorcustomdomain verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorcustomdomain %q still exists after destroy", id)
	}
	return nil
}

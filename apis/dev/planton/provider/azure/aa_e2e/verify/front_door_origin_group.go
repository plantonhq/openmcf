package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorOriginGroupVerifier verifies an AzureFrontDoorOriginGroup via
// the generic ARM resources GetByID, keyed on the group's full ARM ID.
type frontDoorOriginGroupVerifier struct{}

// IDOutputKey is the origin group's full ARM ID.
func (*frontDoorOriginGroupVerifier) IDOutputKey() string {
	return "origin_group_id"
}

func (*frontDoorOriginGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoororigingroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoororigingroup %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorOriginGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoororigingroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoororigingroup %q still exists after destroy", id)
	}
	return nil
}

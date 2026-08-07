package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorRouteVerifier verifies an AzureFrontDoorRoute via the generic
// ARM resources GetByID, keyed on the route's full ARM ID.
type frontDoorRouteVerifier struct{}

// IDOutputKey is the route's full ARM ID.
func (*frontDoorRouteVerifier) IDOutputKey() string {
	return "route_id"
}

func (*frontDoorRouteVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorroute verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorroute %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorRouteVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorroute verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorroute %q still exists after destroy", id)
	}
	return nil
}

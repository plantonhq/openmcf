package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorOriginVerifier verifies an AzureFrontDoorOrigin via the
// generic ARM resources GetByID, keyed on the origin's full ARM ID.
type frontDoorOriginVerifier struct{}

// IDOutputKey is the origin's full ARM ID.
func (*frontDoorOriginVerifier) IDOutputKey() string {
	return "origin_id"
}

func (*frontDoorOriginVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoororigin verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoororigin %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorOriginVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoororigin verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoororigin %q still exists after destroy", id)
	}
	return nil
}

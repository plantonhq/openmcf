package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cdnFrontDoorAPIVersion pins the Microsoft.Cdn resource-provider API
// version for the whole Front Door family (profile, endpoint, origin
// group, origin, route) -- a GA version line that serves every child
// type under profiles/.
const cdnFrontDoorAPIVersion = "2024-02-01"

// frontDoorProfileVerifier verifies an AzureFrontDoorProfile via the
// generic ARM resources GetByID, keyed on the profile's full ARM ID.
type frontDoorProfileVerifier struct{}

// IDOutputKey is the profile's full ARM ID.
func (*frontDoorProfileVerifier) IDOutputKey() string {
	return "profile_id"
}

func (*frontDoorProfileVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorprofile verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorprofile %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorProfileVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorprofile verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorprofile %q still exists after destroy", id)
	}
	return nil
}

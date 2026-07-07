package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorEndpointVerifier verifies an AzureFrontDoorEndpoint via the
// generic ARM resources GetByID, keyed on the endpoint's full ARM ID.
type frontDoorEndpointVerifier struct{}

// IDOutputKey is the endpoint's full ARM ID.
func (*frontDoorEndpointVerifier) IDOutputKey() string {
	return "endpoint_id"
}

func (*frontDoorEndpointVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorendpoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorendpoint %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorEndpointVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorendpoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorendpoint %q still exists after destroy", id)
	}
	return nil
}

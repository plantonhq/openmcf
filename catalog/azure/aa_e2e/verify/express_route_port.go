package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// expressRoutePortVerifier verifies an AzureExpressRoutePort via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// port's full ARM ID. Existence is the honest bar: the physical
// cross-connects are ordered out-of-band with the facility, so the ARM
// object -- with its link pair, billing model, and composed
// authorization children -- is what provisioning can promise. The
// ExpressRoute family shares the pinned expressRouteAPIVersion.
type expressRoutePortVerifier struct{}

// IDOutputKey is the port's full ARM ID.
func (*expressRoutePortVerifier) IDOutputKey() string {
	return "express_route_port_id"
}

func (*expressRoutePortVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, expressRouteAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressrouteport verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureexpressrouteport %q not found after deploy", id)
	}
	return nil
}

func (*expressRoutePortVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, expressRouteAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressrouteport verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureexpressrouteport %q still exists after destroy", id)
	}
	return nil
}

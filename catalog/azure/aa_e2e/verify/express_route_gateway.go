package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// expressRouteGatewayVerifier verifies an AzureExpressRouteGateway via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the gateway's full ARM ID. Existence is the honest bar: connections
// require a provider-PROVISIONED circuit (none exists behind the test
// subscription), so the gateway object in its hub -- billing from
// creation at its scale-unit floor -- is what provisioning can promise.
// The gateway is a Virtual WAN resource (Microsoft.Network/
// expressRouteGateways) and shares the pinned virtualWanAPIVersion.
type expressRouteGatewayVerifier struct{}

// IDOutputKey is the gateway's full ARM ID.
func (*expressRouteGatewayVerifier) IDOutputKey() string {
	return "express_route_gateway_id"
}

func (*expressRouteGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressroutegateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureexpressroutegateway %q not found after deploy", id)
	}
	return nil
}

func (*expressRouteGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressroutegateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureexpressroutegateway %q still exists after destroy", id)
	}
	return nil
}

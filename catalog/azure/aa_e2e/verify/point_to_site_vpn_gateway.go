package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// pointToSiteVpnGatewayVerifier verifies an AzurePointToSiteVpnGateway
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the gateway's full ARM ID. Existence is the honest bar: the
// gateway's instance pair settles during the (30-45 minute) create the
// provider already polls, and no client ever connects in a lane (the
// fixture root certificate has no private key) -- ARM state, never
// connected clients. The Virtual WAN family shares the pinned
// virtualWanAPIVersion.
type pointToSiteVpnGatewayVerifier struct{}

// IDOutputKey is the gateway's full ARM ID.
func (*pointToSiteVpnGatewayVerifier) IDOutputKey() string {
	return "point_to_site_vpn_gateway_id"
}

func (*pointToSiteVpnGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepointtositevpngateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurepointtositevpngateway %q not found after deploy", id)
	}
	return nil
}

func (*pointToSiteVpnGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepointtositevpngateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurepointtositevpngateway %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// vpnGatewayVerifier verifies an AzureVpnGateway via the generic ARM
// resources GetByID (see armResourceExists), keyed on the gateway's
// full ARM ID. Existence is the honest bar: the gateway's instance
// pair and assigned public IPs settle during the (30-45 minute)
// create the provider already polls, and the composed NAT rules are
// ARM children under the same path -- the provider's own read cycle
// gates the deploy on them. The Virtual WAN family shares the pinned
// virtualWanAPIVersion.
type vpnGatewayVerifier struct{}

// IDOutputKey is the gateway's full ARM ID.
func (*vpnGatewayVerifier) IDOutputKey() string {
	return "vpn_gateway_id"
}

func (*vpnGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpngateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevpngateway %q not found after deploy", id)
	}
	return nil
}

func (*vpnGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpngateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevpngateway %q still exists after destroy", id)
	}
	return nil
}

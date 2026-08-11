package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// vpnGatewayConnectionVerifier verifies an AzureVpnGatewayConnection
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the connection's full ARM child ID
// (.../vpnGateways/{gateway}/vpnConnections/{name}). Existence is the
// honest bar: ARM provisions the connection when its parameters are
// valid, while each tunnel reaches Connected only when the branch
// device negotiates -- with no real device behind test lanes, the
// verifier asserts ARM state, never tunnel-up
// (provisioned-is-not-connected). The Virtual WAN family shares the
// pinned virtualWanAPIVersion.
type vpnGatewayConnectionVerifier struct{}

// IDOutputKey is the connection's full ARM child ID.
func (*vpnGatewayConnectionVerifier) IDOutputKey() string {
	return "connection_id"
}

func (*vpnGatewayConnectionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpngatewayconnection verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevpngatewayconnection %q not found after deploy", id)
	}
	return nil
}

func (*vpnGatewayConnectionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpngatewayconnection verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevpngatewayconnection %q still exists after destroy", id)
	}
	return nil
}

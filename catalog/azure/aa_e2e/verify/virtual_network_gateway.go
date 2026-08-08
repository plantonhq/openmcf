package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualNetworkGatewayAPIVersion is the stable Microsoft.Network API
// version the generic existence probes for the gateway family are pinned
// to (gateways, connections, and local network gateways share the
// namespace and GA API line).
const virtualNetworkGatewayAPIVersion = "2024-05-01"

// virtualNetworkGatewayVerifier verifies an AzureVirtualNetworkGateway via
// the generic ARM resources GetByID (see armResourceExists), keyed on the
// gateway's full ARM ID. Existence is the honest bar: a gateway can take
// 25-45 minutes to provision, and the runner's deploy phase has already
// waited for ARM Succeeded by the time this probe runs.
type virtualNetworkGatewayVerifier struct{}

// IDOutputKey is the virtual network gateway's full ARM ID.
func (*virtualNetworkGatewayVerifier) IDOutputKey() string {
	return "virtual_network_gateway_id"
}

func (*virtualNetworkGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetworkgateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualnetworkgateway %q not found after deploy", id)
	}
	return nil
}

func (*virtualNetworkGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetworkgateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualnetworkgateway %q still exists after destroy", id)
	}
	return nil
}

// virtualNetworkGatewayConnectionVerifier verifies an
// AzureVirtualNetworkGatewayConnection by ARM resource existence.
// Deliberately NOT the tunnel's connectionStatus: the e2e fixture site has
// no real device behind it, so the tunnel provisions but never reaches
// Connected -- ARM existence is what provisioning can honestly promise.
type virtualNetworkGatewayConnectionVerifier struct{}

// IDOutputKey is the connection's full ARM ID.
func (*virtualNetworkGatewayConnectionVerifier) IDOutputKey() string {
	return "connection_id"
}

func (*virtualNetworkGatewayConnectionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetworkgatewayconnection verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualnetworkgatewayconnection %q not found after deploy", id)
	}
	return nil
}

func (*virtualNetworkGatewayConnectionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetworkgatewayconnection verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualnetworkgatewayconnection %q still exists after destroy", id)
	}
	return nil
}

// localNetworkGatewayVerifier verifies an AzureLocalNetworkGateway by ARM
// resource existence -- the object is pure ARM metadata (nothing deploys
// on-premises), so existence IS the complete deployment contract.
type localNetworkGatewayVerifier struct{}

// IDOutputKey is the local network gateway's full ARM ID.
func (*localNetworkGatewayVerifier) IDOutputKey() string {
	return "local_network_gateway_id"
}

func (*localNetworkGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurelocalnetworkgateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurelocalnetworkgateway %q not found after deploy", id)
	}
	return nil
}

func (*localNetworkGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurelocalnetworkgateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurelocalnetworkgateway %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// expressRouteAPIVersion is the stable Microsoft.Network API version the
// existence probes for the ExpressRoute family are pinned to (circuits
// and their peering children share the namespace and GA API line).
const expressRouteAPIVersion = "2024-05-01"

// expressRouteCircuitVerifier verifies an AzureExpressRouteCircuit via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the circuit's full ARM ID. Existence is the honest bar: the provider
// side legitimately stays NotProvisioned in the lane (no carrier behind
// the test subscription), so the ARM object -- with its SKU, provider
// trio, and composed authorization children -- is what provisioning can
// promise.
type expressRouteCircuitVerifier struct{}

// IDOutputKey is the circuit's full ARM ID.
func (*expressRouteCircuitVerifier) IDOutputKey() string {
	return "express_route_circuit_id"
}

func (*expressRouteCircuitVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, expressRouteAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressroutecircuit verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureexpressroutecircuit %q not found after deploy", id)
	}
	return nil
}

func (*expressRouteCircuitVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, expressRouteAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressroutecircuit verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureexpressroutecircuit %q still exists after destroy", id)
	}
	return nil
}

// expressRouteCircuitPeeringVerifier verifies an
// AzureExpressRouteCircuitPeering by ARM resource existence, keyed on
// the peering's full ARM ID (the circuit path plus /peerings/{type}).
// Deliberately NOT the BGP session state: the fixture circuit has no
// carrier behind it, so a session can never establish -- the ARM child's
// existence is what configuration can honestly promise.
type expressRouteCircuitPeeringVerifier struct{}

// IDOutputKey is the peering's full ARM ID.
func (*expressRouteCircuitPeeringVerifier) IDOutputKey() string {
	return "express_route_circuit_peering_id"
}

func (*expressRouteCircuitPeeringVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, expressRouteAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressroutecircuitpeering verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureexpressroutecircuitpeering %q not found after deploy", id)
	}
	return nil
}

func (*expressRouteCircuitPeeringVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, expressRouteAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureexpressroutecircuitpeering verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureexpressroutecircuitpeering %q still exists after destroy", id)
	}
	return nil
}

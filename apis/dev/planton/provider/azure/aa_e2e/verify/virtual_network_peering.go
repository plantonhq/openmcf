package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualNetworkPeeringAPIVersion is the stable Microsoft.Network API version
// the generic existence probe is pinned to.
const virtualNetworkPeeringAPIVersion = "2024-05-01"

// virtualNetworkPeeringVerifier verifies an AzureVirtualNetworkPeering via
// the generic ARM resources GetByID (see armResourceExists), keyed on the
// peering's full ARM ID (a child of its local virtual network).
type virtualNetworkPeeringVerifier struct{}

// IDOutputKey is the peering's full ARM ID.
func (*virtualNetworkPeeringVerifier) IDOutputKey() string {
	return "peering_id"
}

func (*virtualNetworkPeeringVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkPeeringAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetworkpeering verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualnetworkpeering %q not found after deploy", id)
	}
	return nil
}

func (*virtualNetworkPeeringVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkPeeringAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetworkpeering verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualnetworkpeering %q still exists after destroy", id)
	}
	return nil
}

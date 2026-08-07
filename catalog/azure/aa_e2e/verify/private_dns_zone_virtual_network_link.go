package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// privateDnsZoneVirtualNetworkLinkVerifier verifies an
// AzurePrivateDnsZoneVirtualNetworkLink via the generic ARM resources
// GetByID (see armResourceExists), keyed on the link's full ARM ID
// ({zone-id}/virtualNetworkLinks/{name}). Links are children of the zone and
// served by the same resource provider (and API version) as the zones
// themselves.
type privateDnsZoneVirtualNetworkLinkVerifier struct{}

// IDOutputKey is the link's full ARM ID.
func (*privateDnsZoneVirtualNetworkLinkVerifier) IDOutputKey() string {
	return "link_id"
}

func (*privateDnsZoneVirtualNetworkLinkVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednszonevirtualnetworklink verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivatednszonevirtualnetworklink %q not found after deploy", id)
	}
	return nil
}

func (*privateDnsZoneVirtualNetworkLinkVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednszonevirtualnetworklink verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivatednszonevirtualnetworklink %q still exists after destroy", id)
	}
	return nil
}

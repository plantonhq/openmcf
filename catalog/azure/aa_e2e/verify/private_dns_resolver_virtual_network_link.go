package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// privateDnsResolverVirtualNetworkLinkVerifier verifies an
// AzurePrivateDnsResolverVirtualNetworkLink via the generic ARM
// resources GetByID (see armResourceExists), keyed on the link's full
// ARM ID ({ruleset_id}/virtualNetworkLinks/{name}). Pinned to the
// dnsresolver family's API version (dnsResolverAPIVersion).
type privateDnsResolverVirtualNetworkLinkVerifier struct{}

// IDOutputKey is the link's full ARM ID.
func (*privateDnsResolverVirtualNetworkLinkVerifier) IDOutputKey() string {
	return "virtual_network_link_id"
}

func (*privateDnsResolverVirtualNetworkLinkVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dnsResolverAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsresolvervirtualnetworklink verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivatednsresolvervirtualnetworklink %q not found after deploy", id)
	}
	return nil
}

func (*privateDnsResolverVirtualNetworkLinkVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dnsResolverAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsresolvervirtualnetworklink verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivatednsresolvervirtualnetworklink %q still exists after destroy", id)
	}
	return nil
}

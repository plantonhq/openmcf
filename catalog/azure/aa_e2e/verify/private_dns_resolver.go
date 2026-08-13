package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dnsResolverAPIVersion is the stable Microsoft.Network API version the
// existence probes are pinned to -- the line the pinned provider vendors
// for dnsResolvers, dnsForwardingRulesets, and their children.
const dnsResolverAPIVersion = "2022-07-01"

// privateDnsResolverVerifier verifies an AzurePrivateDnsResolver via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// resolver's full ARM ID. The composed endpoints live and die with the
// resolver (ARM children) -- the resolver's absence proves theirs.
type privateDnsResolverVerifier struct{}

// IDOutputKey is the resolver's full ARM ID.
func (*privateDnsResolverVerifier) IDOutputKey() string {
	return "dns_resolver_id"
}

func (*privateDnsResolverVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dnsResolverAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsresolver verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivatednsresolver %q not found after deploy", id)
	}
	return nil
}

func (*privateDnsResolverVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dnsResolverAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsresolver verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivatednsresolver %q still exists after destroy", id)
	}
	return nil
}

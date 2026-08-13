package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// privateDnsResolverForwardingRulesetVerifier verifies an
// AzurePrivateDnsResolverForwardingRuleset via the generic ARM resources
// GetByID (see armResourceExists), keyed on the ruleset's full ARM ID.
// The composed forwarding rules live and die with the ruleset (ARM
// children) -- the ruleset's absence proves theirs. Pinned to the
// dnsresolver family's API version (dnsResolverAPIVersion).
type privateDnsResolverForwardingRulesetVerifier struct{}

// IDOutputKey is the ruleset's full ARM ID.
func (*privateDnsResolverForwardingRulesetVerifier) IDOutputKey() string {
	return "dns_forwarding_ruleset_id"
}

func (*privateDnsResolverForwardingRulesetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dnsResolverAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsresolverforwardingruleset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivatednsresolverforwardingruleset %q not found after deploy", id)
	}
	return nil
}

func (*privateDnsResolverForwardingRulesetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dnsResolverAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsresolverforwardingruleset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivatednsresolverforwardingruleset %q still exists after destroy", id)
	}
	return nil
}

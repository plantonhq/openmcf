package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// firewallPolicyAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const firewallPolicyAPIVersion = "2024-05-01"

// firewallPolicyVerifier verifies an AzureFirewallPolicy via the generic
// ARM resources GetByID (see armResourceExists), keyed on the policy's
// full ARM ID.
type firewallPolicyVerifier struct{}

// IDOutputKey is the policy's full ARM ID.
func (*firewallPolicyVerifier) IDOutputKey() string {
	return "firewall_policy_id"
}

func (*firewallPolicyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, firewallPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefirewallpolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefirewallpolicy %q not found after deploy", id)
	}
	return nil
}

func (*firewallPolicyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, firewallPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefirewallpolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefirewallpolicy %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorWafAPIVersion pins the API version for the Front Door WAF
// policy. The policy is a Microsoft.Network type
// (frontDoorWebApplicationFirewallPolicies) -- NOT a Microsoft.Cdn
// profile child -- so the CDN family version does not apply; this is
// the GA version line the azurerm v4.80 provider builds against.
const frontDoorWafAPIVersion = "2025-03-01"

// frontDoorFirewallPolicyVerifier verifies an
// AzureFrontDoorFirewallPolicy via the generic ARM resources GetByID,
// keyed on the policy's full ARM ID.
type frontDoorFirewallPolicyVerifier struct{}

// IDOutputKey is the policy's full ARM ID.
func (*frontDoorFirewallPolicyVerifier) IDOutputKey() string {
	return "firewall_policy_id"
}

func (*frontDoorFirewallPolicyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, frontDoorWafAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorfirewallpolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorfirewallpolicy %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorFirewallPolicyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, frontDoorWafAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorfirewallpolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorfirewallpolicy %q still exists after destroy", id)
	}
	return nil
}

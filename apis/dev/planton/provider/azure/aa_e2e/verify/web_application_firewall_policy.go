package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// wafPolicyAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const wafPolicyAPIVersion = "2024-05-01"

// webApplicationFirewallPolicyVerifier verifies an
// AzureWebApplicationFirewallPolicy via the generic ARM resources GetByID,
// keyed on the policy's full ARM ID.
type webApplicationFirewallPolicyVerifier struct{}

// IDOutputKey is the policy's full ARM ID.
func (*webApplicationFirewallPolicyVerifier) IDOutputKey() string {
	return "policy_id"
}

func (*webApplicationFirewallPolicyVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, wafPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurewebapplicationfirewallpolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurewebapplicationfirewallpolicy %q not found after deploy", id)
	}
	return nil
}

func (*webApplicationFirewallPolicyVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, wafPolicyAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurewebapplicationfirewallpolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurewebapplicationfirewallpolicy %q still exists after destroy", id)
	}
	return nil
}

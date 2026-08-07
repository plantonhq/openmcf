package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// firewallPolicyRuleCollectionGroupAPIVersion is the stable
// Microsoft.Network API version the generic existence probe is pinned to.
const firewallPolicyRuleCollectionGroupAPIVersion = "2024-05-01"

// firewallPolicyRuleCollectionGroupVerifier verifies an
// AzureFirewallPolicyRuleCollectionGroup via the generic ARM resources
// GetByID (see armResourceExists) -- the group is a real nested ARM
// resource, so the child-resource ID probes directly.
type firewallPolicyRuleCollectionGroupVerifier struct{}

// IDOutputKey is the group's full (nested) ARM ID.
func (*firewallPolicyRuleCollectionGroupVerifier) IDOutputKey() string {
	return "rule_collection_group_id"
}

func (*firewallPolicyRuleCollectionGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, firewallPolicyRuleCollectionGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefirewallpolicyrulecollectiongroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefirewallpolicyrulecollectiongroup %q not found after deploy", id)
	}
	return nil
}

func (*firewallPolicyRuleCollectionGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, firewallPolicyRuleCollectionGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefirewallpolicyrulecollectiongroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefirewallpolicyrulecollectiongroup %q still exists after destroy", id)
	}
	return nil
}

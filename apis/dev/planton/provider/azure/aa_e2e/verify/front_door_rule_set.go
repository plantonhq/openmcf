package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorRuleSetVerifier verifies an AzureFrontDoorRuleSet via the
// generic ARM resources GetByID, keyed on the rule set's full ARM ID.
// The folded rules are children of the set in ARM, so the set's
// existence/absence covers the whole delivery policy's lifecycle.
type frontDoorRuleSetVerifier struct{}

// IDOutputKey is the rule set's full ARM ID.
func (*frontDoorRuleSetVerifier) IDOutputKey() string {
	return "rule_set_id"
}

func (*frontDoorRuleSetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorruleset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorruleset %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorRuleSetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorruleset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorruleset %q still exists after destroy", id)
	}
	return nil
}

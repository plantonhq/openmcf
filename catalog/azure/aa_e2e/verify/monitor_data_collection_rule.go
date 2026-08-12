package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorDataCollectionRuleAPIVersion pins the Microsoft.Insights
// dataCollectionRules RP version the verifier reads with -- the line the
// provider vendors (insights/2023-03-11/datacollectionrules). The same
// version serves the association verifier (one vendored line covers the
// whole DCR family).
const monitorDataCollectionRuleAPIVersion = "2023-03-11"

// monitorDataCollectionRuleVerifier verifies an
// AzureMonitorDataCollectionRule via the generic ARM resources GetByID
// (see armResourceExists), keyed on the rule's ARM ID.
type monitorDataCollectionRuleVerifier struct{}

func (*monitorDataCollectionRuleVerifier) IDOutputKey() string {
	return "data_collection_rule_id"
}

func (*monitorDataCollectionRuleVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorDataCollectionRuleAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitordatacollectionrule verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitordatacollectionrule %q not found after deploy", id)
	}
	return nil
}

func (*monitorDataCollectionRuleVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorDataCollectionRuleAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitordatacollectionrule verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitordatacollectionrule %q still exists after destroy", id)
	}
	return nil
}

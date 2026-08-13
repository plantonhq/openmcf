package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorDataCollectionRuleAssociationVerifier verifies an
// AzureMonitorDataCollectionRuleAssociation via the generic ARM
// resources GetByID (see armResourceExists), keyed on the association's
// TARGET-scoped ARM ID
// ({target_id}/providers/Microsoft.Insights/dataCollectionRuleAssociations/{name})
// -- GetByID resolves extension-resource ids like any other. Pinned to
// the same insights 2023-03-11 line as the rule verifier (the version
// the provider vendors for the whole DCR family).
type monitorDataCollectionRuleAssociationVerifier struct{}

func (*monitorDataCollectionRuleAssociationVerifier) IDOutputKey() string {
	return "data_collection_rule_association_id"
}

func (*monitorDataCollectionRuleAssociationVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorDataCollectionRuleAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitordatacollectionruleassociation verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitordatacollectionruleassociation %q not found after deploy", id)
	}
	return nil
}

func (*monitorDataCollectionRuleAssociationVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorDataCollectionRuleAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitordatacollectionruleassociation verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitordatacollectionruleassociation %q still exists after destroy", id)
	}
	return nil
}

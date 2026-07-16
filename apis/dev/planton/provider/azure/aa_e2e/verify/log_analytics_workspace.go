package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// operationalInsightsAPIVersion pins the Microsoft.OperationalInsights RP
// version the workspace verifier reads with.
const operationalInsightsAPIVersion = "2023-09-01"

// logAnalyticsWorkspaceVerifier verifies an AzureLogAnalyticsWorkspace via
// the generic ARM resources GetByID (see armResourceExists), keyed on the
// workspace's ARM ID.
type logAnalyticsWorkspaceVerifier struct{}

// IDOutputKey is the workspace's ARM resource ID (NOT the customer GUID,
// which the catalog exports separately as workspace_customer_id).
func (*logAnalyticsWorkspaceVerifier) IDOutputKey() string {
	return "workspace_id"
}

func (*logAnalyticsWorkspaceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, operationalInsightsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureloganalyticsworkspace verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureloganalyticsworkspace %q not found after deploy", id)
	}
	return nil
}

func (*logAnalyticsWorkspaceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, operationalInsightsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureloganalyticsworkspace verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureloganalyticsworkspace %q still exists after destroy", id)
	}
	return nil
}

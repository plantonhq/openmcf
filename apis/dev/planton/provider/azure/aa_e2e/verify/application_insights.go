package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// applicationInsightsAPIVersion pins the Microsoft.Insights components RP
// version the verifier reads with.
const applicationInsightsAPIVersion = "2020-02-02"

// applicationInsightsVerifier verifies an AzureApplicationInsights via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// component's ARM ID.
type applicationInsightsVerifier struct{}

func (*applicationInsightsVerifier) IDOutputKey() string {
	return "application_insights_id"
}

func (*applicationInsightsVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, applicationInsightsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationinsights verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureapplicationinsights %q not found after deploy", id)
	}
	return nil
}

func (*applicationInsightsVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, applicationInsightsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationinsights verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureapplicationinsights %q still exists after destroy", id)
	}
	return nil
}

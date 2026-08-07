package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// standardWebTestAPIVersion is the stable Microsoft.Insights API version the
// generic existence probe is pinned to.
const standardWebTestAPIVersion = "2022-06-15"

// applicationInsightsStandardWebTestVerifier verifies an
// AzureApplicationInsightsStandardWebTest via the generic ARM resources
// GetByID (see armResourceExists), keyed on the web test's full ARM ID.
type applicationInsightsStandardWebTestVerifier struct{}

// IDOutputKey is the web test's full ARM ID.
func (*applicationInsightsStandardWebTestVerifier) IDOutputKey() string {
	return "web_test_id"
}

func (*applicationInsightsStandardWebTestVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, standardWebTestAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationinsightsstandardwebtest verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureapplicationinsightsstandardwebtest %q not found after deploy", id)
	}
	return nil
}

func (*applicationInsightsStandardWebTestVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, standardWebTestAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationinsightsstandardwebtest verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureapplicationinsightsstandardwebtest %q still exists after destroy", id)
	}
	return nil
}

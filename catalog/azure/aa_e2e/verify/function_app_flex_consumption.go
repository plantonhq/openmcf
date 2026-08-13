package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// functionAppFlexConsumptionAPIVersion is the stable Microsoft.Web sites
// API version the generic existence probe is pinned to -- the same
// version the provider's own webapps client uses. Flex Consumption apps
// are Microsoft.Web sites (kind functionapp,linux) on an FC1 plan.
const functionAppFlexConsumptionAPIVersion = "2023-12-01"

// functionAppFlexConsumptionVerifier verifies an
// AzureFunctionAppFlexConsumption via the generic ARM resources GetByID
// (see armResourceExists), keyed on the site's full ARM ID.
type functionAppFlexConsumptionVerifier struct{}

// IDOutputKey is the flex consumption app's full ARM ID.
func (*functionAppFlexConsumptionVerifier) IDOutputKey() string {
	return "function_app_id"
}

func (*functionAppFlexConsumptionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, functionAppFlexConsumptionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefunctionappflexconsumption verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefunctionappflexconsumption %q not found after deploy", id)
	}
	return nil
}

func (*functionAppFlexConsumptionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, functionAppFlexConsumptionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefunctionappflexconsumption verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefunctionappflexconsumption %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// servicePlanAPIVersion is the stable GA Microsoft.Web API version the
// generic existence probe is pinned to; the probe only needs GetByID,
// so any GA generation the service supports works.
const servicePlanAPIVersion = "2023-12-01"

// servicePlanVerifier verifies an AzureServicePlan via the generic ARM
// resources GetByID (see armResourceExists), keyed on the plan's full
// ARM ID.
type servicePlanVerifier struct{}

// IDOutputKey is the plan's full ARM ID.
func (*servicePlanVerifier) IDOutputKey() string {
	return "service_plan_id"
}

func (*servicePlanVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, servicePlanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureserviceplan verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureserviceplan %q not found after deploy", id)
	}
	return nil
}

func (*servicePlanVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, servicePlanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureserviceplan verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureserviceplan %q still exists after destroy", id)
	}
	return nil
}

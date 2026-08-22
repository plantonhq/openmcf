package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// plantonRunnerVerifier verifies an AzurePlantonRunner via the generic ARM
// resources GetByID (see armResourceExists), keyed on the runner's
// Container App ARM ID. The appliance is a single-replica Container App;
// Container Apps reports the app provisioned independently of replica
// health, so existence is the honest provisioning-level assertion (the
// runner's own join is proven at the control plane, not here).
type plantonRunnerVerifier struct{}

// IDOutputKey is the app's full ARM ID.
func (*plantonRunnerVerifier) IDOutputKey() string {
	return "container_app_id"
}

func (*plantonRunnerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureplantonrunner verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureplantonrunner %q not found after deploy", id)
	}
	return nil
}

func (*plantonRunnerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureplantonrunner verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureplantonrunner %q still exists after destroy", id)
	}
	return nil
}

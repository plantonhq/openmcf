package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cognitiveDeploymentVerifier verifies an AzureCognitiveDeployment via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the deployment's full ARM ID (an account child:
// .../accounts/{account}/deployments/{name}). Existence is the honest
// bar: a returned deployment IS a provisioned model -- ARM rejects the
// create outright when the region or the subscription's quota cannot
// host the model, so there is no provisioned-but-degraded state to
// probe. The cognitive family shares the pinned cognitiveAPIVersion.
type cognitiveDeploymentVerifier struct{}

// IDOutputKey is the deployment's full ARM ID.
func (*cognitiveDeploymentVerifier) IDOutputKey() string {
	return "deployment_id"
}

func (*cognitiveDeploymentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cognitiveAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecognitivedeployment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecognitivedeployment %q not found after deploy", id)
	}
	return nil
}

func (*cognitiveDeploymentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cognitiveAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecognitivedeployment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecognitivedeployment %q still exists after destroy", id)
	}
	return nil
}

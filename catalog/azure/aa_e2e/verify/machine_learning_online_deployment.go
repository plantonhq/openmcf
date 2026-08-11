package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningOnlineDeploymentVerifier verifies an
// AzureMachineLearningOnlineDeployment via the generic ARM resources
// GetByID (see armResourceExists), keyed on the deployment's full ARM
// ID (.../onlineEndpoints/{endpoint}/deployments/{name}). The read
// rides the same pinned api-version the Terraform module writes -- one
// source of truth checks both engines regardless of which api-version
// each engine's client spoke. Absence-after-destroy needs no
// soft-delete caveat -- deployments delete outright with their
// instances.
type machineLearningOnlineDeploymentVerifier struct{}

// IDOutputKey is the deployment's full ARM ID.
func (*machineLearningOnlineDeploymentVerifier) IDOutputKey() string {
	return "online_deployment_id"
}

func (*machineLearningOnlineDeploymentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningonlinedeployment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningonlinedeployment %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningOnlineDeploymentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningonlinedeployment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningonlinedeployment %q still exists after destroy", id)
	}
	return nil
}

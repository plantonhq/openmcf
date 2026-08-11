package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningBatchDeploymentVerifier verifies an
// AzureMachineLearningBatchDeployment via the generic ARM resources
// GetByID (see armResourceExists), keyed on the deployment's full ARM
// ID (.../batchEndpoints/{endpoint}/deployments/{name}). Existence is
// the honest bar: a batch deployment is a job RECIPE -- it provisions
// nothing at create time (compute materializes per job), so the
// deployment OBJECT is what a smoke lane can verify. The read rides
// the same pinned api-version the Terraform module writes -- one
// source of truth checks both engines regardless of which api-version
// each engine's client spoke.
type machineLearningBatchDeploymentVerifier struct{}

// IDOutputKey is the deployment's full ARM ID.
func (*machineLearningBatchDeploymentVerifier) IDOutputKey() string {
	return "batch_deployment_id"
}

func (*machineLearningBatchDeploymentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningbatchdeployment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningbatchdeployment %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningBatchDeploymentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningbatchdeployment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningbatchdeployment %q still exists after destroy", id)
	}
	return nil
}

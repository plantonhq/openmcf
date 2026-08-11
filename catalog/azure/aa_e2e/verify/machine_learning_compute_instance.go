package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningComputeInstanceVerifier verifies an
// AzureMachineLearningComputeInstance via the generic ARM resources
// GetByID (see armResourceExists), keyed on the compute's full ARM ID
// (.../workspaces/{ws}/computes/{name} -- instances and clusters share
// this collection). Absence-after-destroy is what the region-wide NAME
// reservation makes worth verifying: the ARM object must be gone even
// though the name can stay reserved briefly (the scenario's run-id
// token is what handles the reservation).
type machineLearningComputeInstanceVerifier struct{}

// IDOutputKey is the instance's full ARM ID.
func (*machineLearningComputeInstanceVerifier) IDOutputKey() string {
	return "machine_learning_compute_instance_id"
}

func (*machineLearningComputeInstanceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningcomputeinstance verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningcomputeinstance %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningComputeInstanceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningcomputeinstance verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningcomputeinstance %q still exists after destroy", id)
	}
	return nil
}

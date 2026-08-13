package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningBatchEndpointVerifier verifies an
// AzureMachineLearningBatchEndpoint via the generic ARM resources
// GetByID (see armResourceExists), keyed on the endpoint's full ARM ID
// (.../workspaces/{ws}/batchEndpoints/{name}). Existence is the honest
// bar: a batch endpoint is a routing object -- it holds no compute and
// runs nothing until a job is submitted, so the endpoint OBJECT is
// what a smoke lane can verify. The read rides the same pinned
// api-version the Terraform module writes -- one source of truth
// checks both engines regardless of which api-version each engine's
// client spoke. Absence-after-destroy needs no soft-delete caveat --
// endpoints delete outright (their WORKSPACE soft-deletes).
type machineLearningBatchEndpointVerifier struct{}

// IDOutputKey is the endpoint's full ARM ID.
func (*machineLearningBatchEndpointVerifier) IDOutputKey() string {
	return "batch_endpoint_id"
}

func (*machineLearningBatchEndpointVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningbatchendpoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningbatchendpoint %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningBatchEndpointVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningbatchendpoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningbatchendpoint %q still exists after destroy", id)
	}
	return nil
}

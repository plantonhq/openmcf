package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningOnlineEndpointVerifier verifies an
// AzureMachineLearningOnlineEndpoint via the generic ARM resources
// GetByID (see armResourceExists), keyed on the endpoint's full ARM ID
// (.../workspaces/{ws}/onlineEndpoints/{name}). Existence is the honest
// bar: an endpoint without deployments holds no instances, so the
// endpoint OBJECT is what a smoke lane can verify. The read rides the
// same pinned api-version the Terraform module writes -- one source of
// truth checks both engines regardless of which api-version each
// engine's client spoke. Absence-after-destroy needs no soft-delete
// caveat -- endpoints delete outright (their WORKSPACE soft-deletes),
// though the name reservation releases only after deletion completes.
type machineLearningOnlineEndpointVerifier struct{}

// IDOutputKey is the endpoint's full ARM ID.
func (*machineLearningOnlineEndpointVerifier) IDOutputKey() string {
	return "online_endpoint_id"
}

func (*machineLearningOnlineEndpointVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningonlineendpoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningonlineendpoint %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningOnlineEndpointVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningonlineendpoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningonlineendpoint %q still exists after destroy", id)
	}
	return nil
}

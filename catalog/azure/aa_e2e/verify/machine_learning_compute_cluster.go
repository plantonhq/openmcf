package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningComputeClusterVerifier verifies an
// AzureMachineLearningComputeCluster via the generic ARM resources
// GetByID (see armResourceExists), keyed on the compute's full ARM ID
// (.../workspaces/{ws}/computes/{name} -- clusters and instances share
// this collection). Existence is the honest bar: a scale-to-zero
// cluster holds no nodes, so the compute OBJECT is what a smoke lane
// can verify. Absence-after-destroy needs no soft-delete caveat --
// computes delete outright; it is their WORKSPACE that soft-deletes.
type machineLearningComputeClusterVerifier struct{}

// IDOutputKey is the cluster's full ARM ID.
func (*machineLearningComputeClusterVerifier) IDOutputKey() string {
	return "machine_learning_compute_cluster_id"
}

func (*machineLearningComputeClusterVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningcomputecluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningcomputecluster %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningComputeClusterVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningcomputecluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningcomputecluster %q still exists after destroy", id)
	}
	return nil
}

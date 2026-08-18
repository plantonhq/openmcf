package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningAPIVersion pins the Microsoft.MachineLearningServices
// GA line the ML-family verifiers read with -- the version that carries
// workspaces, datastores, and the managed network's outbound rules.
const machineLearningAPIVersion = "2025-06-01"

// machineLearningWorkspaceVerifier verifies an
// AzureMachineLearningWorkspace via the generic ARM resources GetByID
// (see armResourceExists), keyed on the workspace's full ARM ID.
// Existence is the honest bar: the provider's own read cycle gates the
// deploy on the workspace's properties, and the composed outbound
// rules are ARM children under the same path.
// Absence-after-destroy is genuine absence: a soft-deleted workspace
// ghost is not returned by GetByID. Ghosts have NO list API (portal
// "Recently deleted" only); the modules purge on destroy via the
// provider's machine_learning features flag, and the dual-engine
// lanes prove the purge by recreating the same fixed name.
type machineLearningWorkspaceVerifier struct{}

// IDOutputKey is the workspace's full ARM ID.
func (*machineLearningWorkspaceVerifier) IDOutputKey() string {
	return "machine_learning_workspace_id"
}

func (*machineLearningWorkspaceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningworkspace verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningworkspace %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningWorkspaceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningworkspace verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningworkspace %q still exists after destroy", id)
	}
	return nil
}

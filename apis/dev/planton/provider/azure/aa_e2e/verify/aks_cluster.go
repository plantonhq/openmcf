package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// aksClusterAPIVersion is the stable Microsoft.ContainerService API version
// the generic ARM GetByID probe below is pinned to.
const aksClusterAPIVersion = "2024-08-01"

// aksClusterVerifier verifies an AzureAksCluster via generic ARM GetByID,
// keyed on the cluster's full ARM ID (cluster_id output).
type aksClusterVerifier struct{}

func (*aksClusterVerifier) IDOutputKey() string { return "cluster_id" }

func (*aksClusterVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, aksClusterAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureakscluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureakscluster %q not found after deploy", id)
	}
	return nil
}

func (*aksClusterVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, aksClusterAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureakscluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureakscluster %q still exists after destroy", id)
	}
	return nil
}

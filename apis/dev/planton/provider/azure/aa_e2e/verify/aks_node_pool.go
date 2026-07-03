package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// aksNodePoolAPIVersion is the stable Microsoft.ContainerService API version
// the generic ARM GetByID probe below is pinned to. Agent pools are child
// resources of a managed cluster; their ARM id ends in /agentPools/{name}.
const aksNodePoolAPIVersion = "2024-08-01"

// aksNodePoolVerifier verifies an AzureAksNodePool via generic ARM GetByID,
// keyed on the agent pool's full ARM ID (node_pool_id output).
type aksNodePoolVerifier struct{}

func (*aksNodePoolVerifier) IDOutputKey() string { return "node_pool_id" }

func (*aksNodePoolVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, aksNodePoolAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureaksnodepool verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureaksnodepool %q not found after deploy", id)
	}
	return nil
}

func (*aksNodePoolVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, aksNodePoolAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureaksnodepool verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureaksnodepool %q still exists after destroy", id)
	}
	return nil
}

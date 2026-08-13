package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryIntegrationRuntimeVerifier verifies an
// AzureDataFactoryIntegrationRuntime via the generic ARM resources
// GetByID (see armResourceExists), keyed on the runtime's ARM ID --
// one ID shape serves all 3 engine flavors
// ({factory_id}/integrationRuntimes/{name}), on the same
// datafactory/2018-06-01 API line the pinned provider vendors.
type dataFactoryIntegrationRuntimeVerifier struct{}

func (*dataFactoryIntegrationRuntimeVerifier) IDOutputKey() string {
	return "integration_runtime_id"
}

func (*dataFactoryIntegrationRuntimeVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactoryintegrationruntime verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactoryintegrationruntime %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryIntegrationRuntimeVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactoryintegrationruntime verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactoryintegrationruntime %q still exists after destroy", id)
	}
	return nil
}

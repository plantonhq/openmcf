package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppEnvironmentDaprComponentVerifier verifies an
// AzureContainerAppEnvironmentDaprComponent via the generic ARM resources
// GetByID (see armResourceExists), keyed on the component's full ARM ID
// (a child of the managed environment).
type containerAppEnvironmentDaprComponentVerifier struct{}

// IDOutputKey is the Dapr component's full ARM ID.
func (*containerAppEnvironmentDaprComponentVerifier) IDOutputKey() string {
	return "dapr_component_id"
}

func (*containerAppEnvironmentDaprComponentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentdaprcomponent verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentdaprcomponent %q not found after deploy", id)
	}
	return nil
}

func (*containerAppEnvironmentDaprComponentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentdaprcomponent verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentdaprcomponent %q still exists after destroy", id)
	}
	return nil
}

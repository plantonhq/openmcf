package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppsAPIVersion is the stable Microsoft.App API version the
// generic existence probes are pinned to (the same generation azurerm v4's
// containerapps service targets).
const containerAppsAPIVersion = "2025-01-01"

// containerAppEnvironmentVerifier verifies an AzureContainerAppEnvironment
// via the generic ARM resources GetByID (see armResourceExists), keyed on
// the environment's full ARM ID.
type containerAppEnvironmentVerifier struct{}

// IDOutputKey is the environment's full ARM ID.
func (*containerAppEnvironmentVerifier) IDOutputKey() string {
	return "environment_id"
}

func (*containerAppEnvironmentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappenvironment %q not found after deploy", id)
	}
	return nil
}

func (*containerAppEnvironmentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappenvironment %q still exists after destroy", id)
	}
	return nil
}

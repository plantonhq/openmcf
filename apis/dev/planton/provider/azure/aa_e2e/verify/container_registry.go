package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerRegistryAPIVersion is the stable Microsoft.ContainerRegistry API
// version the generic existence probe is pinned to.
const containerRegistryAPIVersion = "2023-07-01"

// containerRegistryVerifier verifies an AzureContainerRegistry via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// registry's full ARM ID.
type containerRegistryVerifier struct{}

// IDOutputKey is the registry's full ARM ID.
func (*containerRegistryVerifier) IDOutputKey() string {
	return "container_registry_id"
}

func (*containerRegistryVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerRegistryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerregistry verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerregistry %q not found after deploy", id)
	}
	return nil
}

func (*containerRegistryVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerRegistryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerregistry verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerregistry %q still exists after destroy", id)
	}
	return nil
}

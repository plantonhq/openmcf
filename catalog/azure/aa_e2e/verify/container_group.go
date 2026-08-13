package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerGroupAPIVersion is the stable Microsoft.ContainerInstance
// containerGroups API version the generic existence probe is pinned to
// (the provider's own containerinstance SDK pin at v5.0.0).
const containerGroupAPIVersion = "2025-09-01"

// containerGroupVerifier verifies an AzureContainerInstance via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// container group's full ARM ID.
type containerGroupVerifier struct{}

// IDOutputKey is the container group's full ARM ID.
func (*containerGroupVerifier) IDOutputKey() string {
	return "container_group_id"
}

func (*containerGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerinstance verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerinstance %q not found after deploy", id)
	}
	return nil
}

func (*containerGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerinstance verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerinstance %q still exists after destroy", id)
	}
	return nil
}

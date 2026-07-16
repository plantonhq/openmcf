package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppEnvironmentStorageVerifier verifies an
// AzureContainerAppEnvironmentStorage via the generic ARM resources
// GetByID (see armResourceExists), keyed on the registration's full ARM ID
// (a child of the managed environment).
type containerAppEnvironmentStorageVerifier struct{}

// IDOutputKey is the storage registration's full ARM ID.
func (*containerAppEnvironmentStorageVerifier) IDOutputKey() string {
	return "storage_id"
}

func (*containerAppEnvironmentStorageVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentstorage verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentstorage %q not found after deploy", id)
	}
	return nil
}

func (*containerAppEnvironmentStorageVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentstorage verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentstorage %q still exists after destroy", id)
	}
	return nil
}

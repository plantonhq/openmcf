package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppVerifier verifies an AzureContainerApp via the generic ARM
// resources GetByID (see armResourceExists), keyed on the app's full ARM ID.
type containerAppVerifier struct{}

// IDOutputKey is the app's full ARM ID.
func (*containerAppVerifier) IDOutputKey() string {
	return "container_app_id"
}

func (*containerAppVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerapp verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerapp %q not found after deploy", id)
	}
	return nil
}

func (*containerAppVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerapp verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerapp %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridNamespaceVerifier verifies an AzureEventgridNamespace via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the namespace's ARM ID. Rides the family's eventgrid API pin (the
// 2025-02-15 line carries namespaces).
type eventgridNamespaceVerifier struct{}

func (*eventgridNamespaceVerifier) IDOutputKey() string {
	return "namespace_id"
}

func (*eventgridNamespaceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridnamespace verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgridnamespace %q not found after deploy", id)
	}
	return nil
}

func (*eventgridNamespaceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridnamespace verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgridnamespace %q still exists after destroy", id)
	}
	return nil
}

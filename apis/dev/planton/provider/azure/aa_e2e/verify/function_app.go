package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// functionAppAPIVersion is the stable Microsoft.Web sites API version the
// generic existence probe is pinned to (function apps are Microsoft.Web
// sites with kind functionapp,linux).
const functionAppAPIVersion = "2023-12-01"

// functionAppVerifier verifies an AzureFunctionApp via the generic ARM
// resources GetByID (see armResourceExists), keyed on the site's full
// ARM ID.
type functionAppVerifier struct{}

// IDOutputKey is the function app's full ARM ID.
func (*functionAppVerifier) IDOutputKey() string {
	return "function_app_id"
}

func (*functionAppVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, functionAppAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefunctionapp verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefunctionapp %q not found after deploy", id)
	}
	return nil
}

func (*functionAppVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, functionAppAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefunctionapp verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefunctionapp %q still exists after destroy", id)
	}
	return nil
}

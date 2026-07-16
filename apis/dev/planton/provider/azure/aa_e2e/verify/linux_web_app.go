package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// webAppAPIVersion is the stable Microsoft.Web sites API version the
// generic existence probe is pinned to.
const webAppAPIVersion = "2023-12-01"

// linuxWebAppVerifier verifies an AzureLinuxWebApp via the generic ARM
// resources GetByID (see armResourceExists), keyed on the site's full
// ARM ID.
type linuxWebAppVerifier struct{}

// IDOutputKey is the web app's full ARM ID.
func (*linuxWebAppVerifier) IDOutputKey() string {
	return "web_app_id"
}

func (*linuxWebAppVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, webAppAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurelinuxwebapp verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurelinuxwebapp %q not found after deploy", id)
	}
	return nil
}

func (*linuxWebAppVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, webAppAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurelinuxwebapp verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurelinuxwebapp %q still exists after destroy", id)
	}
	return nil
}

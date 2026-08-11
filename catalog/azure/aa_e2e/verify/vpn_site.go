package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// vpnSiteVerifier verifies an AzureVpnSite via the generic ARM
// resources GetByID (see armResourceExists), keyed on the site's full
// ARM ID. The site is a pure description object (its links are inline
// properties, and it deploys nothing at the branch), so ARM existence
// IS the complete deployment claim. The Virtual WAN family shares the
// pinned virtualWanAPIVersion.
type vpnSiteVerifier struct{}

// IDOutputKey is the site's full ARM ID.
func (*vpnSiteVerifier) IDOutputKey() string {
	return "vpn_site_id"
}

func (*vpnSiteVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpnsite verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevpnsite %q not found after deploy", id)
	}
	return nil
}

func (*vpnSiteVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpnsite verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevpnsite %q still exists after destroy", id)
	}
	return nil
}

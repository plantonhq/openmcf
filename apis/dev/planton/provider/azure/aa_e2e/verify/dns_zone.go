package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// publicDnsAPIVersion is the stable Microsoft.Network public DNS API
// version the generic existence probes are pinned to (zones and their
// record sets share it -- same resource provider).
const publicDnsAPIVersion = "2018-05-01"

// dnsZoneVerifier verifies an AzureDnsZone via the generic ARM resources
// GetByID (see armResourceExists), keyed on the zone's full ARM ID.
// Public DNS zones are global resources; the subscription-scoped generic
// client resolves them by ID like any other tracked resource.
type dnsZoneVerifier struct{}

// IDOutputKey is the zone's full ARM ID.
func (*dnsZoneVerifier) IDOutputKey() string {
	return "zone_id"
}

func (*dnsZoneVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurednszone verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurednszone %q not found after deploy", id)
	}
	return nil
}

func (*dnsZoneVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurednszone verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurednszone %q still exists after destroy", id)
	}
	return nil
}

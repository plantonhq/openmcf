package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// privateLinkServiceAPIVersion is the stable Microsoft.Network API
// version the existence probe is pinned to (the Private Link surface
// rides the same GA network line as the gateway family).
const privateLinkServiceAPIVersion = "2024-05-01"

// privateLinkServiceVerifier verifies an AzurePrivateLinkService via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// service's full ARM ID. Existence is the honest bar for the provider
// side alone: the alias is ARM-generated with the object, and consumer
// connections (the part that would prove traffic) require a consuming
// private endpoint that is deliberately out of this lane's scope.
type privateLinkServiceVerifier struct{}

// IDOutputKey is the Private Link Service's full ARM ID.
func (*privateLinkServiceVerifier) IDOutputKey() string {
	return "private_link_service_id"
}

func (*privateLinkServiceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateLinkServiceAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatelinkservice verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivatelinkservice %q not found after deploy", id)
	}
	return nil
}

func (*privateLinkServiceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateLinkServiceAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatelinkservice verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivatelinkservice %q still exists after destroy", id)
	}
	return nil
}

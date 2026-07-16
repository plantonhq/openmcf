package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// routeTableAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const routeTableAPIVersion = "2024-05-01"

// routeTableVerifier verifies an AzureRouteTable via the generic ARM
// resources GetByID (see armResourceExists), keyed on the table's full ARM
// ID.
type routeTableVerifier struct{}

// IDOutputKey is the route table's full ARM ID.
func (*routeTableVerifier) IDOutputKey() string {
	return "route_table_id"
}

func (*routeTableVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, routeTableAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureroutetable verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureroutetable %q not found after deploy", id)
	}
	return nil
}

func (*routeTableVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, routeTableAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureroutetable verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureroutetable %q still exists after destroy", id)
	}
	return nil
}

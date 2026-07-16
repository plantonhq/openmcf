package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mssqlFailoverGroupAPIVersion is the stable Microsoft.Sql API version the
// generic existence probe is pinned to.
const mssqlFailoverGroupAPIVersion = "2023-08-01-preview"

// mssqlFailoverGroupVerifier verifies an AzureMssqlFailoverGroup via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// group's full ARM ID (a child of the primary logical server).
type mssqlFailoverGroupVerifier struct{}

// IDOutputKey is the failover group's full ARM ID.
func (*mssqlFailoverGroupVerifier) IDOutputKey() string {
	return "failover_group_id"
}

func (*mssqlFailoverGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlFailoverGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqlfailovergroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremssqlfailovergroup %q not found after deploy", id)
	}
	return nil
}

func (*mssqlFailoverGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlFailoverGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqlfailovergroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremssqlfailovergroup %q still exists after destroy", id)
	}
	return nil
}

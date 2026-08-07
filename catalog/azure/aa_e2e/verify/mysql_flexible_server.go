package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mysqlFlexibleServerAPIVersion is the stable Microsoft.DBforMySQL API
// version the generic existence probe is pinned to.
const mysqlFlexibleServerAPIVersion = "2023-12-30"

// mysqlFlexibleServerVerifier verifies an AzureMysqlFlexibleServer via the
// generic ARM resources GetByID, keyed on the server's full ARM ID.
type mysqlFlexibleServerVerifier struct{}

// IDOutputKey is the server's full ARM ID.
func (*mysqlFlexibleServerVerifier) IDOutputKey() string {
	return "server_id"
}

func (*mysqlFlexibleServerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mysqlFlexibleServerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremysqlflexibleserver verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremysqlflexibleserver %q not found after deploy", id)
	}
	return nil
}

func (*mysqlFlexibleServerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mysqlFlexibleServerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremysqlflexibleserver verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremysqlflexibleserver %q still exists after destroy", id)
	}
	return nil
}

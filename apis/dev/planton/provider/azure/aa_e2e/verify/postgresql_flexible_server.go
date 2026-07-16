package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// postgresqlFlexibleServerAPIVersion is the stable Microsoft.DBforPostgreSQL
// API version the generic existence probe is pinned to.
const postgresqlFlexibleServerAPIVersion = "2024-08-01"

// postgresqlFlexibleServerVerifier verifies an AzurePostgresqlFlexibleServer
// via the generic ARM resources GetByID, keyed on the server's full ARM ID.
type postgresqlFlexibleServerVerifier struct{}

// IDOutputKey is the server's full ARM ID.
func (*postgresqlFlexibleServerVerifier) IDOutputKey() string {
	return "server_id"
}

func (*postgresqlFlexibleServerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, postgresqlFlexibleServerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepostgresqlflexibleserver verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurepostgresqlflexibleserver %q not found after deploy", id)
	}
	return nil
}

func (*postgresqlFlexibleServerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, postgresqlFlexibleServerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepostgresqlflexibleserver verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurepostgresqlflexibleserver %q still exists after destroy", id)
	}
	return nil
}

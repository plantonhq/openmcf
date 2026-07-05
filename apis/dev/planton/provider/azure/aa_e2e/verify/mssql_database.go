package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mssqlDatabaseVerifier verifies an AzureMssqlDatabase via the generic ARM
// resources GetByID, keyed on the database's full ARM ID.
type mssqlDatabaseVerifier struct{}

// IDOutputKey is the database's full ARM ID.
func (*mssqlDatabaseVerifier) IDOutputKey() string {
	return "database_id"
}

func (*mssqlDatabaseVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqldatabase verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremssqldatabase %q not found after deploy", id)
	}
	return nil
}

func (*mssqlDatabaseVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqldatabase verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremssqldatabase %q still exists after destroy", id)
	}
	return nil
}

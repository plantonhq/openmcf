package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mssqlAPIVersion is the Microsoft.Sql API version the generic existence
// probes are pinned to -- the same version the provider deploys with
// (Microsoft.Sql ships no stable-track version covering the current
// server/database/pool surface).
const mssqlAPIVersion = "2023-08-01-preview"

// mssqlServerVerifier verifies an AzureMssqlServer via the generic ARM
// resources GetByID, keyed on the server's full ARM ID.
type mssqlServerVerifier struct{}

// IDOutputKey is the server's full ARM ID.
func (*mssqlServerVerifier) IDOutputKey() string {
	return "server_id"
}

func (*mssqlServerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqlserver verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremssqlserver %q not found after deploy", id)
	}
	return nil
}

func (*mssqlServerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqlserver verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremssqlserver %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// mssqlElasticPoolVerifier verifies an AzureMssqlElasticPool via the
// generic ARM resources GetByID, keyed on the pool's full ARM ID.
type mssqlElasticPoolVerifier struct{}

// IDOutputKey is the pool's full ARM ID.
func (*mssqlElasticPoolVerifier) IDOutputKey() string {
	return "elastic_pool_id"
}

func (*mssqlElasticPoolVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqlelasticpool verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremssqlelasticpool %q not found after deploy", id)
	}
	return nil
}

func (*mssqlElasticPoolVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, mssqlAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremssqlelasticpool verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremssqlelasticpool %q still exists after destroy", id)
	}
	return nil
}

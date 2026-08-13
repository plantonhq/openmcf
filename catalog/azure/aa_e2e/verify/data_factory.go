package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryAPIVersion is the Microsoft.DataFactory/factories API
// line the pinned provider vendors (datafactory/2018-06-01).
const dataFactoryAPIVersion = "2018-06-01"

// dataFactoryVerifier verifies an AzureDataFactory via the generic
// ARM resources GetByID (see armResourceExists), keyed on the
// factory's ARM ID.
type dataFactoryVerifier struct{}

func (*dataFactoryVerifier) IDOutputKey() string {
	return "data_factory_id"
}

func (*dataFactoryVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactory verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactory %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactory verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactory %q still exists after destroy", id)
	}
	return nil
}

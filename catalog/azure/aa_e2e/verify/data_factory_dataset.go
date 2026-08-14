package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryDatasetVerifier verifies an AzureDataFactoryDataset via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the dataset's ARM ID -- one ID shape serves all 13 dataset types
// ({factory_id}/datasets/{name}), on the same datafactory/2018-06-01
// API line the pinned provider vendors.
type dataFactoryDatasetVerifier struct{}

func (*dataFactoryDatasetVerifier) IDOutputKey() string {
	return "dataset_id"
}

func (*dataFactoryDatasetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorydataset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactorydataset %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryDatasetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorydataset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactorydataset %q still exists after destroy", id)
	}
	return nil
}

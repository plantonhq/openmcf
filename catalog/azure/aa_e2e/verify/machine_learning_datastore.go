package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// machineLearningDatastoreVerifier verifies an
// AzureMachineLearningDatastore via the generic ARM resources GetByID
// (see armResourceExists), keyed on the datastore's full ARM ID
// (.../workspaces/{ws}/dataStores/{name}) -- one id shape for all
// three variants (blob / data-lake / file-share); which provider
// resource created it is invisible at the ARM layer, by design.
type machineLearningDatastoreVerifier struct{}

// IDOutputKey is the datastore's full ARM ID.
func (*machineLearningDatastoreVerifier) IDOutputKey() string {
	return "datastore_id"
}

func (*machineLearningDatastoreVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningdatastore verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremachinelearningdatastore %q not found after deploy", id)
	}
	return nil
}

func (*machineLearningDatastoreVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremachinelearningdatastore verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremachinelearningdatastore %q still exists after destroy", id)
	}
	return nil
}

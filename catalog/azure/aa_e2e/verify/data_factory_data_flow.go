package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryDataFlowVerifier verifies an AzureDataFactoryDataFlow via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the data flow's factory-scoped ARM ID
// ({factory_id}/dataflows/{name}) -- the SAME id shape for both the
// mapping and flowlet forms, so one verifier serves both. Data flows
// share the factory's own API line (dataFactoryAPIVersion,
// data_factory.go).
type dataFactoryDataFlowVerifier struct{}

func (*dataFactoryDataFlowVerifier) IDOutputKey() string {
	return "data_flow_id"
}

func (*dataFactoryDataFlowVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorydataflow verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactorydataflow %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryDataFlowVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorydataflow verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactorydataflow %q still exists after destroy", id)
	}
	return nil
}

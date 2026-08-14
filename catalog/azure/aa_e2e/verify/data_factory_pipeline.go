package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryPipelineVerifier verifies an AzureDataFactoryPipeline via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the pipeline's factory-scoped ARM ID
// ({factory_id}/pipelines/{name}). Pipelines share the factory's own
// API line (dataFactoryAPIVersion, data_factory.go).
type dataFactoryPipelineVerifier struct{}

func (*dataFactoryPipelineVerifier) IDOutputKey() string {
	return "pipeline_id"
}

func (*dataFactoryPipelineVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorypipeline verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactorypipeline %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryPipelineVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorypipeline verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactorypipeline %q still exists after destroy", id)
	}
	return nil
}

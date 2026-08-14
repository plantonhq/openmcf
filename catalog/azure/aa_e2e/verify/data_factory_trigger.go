package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataFactoryTriggerVerifier verifies an AzureDataFactoryTrigger via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the trigger's factory-scoped ARM ID ({factory_id}/triggers/{name})
// -- the SAME id shape for all four trigger types (schedule, tumbling
// window, blob event, custom event), so one verifier serves the whole
// union. Triggers share the factory's own API line
// (dataFactoryAPIVersion, data_factory.go).
type dataFactoryTriggerVerifier struct{}

func (*dataFactoryTriggerVerifier) IDOutputKey() string {
	return "trigger_id"
}

func (*dataFactoryTriggerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorytrigger verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredatafactorytrigger %q not found after deploy", id)
	}
	return nil
}

func (*dataFactoryTriggerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataFactoryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredatafactorytrigger verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredatafactorytrigger %q still exists after destroy", id)
	}
	return nil
}

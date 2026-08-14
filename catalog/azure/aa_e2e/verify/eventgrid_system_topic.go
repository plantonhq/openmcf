package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridSystemTopicVerifier verifies an AzureEventgridSystemTopic
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the system topic's ARM ID. Rides the family's eventgrid API pin.
type eventgridSystemTopicVerifier struct{}

func (*eventgridSystemTopicVerifier) IDOutputKey() string {
	return "system_topic_id"
}

func (*eventgridSystemTopicVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridsystemtopic verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgridsystemtopic %q not found after deploy", id)
	}
	return nil
}

func (*eventgridSystemTopicVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridsystemtopic verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgridsystemtopic %q still exists after destroy", id)
	}
	return nil
}

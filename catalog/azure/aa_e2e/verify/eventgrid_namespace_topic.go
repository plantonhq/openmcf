package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridNamespaceTopicVerifier verifies an
// AzureEventgridNamespaceTopic via the generic ARM resources GetByID
// (see armResourceExists), keyed on the topic's namespace-scoped ARM
// ID. Rides the family's eventgrid API pin (the 2025-02-15 line the
// provider vendors for namespace topics).
type eventgridNamespaceTopicVerifier struct{}

func (*eventgridNamespaceTopicVerifier) IDOutputKey() string {
	return "namespace_topic_id"
}

func (*eventgridNamespaceTopicVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridnamespacetopic verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgridnamespacetopic %q not found after deploy", id)
	}
	return nil
}

func (*eventgridNamespaceTopicVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridnamespacetopic verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgridnamespacetopic %q still exists after destroy", id)
	}
	return nil
}

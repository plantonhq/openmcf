package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridAPIVersion pins the Microsoft.EventGrid data-plane control
// API the family's verifiers read with -- the version the pinned
// azurerm provider vendors for topics, domains, and domain topics.
const eventgridAPIVersion = "2025-02-15"

// eventgridTopicVerifier verifies an AzureEventgridTopic via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// topic's ARM ID.
type eventgridTopicVerifier struct{}

func (*eventgridTopicVerifier) IDOutputKey() string {
	return "topic_id"
}

func (*eventgridTopicVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridtopic verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgridtopic %q not found after deploy", id)
	}
	return nil
}

func (*eventgridTopicVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgridtopic verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgridtopic %q still exists after destroy", id)
	}
	return nil
}

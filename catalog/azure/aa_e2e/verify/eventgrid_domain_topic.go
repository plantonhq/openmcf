package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridDomainTopicVerifier verifies an AzureEventgridDomainTopic
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the topic's domain-scoped ARM ID ({domain_id}/topics/{name}) --
// GetByID resolves child-resource ids like any other. Pinned to the
// same eventgrid 2025-02-15 line as the family's other verifiers.
type eventgridDomainTopicVerifier struct{}

func (*eventgridDomainTopicVerifier) IDOutputKey() string {
	return "domain_topic_id"
}

func (*eventgridDomainTopicVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgriddomaintopic verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgriddomaintopic %q not found after deploy", id)
	}
	return nil
}

func (*eventgridDomainTopicVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgriddomaintopic verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgriddomaintopic %q still exists after destroy", id)
	}
	return nil
}

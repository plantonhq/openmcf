package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridEventSubscriptionVerifier verifies an
// AzureEventgridEventSubscription via the generic ARM resources
// GetByID (see armResourceExists), keyed on the subscription's ARM ID.
// The id's shape follows the kind's addressing choice (a scoped
// extension id or a system-topic child id) -- GetByID serves both, so
// one verifier covers both wirings. Rides the family's eventgrid API
// pin.
type eventgridEventSubscriptionVerifier struct{}

func (*eventgridEventSubscriptionVerifier) IDOutputKey() string {
	return "event_subscription_id"
}

func (*eventgridEventSubscriptionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgrideventsubscription verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgrideventsubscription %q not found after deploy", id)
	}
	return nil
}

func (*eventgridEventSubscriptionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgrideventsubscription verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgrideventsubscription %q still exists after destroy", id)
	}
	return nil
}

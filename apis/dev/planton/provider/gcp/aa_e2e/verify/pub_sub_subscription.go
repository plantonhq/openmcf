package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// pubSubSubscriptionVerifier probes a Pub/Sub subscription via the
// Pub/Sub API. The subscription_id output is the fully qualified resource
// path (projects/{p}/subscriptions/{name}). Posture assertions confirm
// the topic attachment resolved (the FK-composition proof) and that the
// platform attribution labels landed (the label-parity proof).
type pubSubSubscriptionVerifier struct{}

func (v *pubSubSubscriptionVerifier) IDOutputKey() string { return "subscription_id" }

func (v *pubSubSubscriptionVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	subscriptionID := outputs["subscription_id"]
	if subscriptionID == "" {
		return errors.New("subscription_id output missing after deploy")
	}

	sub, err := svc.PubSub.Projects.Subscriptions.Get(subscriptionID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "pub/sub subscription %s not found after deploy", subscriptionID)
	}

	// The topic attachment is the whole point of the subscription — an
	// empty (or deleted-sentinel) topic means the FK chain did not
	// actually compose.
	if sub.Topic == "" || sub.Topic == "_deleted-topic_" {
		return errors.Errorf("pub/sub subscription %s has no live topic attachment after deploy (got %q)",
			subscriptionID, sub.Topic)
	}
	if sub.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("pub/sub subscription %s missing the planton-ai_resource attribution label after deploy", subscriptionID)
	}
	return nil
}

func (v *pubSubSubscriptionVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	subscriptionID := outputs["subscription_id"]
	if subscriptionID == "" {
		return nil
	}

	_, err := svc.PubSub.Projects.Subscriptions.Get(subscriptionID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("pub/sub subscription %s still exists after destroy", subscriptionID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing pub/sub subscription %s after destroy", subscriptionID)
}

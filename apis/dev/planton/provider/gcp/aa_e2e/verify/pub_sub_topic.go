package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// pubSubTopicVerifier probes a Pub/Sub topic via the Pub/Sub API. The
// topic_id output is the fully qualified resource path
// (projects/{p}/topics/{name}) — exactly the handle subscriptions, Cloud
// Functions event triggers, and Scheduler targets consume. Posture
// assertions confirm the platform attribution labels landed (the
// label-parity proof — both engines must stamp the identical set).
type pubSubTopicVerifier struct{}

func (v *pubSubTopicVerifier) IDOutputKey() string { return "topic_id" }

func (v *pubSubTopicVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	topicID := outputs["topic_id"]
	if topicID == "" {
		return errors.New("topic_id output missing after deploy")
	}

	topic, err := svc.PubSub.Projects.Topics.Get(topicID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "pub/sub topic %s not found after deploy", topicID)
	}

	// The platform attribution labels are the cross-engine parity canary:
	// a missing set means one engine stamped labels and the other did not.
	if topic.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("pub/sub topic %s missing the planton-ai_resource attribution label after deploy", topicID)
	}
	return nil
}

func (v *pubSubTopicVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	topicID := outputs["topic_id"]
	if topicID == "" {
		return nil
	}

	_, err := svc.PubSub.Projects.Topics.Get(topicID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("pub/sub topic %s still exists after destroy", topicID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing pub/sub topic %s after destroy", topicID)
}

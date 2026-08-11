package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// eventarcTriggerVerifier probes an Eventarc trigger and confirms it
// exists with a destination. The trigger_id output is the full resource
// name (projects/{p}/locations/{l}/triggers/{name}) — exactly what the
// Eventarc GET takes.
type eventarcTriggerVerifier struct{}

func (v *eventarcTriggerVerifier) IDOutputKey() string { return "trigger_id" }

func (v *eventarcTriggerVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["trigger_id"]
	trigger, err := svc.Eventarc.Projects.Locations.Triggers.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "eventarc trigger %s not found after deploy", name)
	}
	if trigger.Destination == nil {
		return errors.Errorf("eventarc trigger %s reports no destination", name)
	}
	return nil
}

func (v *eventarcTriggerVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["trigger_id"]
	_, err := svc.Eventarc.Projects.Locations.Triggers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing eventarc trigger %s after destroy", name)
	}
	return errors.Errorf("eventarc trigger %s still exists after destroy", name)
}

package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// eventarcMessageBusVerifier probes an Eventarc Advanced message bus and
// confirms it exists. The message_bus_name output is the full resource
// name (projects/{p}/locations/{l}/messageBuses/{id}) — exactly what the
// Eventarc GET takes. Satellites (sources, pipelines, enrollments) are
// deployed and destroyed with the bus; the bus probe is the family's
// existence anchor, and satellite posture belongs to the proof lane's
// live API reads.
type eventarcMessageBusVerifier struct{}

func (v *eventarcMessageBusVerifier) IDOutputKey() string { return "message_bus_name" }

func (v *eventarcMessageBusVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["message_bus_name"]
	if _, err := svc.Eventarc.Projects.Locations.MessageBuses.Get(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "eventarc message bus %s not found after deploy", name)
	}
	return nil
}

func (v *eventarcMessageBusVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["message_bus_name"]
	_, err := svc.Eventarc.Projects.Locations.MessageBuses.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing eventarc message bus %s after destroy", name)
	}
	return errors.Errorf("eventarc message bus %s still exists after destroy", name)
}

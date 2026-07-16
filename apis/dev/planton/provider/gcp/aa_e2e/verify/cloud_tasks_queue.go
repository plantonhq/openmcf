package verify

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudTasksQueueVerifier probes a Cloud Tasks queue via the Cloud Tasks
// API. The queue_id output is the fully qualified resource path
// (projects/{p}/locations/{l}/queues/{name}) the Queues.Get call consumes
// directly. Posture assertions confirm the queue is RUNNING and that the
// max_burst_size output matches the live GCP-computed value — the proof
// that both engines export the same effective rate-limit posture.
type cloudTasksQueueVerifier struct{}

func (v *cloudTasksQueueVerifier) IDOutputKey() string { return "queue_id" }

func (v *cloudTasksQueueVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	queueID := outputs["queue_id"]
	if queueID == "" {
		return errors.New("queue_id output missing after deploy")
	}

	queue, err := svc.CloudTasks.Projects.Locations.Queues.Get(queueID).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "cloud tasks queue %s not found after deploy", queueID)
	}

	// A freshly created queue dispatches; anything else means the deploy
	// left the queue in a broken or paused posture.
	if queue.State != "RUNNING" {
		return errors.Errorf("cloud tasks queue %s in state %s after deploy (expected RUNNING)", queueID, queue.State)
	}

	// max_burst_size is computed by GCP from the dispatch rate; the output
	// must report the live value identically on both engines.
	if got := outputs["max_burst_size"]; got != "" {
		var live int64
		if queue.RateLimits != nil {
			live = queue.RateLimits.MaxBurstSize
		}
		if got != strconv.FormatInt(live, 10) {
			return errors.Errorf("cloud tasks queue %s max_burst_size output %q does not match live value %d", queueID, got, live)
		}
	}
	return nil
}

func (v *cloudTasksQueueVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	queueID := outputs["queue_id"]
	if queueID == "" {
		return nil
	}

	_, err := svc.CloudTasks.Projects.Locations.Queues.Get(queueID).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("cloud tasks queue %s still exists after destroy", queueID)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing cloud tasks queue %s after destroy", queueID)
}

package verify

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudRunVerifier probes a Cloud Run service via the Run Admin API v2.
// Posture assertions confirm the service reconciled successfully (the Ready
// terminal condition), that the serving URL and latest ready revision in the
// outputs match live state, and that traffic actually serves 100% — proof the
// deploy produced a routable service, not just an API object.
type cloudRunVerifier struct{}

func (v *cloudRunVerifier) IDOutputKey() string { return "service_name" }

// servicePath builds the projects/{p}/locations/{region}/services/{name}
// resource path the run API addresses services by.
func (v *cloudRunVerifier) servicePath(svc *Services, outputs map[string]string) (string, error) {
	name := outputs["service_name"]
	location := outputs["location"]
	if name == "" || location == "" {
		return "", errors.New("service_name or location output missing")
	}
	return fmt.Sprintf("projects/%s/locations/%s/services/%s", svc.Project, location, name), nil
}

func (v *cloudRunVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.servicePath(svc, outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	service, err := svc.Run.Projects.Locations.Services.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "cloud run service %s not found after deploy", path)
	}

	// The Ready terminal condition is the API's own "reconciled and serving"
	// signal — a service can exist while its latest revision failed to start.
	if service.TerminalCondition == nil || service.TerminalCondition.State != "CONDITION_SUCCEEDED" {
		state := "<none>"
		if service.TerminalCondition != nil {
			state = service.TerminalCondition.State
		}
		return errors.Errorf("cloud run service %s terminal condition is %q, want CONDITION_SUCCEEDED", path, state)
	}

	if url := outputs["url"]; url != "" && service.Uri != url {
		return errors.Errorf("cloud run service %s url mismatch: output %q, live %q", path, url, service.Uri)
	}
	if revision := outputs["revision"]; revision != "" && service.LatestReadyRevision != revision {
		return errors.Errorf("cloud run service %s latest ready revision mismatch: output %q, live %q",
			path, revision, service.LatestReadyRevision)
	}

	// Traffic must actually route: the percents across all statuses sum to
	// 100 for a serving service (an empty traffic spec routes 100% LATEST).
	totalPercent := int64(0)
	for _, status := range service.TrafficStatuses {
		totalPercent += status.Percent
	}
	if totalPercent != 100 {
		return errors.Errorf("cloud run service %s traffic serves %d%%, want 100%%", path, totalPercent)
	}

	return nil
}

func (v *cloudRunVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.servicePath(svc, outputs)
	if err != nil {
		// Without a path there is nothing left to probe; treat as gone.
		return nil
	}

	_, err = svc.Run.Projects.Locations.Services.Get(path).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("cloud run service %s still exists after destroy", path)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing cloud run service %s after destroy", path)
}

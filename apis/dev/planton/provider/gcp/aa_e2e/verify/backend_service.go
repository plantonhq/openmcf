package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// backendServiceVerifier probes a global Compute Engine backend service by
// name via the compute API, additionally confirming its health-check wiring —
// the FK-resolved reference to the GcpHealthCheck prerequisite is the point
// of the composition.
type backendServiceVerifier struct{}

func (v *backendServiceVerifier) IDOutputKey() string { return "self_link" }

func (v *backendServiceVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["backend_service_name"]
	backendService, err := svc.Compute.BackendServices.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "backend service %s not found after deploy", name)
	}
	// The scenario wires exactly one health check by reference; the deployed
	// service must carry it. (GCP itself caps the list at one.)
	if len(backendService.HealthChecks) != 1 {
		return errors.Errorf("backend service %s carries %d health checks, expected exactly 1",
			name, len(backendService.HealthChecks))
	}
	return nil
}

func (v *backendServiceVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["backend_service_name"]
	_, err := svc.Compute.BackendServices.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing backend service %s after destroy", name)
	}
	return errors.Errorf("backend service %s still exists after destroy", name)
}

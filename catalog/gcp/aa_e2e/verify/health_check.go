package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// healthCheckVerifier probes a Compute Engine health check by the self_link
// output. One kind maps to two API collections (global and regional health
// checks), so the verifier routes on the self-link's shape rather than
// requiring the caller to know the scope.
type healthCheckVerifier struct{}

func (v *healthCheckVerifier) IDOutputKey() string { return "self_link" }

func (v *healthCheckVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["health_check_name"]
	region := outputs["region"]

	if region != "" {
		if _, err := svc.Compute.RegionHealthChecks.Get(svc.Project, region, name).Context(ctx).Do(); err != nil {
			return errors.Wrapf(err, "regional health check %s/%s not found after deploy", region, name)
		}
		return nil
	}
	if _, err := svc.Compute.HealthChecks.Get(svc.Project, name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "global health check %s not found after deploy", name)
	}
	return nil
}

func (v *healthCheckVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["health_check_name"]
	region := outputs["region"]

	var err error
	if region != "" {
		_, err = svc.Compute.RegionHealthChecks.Get(svc.Project, region, name).Context(ctx).Do()
	} else {
		_, err = svc.Compute.HealthChecks.Get(svc.Project, name).Context(ctx).Do()
	}
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing health check %s after destroy", name)
	}
	return errors.Errorf("health check %s still exists after destroy", strings.TrimSpace(name))
}

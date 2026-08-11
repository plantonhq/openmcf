package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// monitoringSloVerifier probes a Cloud Monitoring SLO by its
// server-assigned resource name and confirms the objective exists with a
// real goal. Both outputs derive from the SLO's own name, so the probe is
// correct on every service arm (existing, custom, or basic service); the
// created service's existence is implied — an SLO cannot exist under a
// service that does not.
type monitoringSloVerifier struct{}

func (v *monitoringSloVerifier) IDOutputKey() string { return "slo_name" }

func (v *monitoringSloVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["slo_name"]
	slo, err := svc.Monitoring.Services.ServiceLevelObjectives.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "slo %s not found after deploy", name)
	}
	if slo.Goal == 0 {
		return errors.Errorf("slo %s reports no goal", name)
	}
	return nil
}

func (v *monitoringSloVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["slo_name"]
	_, err := svc.Monitoring.Services.ServiceLevelObjectives.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing slo %s after destroy", name)
	}
	return errors.Errorf("slo %s still exists after destroy", name)
}

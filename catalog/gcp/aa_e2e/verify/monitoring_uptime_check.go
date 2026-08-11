package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// monitoringUptimeCheckVerifier probes a Cloud Monitoring uptime check by
// its full resource name and confirms a probe surface landed (an HTTP or
// TCP check block) — the config GCP stores must carry the check arm the
// scenario configured, not just exist as an empty shell.
type monitoringUptimeCheckVerifier struct{}

func (v *monitoringUptimeCheckVerifier) IDOutputKey() string { return "uptime_check_name" }

func (v *monitoringUptimeCheckVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["uptime_check_name"]
	config, err := svc.Monitoring.Projects.UptimeCheckConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "uptime check %s not found after deploy", name)
	}
	if config.HttpCheck == nil && config.TcpCheck == nil && config.SyntheticMonitor == nil {
		return errors.Errorf("uptime check %s carries no probe surface (http/tcp/synthetic)", name)
	}
	return nil
}

func (v *monitoringUptimeCheckVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["uptime_check_name"]
	_, err := svc.Monitoring.Projects.UptimeCheckConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing uptime check %s after destroy", name)
	}
	return errors.Errorf("uptime check %s still exists after destroy", name)
}

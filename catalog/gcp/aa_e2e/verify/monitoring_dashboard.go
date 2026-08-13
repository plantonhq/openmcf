package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// monitoringDashboardVerifier probes a Cloud Monitoring dashboard by its
// server-assigned resource name and confirms the JSON document became a
// real dashboard with a display name. Dashboards live on the Monitoring
// API's v1 surface (a different API version from the v3 client the other
// monitoring verifiers use).
type monitoringDashboardVerifier struct{}

func (v *monitoringDashboardVerifier) IDOutputKey() string { return "dashboard_name" }

func (v *monitoringDashboardVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["dashboard_name"]
	dashboard, err := svc.MonitoringDashboards.Projects.Dashboards.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "dashboard %s not found after deploy", name)
	}
	if dashboard.DisplayName == "" {
		return errors.Errorf("dashboard %s reports no display name", name)
	}
	return nil
}

func (v *monitoringDashboardVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["dashboard_name"]
	_, err := svc.MonitoringDashboards.Projects.Dashboards.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing dashboard %s after destroy", name)
	}
	return errors.Errorf("dashboard %s still exists after destroy", name)
}

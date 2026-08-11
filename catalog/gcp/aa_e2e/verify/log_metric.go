package verify

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// logMetricVerifier probes a Cloud Logging log-based metric and confirms
// it exists with a real filter. The metric_name output is the bare metric
// name (the segment after "metrics/"); it may contain forward slashes
// (namespaced metrics like "checkout/errors"), which the GET path
// requires URL-encoded — the probe escapes it explicitly rather than
// trusting path substitution.
type logMetricVerifier struct{}

func (v *logMetricVerifier) IDOutputKey() string { return "metric_name" }

func (v *logMetricVerifier) metricPath(svc *Services, outputs map[string]string) string {
	return fmt.Sprintf("projects/%s/metrics/%s", svc.Project, url.PathEscape(outputs["metric_name"]))
}

func (v *logMetricVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := v.metricPath(svc, outputs)
	metric, err := svc.Logging.Projects.Metrics.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "log metric %s not found after deploy", name)
	}
	if metric.Filter == "" {
		return errors.Errorf("log metric %s reports no filter", name)
	}
	return nil
}

func (v *logMetricVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := v.metricPath(svc, outputs)
	_, err := svc.Logging.Projects.Metrics.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing log metric %s after destroy", name)
	}
	return errors.Errorf("log metric %s still exists after destroy", name)
}

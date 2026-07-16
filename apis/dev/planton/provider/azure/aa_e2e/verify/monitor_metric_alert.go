package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorMetricAlertAPIVersion pins the Microsoft.Insights metricAlerts RP
// version the verifier reads with.
const monitorMetricAlertAPIVersion = "2018-03-01"

// monitorMetricAlertVerifier verifies an AzureMonitorMetricAlert via the
// generic ARM resources GetByID (see armResourceExists), keyed on the alert
// rule's ARM ID.
type monitorMetricAlertVerifier struct{}

func (*monitorMetricAlertVerifier) IDOutputKey() string {
	return "metric_alert_id"
}

func (*monitorMetricAlertVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorMetricAlertAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitormetricalert verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitormetricalert %q not found after deploy", id)
	}
	return nil
}

func (*monitorMetricAlertVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorMetricAlertAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitormetricalert verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitormetricalert %q still exists after destroy", id)
	}
	return nil
}

package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorActivityLogAlertAPIVersion is the stable Microsoft.Insights API
// version the generic existence probe is pinned to.
const monitorActivityLogAlertAPIVersion = "2020-10-01"

// monitorActivityLogAlertVerifier verifies an AzureMonitorActivityLogAlert
// via the generic ARM resources GetByID (see armResourceExists), keyed on
// the alert's full ARM ID.
type monitorActivityLogAlertVerifier struct{}

// IDOutputKey is the alert's full ARM ID.
func (*monitorActivityLogAlertVerifier) IDOutputKey() string {
	return "activity_log_alert_id"
}

func (*monitorActivityLogAlertVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorActivityLogAlertAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitoractivitylogalert verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitoractivitylogalert %q not found after deploy", id)
	}
	return nil
}

func (*monitorActivityLogAlertVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorActivityLogAlertAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitoractivitylogalert verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitoractivitylogalert %q still exists after destroy", id)
	}
	return nil
}

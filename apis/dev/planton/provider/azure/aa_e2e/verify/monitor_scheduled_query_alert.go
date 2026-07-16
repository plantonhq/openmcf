package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorScheduledQueryAlertAPIVersion pins the Microsoft.Insights
// scheduledQueryRules RP version the verifier reads with.
const monitorScheduledQueryAlertAPIVersion = "2023-03-15-preview"

// monitorScheduledQueryAlertVerifier verifies an
// AzureMonitorScheduledQueryAlert via the generic ARM resources GetByID
// (see armResourceExists), keyed on the rule's ARM ID.
type monitorScheduledQueryAlertVerifier struct{}

func (*monitorScheduledQueryAlertVerifier) IDOutputKey() string {
	return "scheduled_query_alert_id"
}

func (*monitorScheduledQueryAlertVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorScheduledQueryAlertAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitorscheduledqueryalert verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitorscheduledqueryalert %q not found after deploy", id)
	}
	return nil
}

func (*monitorScheduledQueryAlertVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorScheduledQueryAlertAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitorscheduledqueryalert verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitorscheduledqueryalert %q still exists after destroy", id)
	}
	return nil
}

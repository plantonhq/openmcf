package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorAutoscaleSettingAPIVersion pins the Microsoft.Insights
// autoScaleSettings RP version the verifier reads with -- the line the
// provider vendors (insights/2022-10-01/autoscalesettings).
const monitorAutoscaleSettingAPIVersion = "2022-10-01"

// monitorAutoscaleSettingVerifier verifies an AzureMonitorAutoscaleSetting
// via the generic ARM resources GetByID (see armResourceExists), keyed on
// the setting's ARM ID.
type monitorAutoscaleSettingVerifier struct{}

func (*monitorAutoscaleSettingVerifier) IDOutputKey() string {
	return "autoscale_setting_id"
}

func (*monitorAutoscaleSettingVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorAutoscaleSettingAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitorautoscalesetting verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitorautoscalesetting %q not found after deploy", id)
	}
	return nil
}

func (*monitorAutoscaleSettingVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorAutoscaleSettingAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitorautoscalesetting verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitorautoscalesetting %q still exists after destroy", id)
	}
	return nil
}

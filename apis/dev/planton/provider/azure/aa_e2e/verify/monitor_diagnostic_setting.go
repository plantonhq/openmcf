package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorDiagnosticSettingAPIVersion pins the Microsoft.Insights
// diagnosticSettings RP version the verifier reads with.
const monitorDiagnosticSettingAPIVersion = "2021-05-01-preview"

// monitorDiagnosticSettingVerifier verifies an AzureMonitorDiagnosticSetting
// via the generic ARM resources GetByID (see armResourceExists). The module
// exports the CONSTRUCTED extension-resource ARM ID
// ({target}/providers/Microsoft.Insights/diagnosticSettings/{name}) --
// never the provider's "{target}|{name}" composite state id, which no API
// consumes.
type monitorDiagnosticSettingVerifier struct{}

func (*monitorDiagnosticSettingVerifier) IDOutputKey() string {
	return "diagnostic_setting_id"
}

func (*monitorDiagnosticSettingVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorDiagnosticSettingAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitordiagnosticsetting verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitordiagnosticsetting %q not found after deploy", id)
	}
	return nil
}

func (*monitorDiagnosticSettingVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorDiagnosticSettingAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitordiagnosticsetting verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitordiagnosticsetting %q still exists after destroy", id)
	}
	return nil
}

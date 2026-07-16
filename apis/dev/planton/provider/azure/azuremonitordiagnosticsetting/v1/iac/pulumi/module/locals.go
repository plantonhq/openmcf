package module

import (
	azuremonitordiagnosticsettingv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremonitordiagnosticsetting/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMonitorDiagnosticSetting *azuremonitordiagnosticsettingv1.AzureMonitorDiagnosticSetting
}

// destinationTypeStrings maps the workspace-table-layout enum to ARM's wire
// values.
var destinationTypeStrings = map[azuremonitordiagnosticsettingv1.AzureMonitorDiagnosticSettingLogAnalyticsDestinationType]string{
	azuremonitordiagnosticsettingv1.AzureMonitorDiagnosticSettingLogAnalyticsDestinationType_DEDICATED:         "Dedicated",
	azuremonitordiagnosticsettingv1.AzureMonitorDiagnosticSettingLogAnalyticsDestinationType_AZURE_DIAGNOSTICS: "AzureDiagnostics",
}

// The diagnostic setting carries no tags (the ARM extension resource does
// not support them), so locals stay minimal.
func initializeLocals(ctx *pulumi.Context, stackInput *azuremonitordiagnosticsettingv1.AzureMonitorDiagnosticSettingStackInput) *Locals {
	return &Locals{
		AzureMonitorDiagnosticSetting: stackInput.Target,
	}
}

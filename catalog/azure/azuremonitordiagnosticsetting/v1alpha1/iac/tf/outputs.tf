# The provider's own state id is a "{target}|{name}" composite, not an
# ARM id -- the output constructs the real ARM extension-resource id so
# downstream consumers and verification address the setting the way the
# API does.
output "diagnostic_setting_id" {
  description = "The ARM extension-resource ID of the diagnostic setting"
  value       = "${var.spec.target_resource_id}/providers/Microsoft.Insights/diagnosticSettings/${var.spec.setting_name}"
}

output "diagnostic_setting_name" {
  description = "The name of the diagnostic setting"
  value       = azurerm_monitor_diagnostic_setting.main.name
}

output "target_resource_id" {
  description = "The ARM ID of the target resource the setting routes telemetry from"
  value       = azurerm_monitor_diagnostic_setting.main.target_resource_id
}

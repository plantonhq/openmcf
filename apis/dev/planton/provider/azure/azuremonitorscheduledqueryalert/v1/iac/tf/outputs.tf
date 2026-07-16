output "scheduled_query_alert_id" {
  description = "The Azure Resource Manager ID of the scheduled query alert rule"
  value       = azurerm_monitor_scheduled_query_rules_alert_v2.main.id
}

output "scheduled_query_alert_name" {
  description = "The name of the scheduled query alert rule"
  value       = azurerm_monitor_scheduled_query_rules_alert_v2.main.name
}

output "identity_principal_id" {
  description = "The principal ID of the system-assigned managed identity (empty unless SYSTEM_ASSIGNED is enabled)"
  value       = try(azurerm_monitor_scheduled_query_rules_alert_v2.main.identity[0].principal_id, "")
}

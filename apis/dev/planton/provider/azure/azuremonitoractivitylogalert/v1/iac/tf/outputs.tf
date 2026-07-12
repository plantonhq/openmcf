output "activity_log_alert_id" {
  description = "The Azure Resource Manager ID of the activity log alert"
  value       = azurerm_monitor_activity_log_alert.main.id
}

output "activity_log_alert_name" {
  description = "The name of the activity log alert resource"
  value       = azurerm_monitor_activity_log_alert.main.name
}

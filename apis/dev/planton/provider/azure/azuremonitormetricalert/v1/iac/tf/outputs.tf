output "metric_alert_id" {
  description = "The Azure Resource Manager ID of the metric alert rule"
  value       = azurerm_monitor_metric_alert.main.id
}

output "metric_alert_name" {
  description = "The name of the metric alert rule"
  value       = azurerm_monitor_metric_alert.main.name
}

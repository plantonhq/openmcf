output "autoscale_setting_id" {
  description = "The Azure Resource Manager ID of the autoscale setting"
  value       = azurerm_monitor_autoscale_setting.main.id
}

output "autoscale_setting_name" {
  description = "The name of the autoscale setting resource"
  value       = azurerm_monitor_autoscale_setting.main.name
}

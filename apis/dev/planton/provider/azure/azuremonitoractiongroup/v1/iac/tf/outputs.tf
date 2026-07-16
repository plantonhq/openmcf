output "action_group_id" {
  description = "The Azure Resource Manager ID of the action group -- the FK seam alert rules reference"
  value       = azurerm_monitor_action_group.main.id
}

output "action_group_name" {
  description = "The name of the action group"
  value       = azurerm_monitor_action_group.main.name
}

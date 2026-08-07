output "dapr_component_id" {
  description = "The Azure Resource Manager ID of the Dapr component"
  value       = azurerm_container_app_environment_dapr_component.main.id
}

output "component_name" {
  description = "The Dapr component name application code passes to the Dapr API"
  value       = azurerm_container_app_environment_dapr_component.main.name
}

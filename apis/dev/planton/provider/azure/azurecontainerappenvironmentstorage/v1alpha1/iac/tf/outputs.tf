output "storage_id" {
  description = "The Azure Resource Manager ID of the storage registration"
  value       = azurerm_container_app_environment_storage.main.id
}

output "storage_name" {
  description = "The registration name app and job volumes reference in storage_name"
  value       = azurerm_container_app_environment_storage.main.name
}

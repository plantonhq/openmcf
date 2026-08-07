output "container_registry_id" {
  description = "The Azure Resource Manager ID of the registry"
  value       = azurerm_container_registry.main.id
}

output "container_registry_name" {
  description = "The name of the registry"
  value       = azurerm_container_registry.main.name
}

output "login_server" {
  description = "The registry's login server hostname, e.g. myregistry.azurecr.io"
  value       = azurerm_container_registry.main.login_server
}

output "admin_username" {
  description = "The admin account's username (populated only when admin_user_enabled is true)"
  value       = azurerm_container_registry.main.admin_username
}

output "admin_password" {
  description = "One of the admin account's two rotatable passwords (populated only when admin_user_enabled is true)"
  value       = azurerm_container_registry.main.admin_password
  sensitive   = true
}

output "system_assigned_identity_principal_id" {
  description = "The principal ID of the registry's system-assigned identity (populated only when the identity type includes SYSTEM_ASSIGNED)"
  value       = try(azurerm_container_registry.main.identity[0].principal_id, "")
}

output "data_endpoint_host_names" {
  description = "The dedicated regional data-endpoint hostnames (populated only when data_endpoint_enabled is true)"
  value       = azurerm_container_registry.main.data_endpoint_host_names
}

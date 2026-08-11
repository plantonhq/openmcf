output "resource_guard_id" {
  description = "The Azure Resource Manager ID of the guard -- what backup vaults reference to enable Multi-User Authorization"
  value       = azurerm_data_protection_resource_guard.main.id
}

output "resource_guard_name" {
  description = "The guard's name"
  value       = azurerm_data_protection_resource_guard.main.name
}

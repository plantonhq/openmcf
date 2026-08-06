output "disaster_recovery_config_id" {
  description = "The Azure Resource Manager ID of the disaster-recovery config"
  value       = azurerm_servicebus_namespace_disaster_recovery_config.main.id
}

output "alias_name" {
  description = "The failover-stable alias DNS identity"
  value       = azurerm_servicebus_namespace_disaster_recovery_config.main.name
}

# The alias connection strings are what DR-aware clients hold: they
# address the alias DNS name, so a failover needs no client
# reconfiguration.
output "primary_connection_string_alias" {
  description = "The primary connection string addressing the alias"
  value       = azurerm_servicebus_namespace_disaster_recovery_config.main.primary_connection_string_alias
  sensitive   = true
}

output "secondary_connection_string_alias" {
  description = "The secondary alias connection string (rotation partner)"
  value       = azurerm_servicebus_namespace_disaster_recovery_config.main.secondary_connection_string_alias
  sensitive   = true
}

output "default_primary_key" {
  description = "The paired rule's primary key"
  value       = azurerm_servicebus_namespace_disaster_recovery_config.main.default_primary_key
  sensitive   = true
}

output "default_secondary_key" {
  description = "The paired rule's secondary key (rotation partner)"
  value       = azurerm_servicebus_namespace_disaster_recovery_config.main.default_secondary_key
  sensitive   = true
}

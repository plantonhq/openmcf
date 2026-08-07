output "disaster_recovery_config_id" {
  description = "The Azure Resource Manager ID of the disaster-recovery config"
  value       = azurerm_eventhub_namespace_disaster_recovery_config.main.id
}

# No credential outputs here: Azure's Event Hubs DR resource exposes
# none. Alias-addressed connection strings surface on the namespace and
# authorization-rule kinds instead.
output "alias_name" {
  description = "The failover-stable alias DNS identity"
  value       = azurerm_eventhub_namespace_disaster_recovery_config.main.name
}

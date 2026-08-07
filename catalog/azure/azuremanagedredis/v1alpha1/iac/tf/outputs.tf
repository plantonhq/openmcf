output "managed_redis_id" {
  description = "The Azure Resource Manager ID of the Managed Redis cluster"
  value       = azurerm_managed_redis.main.id
}

output "managed_redis_name" {
  description = "The name of the Managed Redis instance"
  value       = azurerm_managed_redis.main.name
}

output "region" {
  description = "The Azure region the instance lives in"
  value       = azurerm_managed_redis.main.location
}

output "resource_group_name" {
  description = "The resource group the instance lives in"
  value       = azurerm_managed_redis.main.resource_group_name
}

# Keyless (Entra) clients need only the hostname and port -- tokens
# replace keys.
output "hostname" {
  description = "The DNS hostname of the Managed Redis instance"
  value       = azurerm_managed_redis.main.hostname
}

# The scope Entra access-policy grants and geo-replication links operate
# on: the cluster ID with /databases/default appended.
output "database_id" {
  description = "The Azure Resource Manager ID of the default database"
  value       = try(azurerm_managed_redis.main.default_database[0].id, "")
}

output "port" {
  description = "The TCP port the database listens on (10000)"
  value       = try(azurerm_managed_redis.main.default_database[0].port, 0)
}

# The keys are secret-bearing: they are the database password. Both stay
# empty under the keyless default (access keys disabled).
output "primary_access_key" {
  description = "The primary access key; empty while access-keys authentication is disabled"
  value       = try(azurerm_managed_redis.main.default_database[0].primary_access_key, "")
  sensitive   = true
}

# Kept live for zero-downtime rotation: move clients here while the
# primary is regenerated, and vice versa.
output "secondary_access_key" {
  description = "The secondary access key; empty while access-keys authentication is disabled"
  value       = try(azurerm_managed_redis.main.default_database[0].secondary_access_key, "")
  sensitive   = true
}

# Populated only when the identity block requests SYSTEM_ASSIGNED.
output "identity_principal_id" {
  description = "The system-assigned identity's principal ID, when enabled"
  value       = try(azurerm_managed_redis.main.identity[0].principal_id, "")
}

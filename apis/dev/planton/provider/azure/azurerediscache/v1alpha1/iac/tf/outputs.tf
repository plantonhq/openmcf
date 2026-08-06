output "redis_cache_id" {
  description = "The Azure Resource Manager ID of the Redis cache"
  value       = azurerm_redis_cache.main.id
}

output "redis_cache_name" {
  description = "The name of the Redis cache"
  value       = azurerm_redis_cache.main.name
}

# The linked-server location seam: AzureRedisLinkedServer references this
# so geo-replication composes without hand-repeating the region.
output "region" {
  description = "The Azure region the cache lives in"
  value       = azurerm_redis_cache.main.location
}

output "resource_group_name" {
  description = "The resource group the cache lives in"
  value       = azurerm_redis_cache.main.resource_group_name
}

# Keyless (Entra) clients need only the hostname -- tokens replace keys.
output "hostname" {
  description = "The DNS hostname of the Redis cache"
  value       = azurerm_redis_cache.main.hostname
}

output "port" {
  description = "The plaintext non-SSL port (6379); only open when non_ssl_port_enabled"
  value       = azurerm_redis_cache.main.port
}

output "ssl_port" {
  description = "The TLS port (6380) every production client should use"
  value       = azurerm_redis_cache.main.ssl_port
}

# The keys and connection strings are secret-bearing: they are the cache
# password. Empty when access-keys authentication is disabled.
output "primary_access_key" {
  description = "The primary access key for authentication"
  value       = azurerm_redis_cache.main.primary_access_key
  sensitive   = true
}

# Kept live for zero-downtime rotation: move clients here while the
# primary is regenerated, and vice versa.
output "secondary_access_key" {
  description = "The secondary access key for authentication"
  value       = azurerm_redis_cache.main.secondary_access_key
  sensitive   = true
}

output "primary_connection_string" {
  description = "Ready-to-use primary connection string (embeds the primary key)"
  value       = azurerm_redis_cache.main.primary_connection_string
  sensitive   = true
}

output "secondary_connection_string" {
  description = "Ready-to-use secondary connection string (embeds the secondary key)"
  value       = azurerm_redis_cache.main.secondary_connection_string
  sensitive   = true
}

# Populated only when the identity block requests SYSTEM_ASSIGNED -- the
# principal RBAC grants target (e.g. Storage Blob Data Contributor for
# managed-identity persistence).
output "identity_principal_id" {
  description = "The system-assigned identity's principal ID, when enabled"
  value       = try(azurerm_redis_cache.main.identity[0].principal_id, "")
}

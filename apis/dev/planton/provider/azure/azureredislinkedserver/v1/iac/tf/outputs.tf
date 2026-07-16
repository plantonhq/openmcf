output "linked_server_id" {
  description = "The Azure Resource Manager ID of the linked-server resource"
  value       = azurerm_redis_linked_server.main.id
}

# Azure names the link after the LINKED (secondary) cache.
output "linked_server_name" {
  description = "The link's name (equals the secondary cache's name)"
  value       = azurerm_redis_linked_server.main.name
}

# The DNS name that always resolves to the CURRENT primary of the pair --
# what applications should point at to survive failovers without a
# connection-string change.
output "geo_replicated_primary_host_name" {
  description = "The geo-replication hostname that follows the current primary"
  value       = azurerm_redis_linked_server.main.geo_replicated_primary_host_name
}

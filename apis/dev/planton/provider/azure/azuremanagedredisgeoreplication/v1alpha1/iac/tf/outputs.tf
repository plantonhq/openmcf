# The group has no ARM object of its own -- its resource ID is the
# managing cluster's ARM ID; membership lives on every member's default
# database.
output "geo_replication_id" {
  description = "The group's resource ID (the managing Managed Redis cluster's ARM ID)"
  value       = azurerm_managed_redis_geo_replication.main.id
}

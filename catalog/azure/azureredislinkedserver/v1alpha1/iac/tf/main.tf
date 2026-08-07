# The geo-replication link between two PREMIUM caches. Azure names the
# link after the LINKED (secondary) cache -- there is no name argument.
# Every argument is ForceNew; replacing the link re-establishes
# replication without touching cached data on the primary.
#
# DELETING this resource IS the failover operation: unlinking makes the
# secondary writable. To fail over, destroy the link, repoint traffic,
# and create a new link in the opposite direction once the region
# recovers. Applications that point at the
# geo_replicated_primary_host_name output (instead of either cache's own
# hostname) keep working across failovers without a config change.
#
# Azure's requirements, enforced at link time: both caches PREMIUM, in
# different regions, and the secondary at least as large as the primary.
# Establishing the link takes several minutes on top of the caches' own
# provisioning; the secondary rejects writes while linked. No tags: ARM
# does not support tags on linked servers.
resource "azurerm_redis_linked_server" "main" {
  target_redis_cache_name = local.target_cache_name
  resource_group_name     = local.target_resource_group_name

  linked_redis_cache_id       = var.spec.linked_redis_cache_id
  linked_redis_cache_location = var.spec.linked_redis_cache_location

  server_role = local.server_role
}

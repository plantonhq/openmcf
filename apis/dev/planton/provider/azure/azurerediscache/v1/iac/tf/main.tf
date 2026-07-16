# The cache itself. Azure sizes it as {family}{capacity} within the tier;
# the family letter is derived from the tier in locals so the spec never
# spells the same fact twice. Provisioning is the slowest in the Azure
# catalog -- 15-40 minutes is normal, and azurerm's own timeout is 3 hours.
resource "azurerm_redis_cache" "main" {
  name                = var.spec.cache_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  sku_name = local.sku_name
  family   = local.family
  capacity = var.spec.capacity

  redis_version = var.spec.redis_version

  # The plaintext port sends commands AND keys unencrypted -- off unless a
  # legacy client genuinely cannot speak TLS.
  non_ssl_port_enabled = var.spec.non_ssl_port_enabled

  # false forces all traffic through Private Link (AzurePrivateEndpoint).
  public_network_access_enabled = var.spec.public_network_access_enabled

  # The keyless posture: keys can only be turned off once Entra auth is on
  # (spec-enforced, mirroring ARM's own contract); identities then connect
  # with tokens under access-policy assignments.
  access_keys_authentication_enabled = var.spec.access_keys_authentication_enabled

  # VNet injection (Premium): the cache gets a private IP inside a subnet
  # dedicated to Redis. The legacy isolation mechanism -- Private Link is
  # the recommendation for new designs; both are modeled. Fixed at
  # creation.
  subnet_id                 = var.spec.subnet_id
  private_static_ip_address = var.spec.private_static_ip_address

  # Zone pinning for datacenter-failure resilience. Fixed at creation.
  zones = length(var.spec.zones) > 0 ? var.spec.zones : null

  # Clustering (Premium): each shard is a primary/replica pair, so memory
  # and throughput scale with the shard count. Mutually exclusive with
  # extra replicas (spec-enforced).
  shard_count = var.spec.shard_count

  # Extra read replicas per primary (Premium). Only ARM's modern name is
  # modeled; the legacy replicas_per_master alias mirrors it server-side.
  replicas_per_primary = var.spec.replicas_per_primary

  # Tenant-level platform settings -- distinct from redis_configuration
  # (the Redis engine's own settings); used by support scenarios.
  tenant_settings = length(var.spec.tenant_settings) > 0 ? var.spec.tenant_settings : null

  # The managed identity, used for keyless persistence-storage access
  # (data_persistence_authentication_method = MANAGED_IDENTITY).
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.user_assigned_identity_ids) > 0 ? identity.value.user_assigned_identity_ids : null
    }
  }

  # Redis engine and platform behavior. The block is emitted only when the
  # spec carries it, so an omitted block deploys Azure's defaults
  # identically on both engines. Null attributes inside the block are not
  # sent -- Azure sizes the memory dials from total memory.
  dynamic "redis_configuration" {
    for_each = local.redis_configuration != null ? [local.redis_configuration] : []
    content {
      active_directory_authentication_enabled = redis_configuration.value.active_directory_authentication_enabled
      maxmemory_policy                        = redis_configuration.value.maxmemory_policy
      maxmemory_reserved                      = redis_configuration.value.maxmemory_reserved
      maxmemory_delta                         = redis_configuration.value.maxmemory_delta
      maxfragmentationmemory_reserved         = redis_configuration.value.maxfragmentationmemory_reserved
      notify_keyspace_events                  = redis_configuration.value.notify_keyspace_events

      # false is only legal inside a VNet-injected cache (spec-enforced);
      # azurerm only transmits the setting when subnet_id is set.
      authentication_enabled = redis_configuration.value.authentication_enabled

      data_persistence_authentication_method = (
        redis_configuration.value.data_persistence_authentication_method != null
        ? local.persistence_auth_map[redis_configuration.value.data_persistence_authentication_method]
        : null
      )

      # RDB snapshots (Premium): periodic full dumps to a storage account.
      # The connection string is secret-bearing and never echoed back by
      # ARM.
      rdb_backup_enabled            = redis_configuration.value.rdb_backup_enabled
      rdb_backup_frequency          = redis_configuration.value.rdb_backup_frequency
      rdb_backup_max_snapshot_count = redis_configuration.value.rdb_backup_max_snapshot_count
      rdb_storage_connection_string = redis_configuration.value.rdb_storage_connection_string

      # AOF log (Premium): near-synchronous write logging for tight
      # recovery points; Azure alternates between the two accounts during
      # storage maintenance.
      aof_backup_enabled              = redis_configuration.value.aof_backup_enabled
      aof_storage_connection_string_0 = redis_configuration.value.aof_storage_connection_string_0
      aof_storage_connection_string_1 = redis_configuration.value.aof_storage_connection_string_1

      storage_account_subscription_id = redis_configuration.value.storage_account_subscription_id
    }
  }

  # Weekly maintenance windows during which Azure may patch the Redis
  # engine and platform.
  dynamic "patch_schedule" {
    for_each = var.spec.patch_schedules
    content {
      day_of_week        = local.day_of_week_map[patch_schedule.value.day_of_week]
      start_hour_utc     = patch_schedule.value.start_hour_utc
      maintenance_window = patch_schedule.value.maintenance_window
    }
  }

  tags = local.final_tags
}

# Public-endpoint IP allow-list. One ARM sub-resource per rule; only
# effective while public network access is on and the cache is not
# VNet-injected. ARM rejects hyphens in rule names (spec-enforced).
resource "azurerm_redis_firewall_rule" "rules" {
  for_each = { for rule in var.spec.firewall_rules : rule.name => rule }

  name                = each.value.name
  redis_cache_name    = azurerm_redis_cache.main.name
  resource_group_name = var.spec.resource_group
  start_ip            = each.value.start_ip
  end_ip              = each.value.end_ip
}

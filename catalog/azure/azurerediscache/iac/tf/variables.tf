variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Cache for Redis specification"
  type = object({
    # The Azure region the cache lives in.
    region = string

    # The resource group the cache lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The cache's name: 1-63 letters/digits/hyphens, globally unique (it
    # becomes {cache_name}.redis.cache.windows.net).
    cache_name = string

    # The pricing tier, as the spec enum's name string (BASIC / STANDARD /
    # PREMIUM). Absent means STANDARD -- the tfvars wire format drops
    # zero-valued proto fields, so the module materializes the default.
    sku_name = optional(string)

    # Cache size within the tier's family: C0-C6 for Basic/Standard,
    # P1-P5 for Premium. Zero is meaningful (C0), so the attribute
    # defaults to 0 rather than null.
    capacity = optional(number, 0)

    # Redis engine major version ("4" or "6").
    redis_version = optional(string, "6")

    # Availability zones to pin the cache's nodes to. Fixed at creation.
    zones = optional(list(string), [])

    # PREMIUM ONLY -- the dedicated subnet for VNet injection. References
    # are resolved to a literal ARM ID by the platform before the module
    # runs. Fixed at creation.
    subnet_id = optional(string)

    # The static private IP inside the injected subnet. Only meaningful
    # with subnet_id. Fixed at creation.
    private_static_ip_address = optional(string)

    # PREMIUM ONLY -- shard count for a clustered cache (1-10).
    shard_count = optional(number)

    # PREMIUM ONLY -- extra read replicas per primary (1-3). Only the
    # modern ARM name is modeled; the legacy replicas_per_master alias
    # mirrors it server-side.
    replicas_per_primary = optional(number)

    # Enable the plaintext non-SSL port (6379). Off by default.
    non_ssl_port_enabled = optional(bool, false)

    # Whether the cache answers on its public endpoint.
    public_network_access_enabled = optional(bool, true)

    # Whether the shared access keys authenticate clients. Can only be
    # false once Entra auth is on (enforced in the spec).
    access_keys_authentication_enabled = optional(bool, true)

    # Redis engine and platform behavior. Every field has a safe Azure
    # default; the whole block may be absent.
    redis_configuration = optional(object({
      # Microsoft Entra token authentication.
      active_directory_authentication_enabled = optional(bool, false)

      # Eviction policy, in Redis's own configuration vocabulary.
      maxmemory_policy = optional(string, "volatile-lru")

      # Memory reservations in MB. Azure sizes defaults from total
      # memory; not configurable on the BASIC tier (enforced in the
      # spec).
      maxmemory_reserved              = optional(number)
      maxmemory_delta                 = optional(number)
      maxfragmentationmemory_reserved = optional(number)

      # Keyspace event notification classes, in Redis's flag syntax.
      notify_keyspace_events = optional(string)

      # Whether Redis requires authentication at all. false only on a
      # VNet-injected cache (enforced in the spec).
      authentication_enabled = optional(bool, true)

      # How the cache authenticates to persistence storage, as the spec
      # enum's name string (SAS / MANAGED_IDENTITY).
      data_persistence_authentication_method = optional(string)

      # PREMIUM ONLY -- RDB snapshot persistence. The connection string
      # is secret-bearing.
      rdb_backup_enabled            = optional(bool, false)
      rdb_backup_frequency          = optional(number)
      rdb_backup_max_snapshot_count = optional(number)
      rdb_storage_connection_string = optional(string)

      # PREMIUM ONLY -- AOF (append-only file) persistence. Both
      # connection strings are secret-bearing; Azure alternates between
      # the two storage accounts during maintenance.
      aof_backup_enabled              = optional(bool, false)
      aof_storage_connection_string_0 = optional(string)
      aof_storage_connection_string_1 = optional(string)

      # The subscription holding the persistence storage account, when
      # not the cache's own.
      storage_account_subscription_id = optional(string)
    }))

    # The cache's managed identity: type arrives as the spec enum's name
    # string (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED)
    # with the user-assigned identity ARM IDs.
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Weekly maintenance windows. day_of_week arrives as the spec enum's
    # name string (MONDAY..SUNDAY).
    patch_schedules = optional(list(object({
      day_of_week        = string
      start_hour_utc     = optional(number, 0)
      maintenance_window = optional(string, "PT5H")
    })), [])

    # Public-endpoint IP allow-list, one Azure sub-resource each.
    firewall_rules = optional(list(object({
      name     = string
      start_ip = string
      end_ip   = string
    })), [])

    # Tenant-level platform settings passed through as raw key/values.
    tenant_settings = optional(map(string), {})

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}

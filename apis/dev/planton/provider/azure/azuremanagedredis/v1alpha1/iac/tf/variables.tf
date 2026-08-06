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
  description = "Azure Managed Redis specification"
  type = object({
    # The Azure region the instance lives in.
    region = string

    # The resource group the instance lives in. References are resolved
    # to a literal name by the platform before the module runs.
    resource_group = string

    # The instance's name: 3-63 letters/digits/hyphens; it becomes the
    # endpoint {cluster_name}.{region}.redis.azure.net.
    cluster_name = string

    # The tier family and size in one value, as the spec enum's name
    # string (e.g. BALANCED_B0, MEMORY_OPTIMIZED_M10). Required.
    sku_name = string

    # Whether the instance runs with a replica for HA and the
    # zone-redundant SLA. Fixed at creation.
    high_availability_enabled = optional(bool, true)

    # Customer-managed-key encryption. The key id is the VERSIONED Key
    # Vault data-plane id; the identity must also be attached through
    # the identity block (an ARM pairing enforced at apply time).
    # References are resolved to literals before the module runs.
    customer_managed_key = optional(object({
      key_vault_key_id          = string
      user_assigned_identity_id = string
    }))

    # The instance's managed identity: type arrives as the spec enum's
    # name string (SYSTEM_ASSIGNED / USER_ASSIGNED /
    # SYSTEM_AND_USER_ASSIGNED) with the user-assigned identity ARM IDs.
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Whether the instance answers on its public endpoint. false forces
    # all traffic through Private Link (AzurePrivateEndpoint) -- Managed
    # Redis has no VNet injection or IP firewall.
    public_network_access_enabled = optional(bool, true)

    # The Redis process itself. Required -- Azure rejects creating an
    # instance without its database.
    default_database = object({
      # Managed Redis is keyless-first: access keys are OFF by default
      # and clients authenticate with Entra tokens under access-policy
      # assignments.
      access_keys_authentication_enabled = optional(bool, false)

      # TLS posture, as the spec enum's name string (ENCRYPTED /
      # PLAINTEXT). Absent means ENCRYPTED.
      client_protocol = optional(string)

      # Shard distribution, as the spec enum's name string
      # (ENTERPRISE_CLUSTER / OSS_CLUSTER / NO_CLUSTER). Absent means
      # OSS_CLUSTER. Changing it recreates the DATABASE.
      clustering_policy = optional(string)

      # Eviction policy, as the spec enum's name string (e.g.
      # VOLATILE_LRU, NO_EVICTION). Absent means VOLATILE_LRU.
      eviction_policy = optional(string)

      # Joining a named ACTIVE geo-replication group; membership is
      # linked by AzureManagedRedisGeoReplication. Changing it recreates
      # the DATABASE.
      geo_replication_group_name = optional(string)

      # Redis modules to enable (RediSearch / RedisJSON / RedisBloom /
      # RedisTimeSeries), up to 4. Changing the set recreates the
      # DATABASE.
      modules = optional(list(object({
        name = string
        args = optional(string)
      })), [])

      # Setting a frequency ENABLES the matching persistence method.
      # AOF and RDB are mutually exclusive, and both conflict with
      # geo-replication (spec-enforced).
      persistence_append_only_file_backup_frequency = optional(string)
      persistence_redis_database_backup_frequency   = optional(string)
    })

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}

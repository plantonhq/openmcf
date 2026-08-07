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
  description = "Azure Cosmos DB account specification"
  type = object({
    # The Azure region the account is homed in (the write region should
    # match the failover_priority-0 geo_location).
    region = string

    # The resource group name (resolved from the AzureResourceGroup
    # reference before the module runs).
    resource_group = string

    # The globally unique account name (becomes the DNS endpoint).
    account_name = string

    # The API the account speaks, as the proto enum value name:
    # GLOBAL_DOCUMENT_DB (default when empty) or MONGO_DB.
    kind = optional(string, "")

    # The default consistency policy (required by Azure at creation).
    consistency_policy = object({
      consistency_level       = optional(string, "")
      max_interval_in_seconds = optional(number)
      max_staleness_prefix    = optional(number)
    })

    # The replicated regions; exactly one carries failover_priority 0.
    geo_locations = list(object({
      location          = string
      failover_priority = optional(number, 0)
      zone_redundant    = optional(bool, false)
    }))

    # Account capabilities, as proto enum value names (e.g.
    # ENABLE_SERVERLESS, ENABLE_MONGO).
    capabilities = optional(list(string), [])

    free_tier_enabled                = optional(bool, false)
    automatic_failover_enabled       = optional(bool, false)
    multiple_write_locations_enabled = optional(bool, false)
    public_network_access_enabled    = optional(bool, true)

    is_virtual_network_filter_enabled = optional(bool, false)

    # Subnets allowed through the virtual-network filter (subnet_id
    # resolved from AzureSubnet references before the module runs).
    virtual_network_rules = optional(list(object({
      subnet_id                            = string
      ignore_missing_vnet_service_endpoint = optional(bool, false)
    })), [])

    # IPv4 addresses / CIDR ranges allowed by the IP firewall.
    ip_range_filter = optional(list(string), [])

    # Backup configuration; empty means Azure's periodic default.
    backup = optional(object({
      type                = string
      tier                = optional(string, "")
      interval_in_minutes = optional(number)
      retention_in_hours  = optional(number)
      storage_redundancy  = optional(string, "")
    }))

    # The MongoDB wire-protocol version enum name (MONGO_3_2 .. MONGO_7_0).
    mongo_server_version = optional(string, "")

    # The account's managed identity.
    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))

    # The default identity the account acts as against other services.
    default_identity = optional(object({
      type                      = string
      user_assigned_identity_id = optional(string, "")
    }))

    # The versionless Key Vault key for CMK encryption (resolved from
    # the AzureKeyVaultKey reference before the module runs).
    key_vault_key_id = optional(string, "")

    analytical_storage_enabled = optional(bool, false)

    analytical_storage = optional(object({
      schema_type = string
    }))

    capacity = optional(object({
      total_throughput_limit = number
    }))

    access_key_metadata_writes_enabled = optional(bool, true)
    local_authentication_enabled       = optional(bool, true)

    # The minimum TLS version enum name (TLS_1_0 / TLS_1_1 / TLS_1_2);
    # empty means Azure's TLS 1.2 default.
    minimal_tls_version = optional(string, "")

    network_acl_bypass_for_azure_services = optional(bool, false)
    network_acl_bypass_ids                = optional(list(string), [])

    burst_capacity_enabled  = optional(bool, false)
    partition_merge_enabled = optional(bool, false)

    # One CORS rule for browser-based data-plane access.
    cors_rule = optional(object({
      allowed_origins    = list(string)
      allowed_methods    = list(string)
      allowed_headers    = list(string)
      exposed_headers    = list(string)
      max_age_in_seconds = optional(number)
    }))

    # How the account is created (DEFAULT or RESTORE); only valid with
    # CONTINUOUS backup.
    create_mode = optional(string, "")

    # The restore source and scope (create_mode RESTORE).
    restore = optional(object({
      source_cosmosdb_account_id = string
      restore_timestamp_in_utc   = string
      databases = optional(list(object({
        name             = string
        collection_names = optional(list(string), [])
      })), [])
      gremlin_databases = optional(list(object({
        name        = string
        graph_names = optional(list(string), [])
      })), [])
      tables_to_restore = optional(list(string), [])
    }))

    # User tags merged over the platform's identity tags (user wins).
    tags = optional(map(string), {})
  })
}

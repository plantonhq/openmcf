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
  description = "Azure Storage Account specification"
  type = object({
    # The Azure region the account lives in.
    region = string

    # The resource group the account lives in. References are resolved to
    # a literal name by the platform before the module runs.
    resource_group = string

    # The account's name: 3-24 lowercase letters and digits only, globally
    # unique (it becomes the DNS prefix of every service endpoint).
    account_name = string

    # The account kind, as the spec enum's name string (STORAGE_V2,
    # BLOB_STORAGE, BLOCK_BLOB_STORAGE, FILE_STORAGE, STORAGE). Unset
    # means STORAGE_V2.
    account_kind = optional(string)

    # The performance tier (STANDARD, PREMIUM). Unset means STANDARD.
    account_tier = optional(string)

    # The replication type (LRS, ZRS, GRS, GZRS, RA_GRS, RA_GZRS). Unset
    # means LRS.
    replication_type = optional(string)

    # The default blob access tier (HOT, COOL, COLD,
    # ACCESS_TIER_PREMIUM). Unset lets Azure apply Hot on the kinds that
    # support tiers.
    access_tier = optional(string)

    # The provisioned billing model version ("V2" or unset).
    provisioned_billing_model_version = optional(string)

    # The Azure Edge Zone for edge deployments; unset for regional.
    edge_zone = optional(string)

    # Whether plaintext HTTP is rejected (Azure's default is true).
    https_traffic_only_enabled = optional(bool, true)

    # The minimum TLS version (TLS1_0, TLS1_1, TLS1_2). Unset means
    # TLS1_2.
    min_tls_version = optional(string)

    # Whether shared access keys authorize requests (Azure's default is
    # true; false forces Entra-only auth).
    shared_access_key_enabled = optional(bool, true)

    # Whether the portal defaults to Entra authorization when browsing
    # data (Azure's default is false).
    default_to_oauth_authentication = optional(bool, false)

    # Whether containers may opt into anonymous public access (Azure's
    # current default is true; this only PERMITS it per container).
    allow_nested_items_to_be_public = optional(bool, true)

    # Whether the account's public endpoints accept traffic at all.
    public_network_access_enabled = optional(bool, true)

    # Copy-scope restriction (AAD, PRIVATE_LINK). Unset leaves copy
    # unrestricted.
    allowed_copy_scope = optional(string)

    # The account-wide SAS expiration policy. expiration_action arrives
    # as the spec enum's name string (LOG, BLOCK); unset means LOG.
    sas_policy = optional(object({
      expiration_period = string
      expiration_action = optional(string)
    }))

    # Whether local (SFTP) user identities may be created (Azure's
    # default is true).
    local_user_enabled = optional(bool, true)

    # Whether the SFTP endpoint is enabled (requires is_hns_enabled).
    sftp_enabled = optional(bool, false)

    # Whether object replication may cross tenants (Azure's v4 default is
    # false).
    cross_tenant_replication_enabled = optional(bool, false)

    # Whether the account has a hierarchical namespace (Data Lake Gen2).
    # Fixed at creation.
    is_hns_enabled = optional(bool, false)

    # Whether the blob service accepts NFSv3 mounts. Fixed at creation.
    nfsv3_enabled = optional(bool, false)

    # Whether file shares may grow to 100 TiB (one-way).
    large_file_share_enabled = optional(bool, false)

    # The DNS endpoint architecture (DNS_ENDPOINT_STANDARD,
    # AZURE_DNS_ZONE). Unset means the classic shared DNS.
    dns_endpoint_type = optional(string)

    # Whether data is double-encrypted at rest. Fixed at creation.
    infrastructure_encryption_enabled = optional(bool, false)

    # Queue/table service encryption key scopes (SERVICE, ACCOUNT).
    # Unset means SERVICE. Fixed at creation.
    queue_encryption_key_type = optional(string)
    table_encryption_key_type = optional(string)

    # The account's managed identity. type arrives as the spec enum's
    # name string (SYSTEM_ASSIGNED, USER_ASSIGNED,
    # SYSTEM_AND_USER_ASSIGNED); identity_ids are resolved
    # user-assigned-identity ARM IDs.
    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))

    # Customer-managed-key encryption. key_vault_key_id is a resolved
    # Key Vault data-plane key ID; user_assigned_identity_id is a
    # resolved user-assigned-identity ARM ID (must be attached via
    # identity).
    customer_managed_key = optional(object({
      key_vault_key_id          = string
      user_assigned_identity_id = string
    }))

    # Data-plane network access control. default_action arrives as the
    # spec enum's name string (ALLOW, DENY); bypass entries as
    # AZURE_SERVICES/LOGGING/METRICS/NONE; virtual_network_subnet_ids as
    # resolved subnet ARM IDs.
    network_rules = optional(object({
      default_action             = string
      bypass                     = optional(list(string), [])
      ip_rules                   = optional(list(string), [])
      virtual_network_subnet_ids = optional(list(string), [])
      private_link_access = optional(list(object({
        endpoint_resource_id = string
        endpoint_tenant_id   = optional(string)
      })), [])
    }))

    # Blob service settings. Retention days carry Azure's default of 7
    # so an empty policy block still renders a valid window.
    blob_properties = optional(object({
      versioning_enabled            = optional(bool, false)
      change_feed_enabled           = optional(bool, false)
      change_feed_retention_in_days = optional(number)
      default_service_version       = optional(string)
      last_access_time_enabled      = optional(bool, false)
      delete_retention_policy = optional(object({
        days                     = optional(number, 7)
        permanent_delete_enabled = optional(bool, false)
      }))
      container_delete_retention_policy = optional(object({
        days = optional(number, 7)
      }))
      restore_policy = optional(object({
        days = number
      }))
      cors_rules = optional(list(object({
        allowed_origins    = list(string)
        allowed_methods    = list(string)
        allowed_headers    = list(string)
        exposed_headers    = list(string)
        max_age_in_seconds = optional(number, 0)
      })), [])
    }))

    # File service settings.
    share_properties = optional(object({
      retention_policy = optional(object({
        days = optional(number, 7)
      }))
      smb = optional(object({
        versions                        = optional(list(string), [])
        authentication_types            = optional(list(string), [])
        kerberos_ticket_encryption_type = optional(list(string), [])
        channel_encryption_type         = optional(list(string), [])
        multichannel_enabled            = optional(bool, false)
      }))
      cors_rules = optional(list(object({
        allowed_origins    = list(string)
        allowed_methods    = list(string)
        allowed_headers    = list(string)
        exposed_headers    = list(string)
        max_age_in_seconds = optional(number, 0)
      })), [])
    }))

    # Static website hosting (realized as the standalone
    # azurerm_storage_account_static_website resource).
    static_website = optional(object({
      index_document     = optional(string)
      error_404_document = optional(string)
    }))

    # Network routing preference. choice arrives as the spec enum's name
    # string (MICROSOFT_ROUTING, INTERNET_ROUTING); unset means
    # Microsoft routing.
    routing = optional(object({
      choice                      = optional(string)
      publish_internet_endpoints  = optional(bool, false)
      publish_microsoft_endpoints = optional(bool, false)
    }))

    # A custom domain (CNAME) for the blob endpoint.
    custom_domain = optional(object({
      name          = string
      use_subdomain = optional(bool, false)
    }))

    # Identity-based authentication for Azure Files. directory_type
    # arrives as the spec enum's name string (AADDS, AADKERB, AD);
    # default_share_level_permission as SHARE_PERMISSION_*.
    azure_files_authentication = optional(object({
      directory_type = string
      active_directory = optional(object({
        domain_name         = string
        domain_guid         = string
        domain_sid          = optional(string)
        storage_sid         = optional(string)
        forest_name         = optional(string)
        netbios_domain_name = optional(string)
      }))
      default_share_level_permission = optional(string)
    }))

    # Account-level immutability (WORM) policy. state arrives as the
    # spec enum's name string (DISABLED, UNLOCKED, LOCKED).
    immutability_policy = optional(object({
      state                         = string
      period_since_creation_in_days = number
      allow_protected_append_writes = optional(bool, false)
    }))

    # Blob lifecycle management rules (realized as the singleton
    # azurerm_storage_management_policy resource). Aging thresholds are
    # optional numbers -- absent means the transition is not configured.
    lifecycle_rules = optional(list(object({
      name    = string
      enabled = optional(bool, true)
      filters = object({
        blob_types   = list(string)
        prefix_match = optional(list(string), [])
        match_blob_index_tags = optional(list(object({
          name      = string
          operation = optional(string, "==")
          value     = string
        })), [])
      })
      actions = object({
        base_blob = optional(object({
          tier_to_cool_after_days_since_modification_greater_than        = optional(number)
          tier_to_cool_after_days_since_last_access_time_greater_than    = optional(number)
          tier_to_cool_after_days_since_creation_greater_than            = optional(number)
          auto_tier_to_hot_from_cool_enabled                             = optional(bool, false)
          tier_to_cold_after_days_since_modification_greater_than        = optional(number)
          tier_to_cold_after_days_since_last_access_time_greater_than    = optional(number)
          tier_to_cold_after_days_since_creation_greater_than            = optional(number)
          tier_to_archive_after_days_since_modification_greater_than     = optional(number)
          tier_to_archive_after_days_since_last_access_time_greater_than = optional(number)
          tier_to_archive_after_days_since_creation_greater_than         = optional(number)
          tier_to_archive_after_days_since_last_tier_change_greater_than = optional(number)
          delete_after_days_since_modification_greater_than              = optional(number)
          delete_after_days_since_last_access_time_greater_than          = optional(number)
          delete_after_days_since_creation_greater_than                  = optional(number)
        }))
        snapshot = optional(object({
          change_tier_to_cool_after_days_since_creation                  = optional(number)
          tier_to_cold_after_days_since_creation_greater_than            = optional(number)
          change_tier_to_archive_after_days_since_creation               = optional(number)
          tier_to_archive_after_days_since_last_tier_change_greater_than = optional(number)
          delete_after_days_since_creation_greater_than                  = optional(number)
        }))
        version = optional(object({
          change_tier_to_cool_after_days_since_creation                  = optional(number)
          tier_to_cold_after_days_since_creation_greater_than            = optional(number)
          change_tier_to_archive_after_days_since_creation               = optional(number)
          tier_to_archive_after_days_since_last_tier_change_greater_than = optional(number)
          delete_after_days_since_creation                               = optional(number)
        }))
      })
    })), [])

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}

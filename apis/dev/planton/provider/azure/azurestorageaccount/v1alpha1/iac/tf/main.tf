# The storage account. Blob containers are first-class
# AzureStorageContainer resources referencing this account's ARM id --
# the account module deliberately creates none. The lifecycle policy and
# static website are separate resources below (Azure models the policy
# as a singleton per-account document, and azurerm's inline
# static_website block is deprecated in favor of the standalone
# resource).
resource "azurerm_storage_account" "main" {
  name                = var.spec.account_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The SKU trio. Kind and tier are fixed shapes (only the legacy
  # Storage -> StorageV2 upgrade changes kind in place); replication may
  # move within its zonal/non-zonal family in place, across families by
  # replacement.
  account_kind                      = local.account_kind
  account_tier                      = local.account_tier
  account_replication_type          = local.replication_type
  access_tier                       = local.access_tier
  provisioned_billing_model_version = local.provisioned_billing_model_version
  edge_zone                         = local.edge_zone

  # Transport and authorization posture. https_traffic_only and the TLS
  # floor govern transport; shared_access_key_enabled false forces every
  # data-plane request through Entra; allow_nested_items_to_be_public
  # false makes anonymous container access unrepresentable account-wide.
  https_traffic_only_enabled      = var.spec.https_traffic_only_enabled
  min_tls_version                 = local.min_tls_version
  shared_access_key_enabled       = var.spec.shared_access_key_enabled
  default_to_oauth_authentication = var.spec.default_to_oauth_authentication
  allow_nested_items_to_be_public = var.spec.allow_nested_items_to_be_public
  public_network_access_enabled   = var.spec.public_network_access_enabled
  allowed_copy_scope              = local.allowed_copy_scope
  local_user_enabled              = var.spec.local_user_enabled
  sftp_enabled                    = var.spec.sftp_enabled

  cross_tenant_replication_enabled = var.spec.cross_tenant_replication_enabled

  # Data-lake / protocol switches -- all create-time architectural
  # choices (HNS and NFSv3 are ForceNew). large_file_share_enabled is
  # sent only when true: the flag is one-way and Computed -- premium
  # FileStorage accounts have it on inherently, so an explicit false
  # would fight Azure. False means "leave it to Azure", never "disable".
  is_hns_enabled           = var.spec.is_hns_enabled
  nfsv3_enabled            = var.spec.nfsv3_enabled
  large_file_share_enabled = var.spec.large_file_share_enabled ? true : null
  dns_endpoint_type        = local.dns_endpoint_type

  # Encryption depth: infrastructure encryption double-encrypts at rest;
  # the key-type fields move queue/table data under the account key
  # scope so the CMK below covers them too. All fixed at creation.
  infrastructure_encryption_enabled = var.spec.infrastructure_encryption_enabled
  queue_encryption_key_type         = local.queue_encryption_key_type
  table_encryption_key_type         = local.table_encryption_key_type

  # The account-wide SAS lifetime policy: violations are logged (Log) or
  # tokens rejected outright (Block).
  dynamic "sas_policy" {
    for_each = var.spec.sas_policy != null ? [var.spec.sas_policy] : []
    content {
      expiration_period = sas_policy.value.expiration_period
      expiration_action = (
        sas_policy.value.expiration_action == null || sas_policy.value.expiration_action == "" ? "Log" :
        local.sas_expiration_action_map[sas_policy.value.expiration_action]
      )
    }
  }

  # The account's managed identity. A user-assigned identity must be
  # attached here for customer_managed_key to unwrap its key.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = identity.value.identity_ids
    }
  }

  # Customer-managed-key encryption. The key ID references the Key Vault
  # data plane (versionless IDs track rotations automatically); the
  # unwrapping identity must have wrap/unwrap access on the key's vault
  # BEFORE this account is created -- compose the grant in the same
  # manifest set.
  dynamic "customer_managed_key" {
    for_each = var.spec.customer_managed_key != null ? [var.spec.customer_managed_key] : []
    content {
      key_vault_key_id          = customer_managed_key.value.key_vault_key_id
      user_assigned_identity_id = customer_managed_key.value.user_assigned_identity_id
    }
  }

  # Data-plane firewall. ARM (control-plane) operations are never
  # subject to these rules; unset bypass lets Azure default to
  # AzureServices.
  dynamic "network_rules" {
    for_each = var.spec.network_rules != null ? [var.spec.network_rules] : []
    content {
      default_action             = local.network_default_action_map[network_rules.value.default_action]
      bypass                     = local.network_bypass
      ip_rules                   = network_rules.value.ip_rules
      virtual_network_subnet_ids = network_rules.value.virtual_network_subnet_ids

      dynamic "private_link_access" {
        for_each = network_rules.value.private_link_access
        content {
          endpoint_resource_id = private_link_access.value.endpoint_resource_id
          endpoint_tenant_id   = private_link_access.value.endpoint_tenant_id
        }
      }
    }
  }

  # Blob service settings. Versioning is the foundation the restore and
  # immutability features build on; the retention blocks are Azure's
  # recycle bin for blobs and containers.
  dynamic "blob_properties" {
    for_each = var.spec.blob_properties != null ? [var.spec.blob_properties] : []
    content {
      versioning_enabled            = blob_properties.value.versioning_enabled
      change_feed_enabled           = blob_properties.value.change_feed_enabled
      change_feed_retention_in_days = blob_properties.value.change_feed_retention_in_days
      default_service_version       = blob_properties.value.default_service_version
      last_access_time_enabled      = blob_properties.value.last_access_time_enabled

      dynamic "delete_retention_policy" {
        for_each = blob_properties.value.delete_retention_policy != null ? [blob_properties.value.delete_retention_policy] : []
        content {
          days                     = delete_retention_policy.value.days
          permanent_delete_enabled = delete_retention_policy.value.permanent_delete_enabled
        }
      }

      dynamic "container_delete_retention_policy" {
        for_each = blob_properties.value.container_delete_retention_policy != null ? [blob_properties.value.container_delete_retention_policy] : []
        content {
          days = container_delete_retention_policy.value.days
        }
      }

      dynamic "restore_policy" {
        for_each = blob_properties.value.restore_policy != null ? [blob_properties.value.restore_policy] : []
        content {
          days = restore_policy.value.days
        }
      }

      dynamic "cors_rule" {
        for_each = blob_properties.value.cors_rules
        content {
          allowed_origins    = cors_rule.value.allowed_origins
          allowed_methods    = cors_rule.value.allowed_methods
          allowed_headers    = cors_rule.value.allowed_headers
          exposed_headers    = cors_rule.value.exposed_headers
          max_age_in_seconds = cors_rule.value.max_age_in_seconds
        }
      }
    }
  }

  # File service settings: the share recycle bin and the SMB protocol
  # dials (multichannel is premium-only -- enforced by spec validation
  # before the plan ever runs).
  dynamic "share_properties" {
    for_each = var.spec.share_properties != null ? [var.spec.share_properties] : []
    content {
      dynamic "retention_policy" {
        for_each = share_properties.value.retention_policy != null ? [share_properties.value.retention_policy] : []
        content {
          days = retention_policy.value.days
        }
      }

      dynamic "smb" {
        for_each = share_properties.value.smb != null ? [share_properties.value.smb] : []
        content {
          versions                        = length(smb.value.versions) > 0 ? smb.value.versions : null
          authentication_types            = length(smb.value.authentication_types) > 0 ? smb.value.authentication_types : null
          kerberos_ticket_encryption_type = length(smb.value.kerberos_ticket_encryption_type) > 0 ? smb.value.kerberos_ticket_encryption_type : null
          channel_encryption_type         = length(smb.value.channel_encryption_type) > 0 ? smb.value.channel_encryption_type : null
          multichannel_enabled            = smb.value.multichannel_enabled
        }
      }

      dynamic "cors_rule" {
        for_each = share_properties.value.cors_rules
        content {
          allowed_origins    = cors_rule.value.allowed_origins
          allowed_methods    = cors_rule.value.allowed_methods
          allowed_headers    = cors_rule.value.allowed_headers
          exposed_headers    = cors_rule.value.exposed_headers
          max_age_in_seconds = cors_rule.value.max_age_in_seconds
        }
      }
    }
  }

  # Traffic routing preference and the optional routing-specific
  # endpoint publication.
  dynamic "routing" {
    for_each = var.spec.routing != null ? [var.spec.routing] : []
    content {
      choice = (
        routing.value.choice == null || routing.value.choice == "" ? "MicrosoftRouting" :
        local.routing_choice_map[routing.value.choice]
      )
      publish_internet_endpoints  = routing.value.publish_internet_endpoints
      publish_microsoft_endpoints = routing.value.publish_microsoft_endpoints
    }
  }

  # Custom domain for the blob endpoint; use_subdomain validates
  # ownership via the asverify CNAME to avoid downtime on live domains.
  dynamic "custom_domain" {
    for_each = var.spec.custom_domain != null ? [var.spec.custom_domain] : []
    content {
      name          = custom_domain.value.name
      use_subdomain = custom_domain.value.use_subdomain
    }
  }

  # Identity-based SMB authentication for Azure Files.
  dynamic "azure_files_authentication" {
    for_each = var.spec.azure_files_authentication != null ? [var.spec.azure_files_authentication] : []
    content {
      directory_type = local.directory_type_map[azure_files_authentication.value.directory_type]
      default_share_level_permission = (
        azure_files_authentication.value.default_share_level_permission == null || azure_files_authentication.value.default_share_level_permission == "" ? null :
        local.default_share_permission_map[azure_files_authentication.value.default_share_level_permission]
      )

      dynamic "active_directory" {
        for_each = azure_files_authentication.value.active_directory != null ? [azure_files_authentication.value.active_directory] : []
        content {
          domain_name         = active_directory.value.domain_name
          domain_guid         = active_directory.value.domain_guid
          domain_sid          = active_directory.value.domain_sid
          storage_sid         = active_directory.value.storage_sid
          forest_name         = active_directory.value.forest_name
          netbios_domain_name = active_directory.value.netbios_domain_name
        }
      }
    }
  }

  # Account-level WORM policy. LOCKED is irreversible -- Azure itself
  # cannot shorten a locked retention window; the block itself is
  # ForceNew.
  dynamic "immutability_policy" {
    for_each = var.spec.immutability_policy != null ? [var.spec.immutability_policy] : []
    content {
      state                         = local.immutability_state_map[immutability_policy.value.state]
      period_since_creation_in_days = immutability_policy.value.period_since_creation_in_days
      allow_protected_append_writes = immutability_policy.value.allow_protected_append_writes
    }
  }

  tags = local.final_tags
}

# Blob lifecycle management: ARM models this as ONE policy document per
# account (name hardcoded "default" server-side), which is why the rules
# fold into the account spec instead of being their own kind. Absent
# aging thresholds are simply not rendered -- azurerm's -1 sentinel is
# an HCL-ergonomics artifact this module never needs.
resource "azurerm_storage_management_policy" "main" {
  count = length(var.spec.lifecycle_rules) > 0 ? 1 : 0

  storage_account_id = azurerm_storage_account.main.id

  dynamic "rule" {
    for_each = var.spec.lifecycle_rules
    content {
      name    = rule.value.name
      enabled = rule.value.enabled

      filters {
        blob_types   = [for t in rule.value.filters.blob_types : local.lifecycle_blob_type_map[t]]
        prefix_match = rule.value.filters.prefix_match

        dynamic "match_blob_index_tag" {
          for_each = rule.value.filters.match_blob_index_tags
          content {
            name      = match_blob_index_tag.value.name
            operation = match_blob_index_tag.value.operation
            value     = match_blob_index_tag.value.value
          }
        }
      }

      actions {
        dynamic "base_blob" {
          for_each = rule.value.actions.base_blob != null ? [rule.value.actions.base_blob] : []
          content {
            tier_to_cool_after_days_since_modification_greater_than        = base_blob.value.tier_to_cool_after_days_since_modification_greater_than
            tier_to_cool_after_days_since_last_access_time_greater_than    = base_blob.value.tier_to_cool_after_days_since_last_access_time_greater_than
            tier_to_cool_after_days_since_creation_greater_than            = base_blob.value.tier_to_cool_after_days_since_creation_greater_than
            auto_tier_to_hot_from_cool_enabled                             = base_blob.value.auto_tier_to_hot_from_cool_enabled
            tier_to_cold_after_days_since_modification_greater_than        = base_blob.value.tier_to_cold_after_days_since_modification_greater_than
            tier_to_cold_after_days_since_last_access_time_greater_than    = base_blob.value.tier_to_cold_after_days_since_last_access_time_greater_than
            tier_to_cold_after_days_since_creation_greater_than            = base_blob.value.tier_to_cold_after_days_since_creation_greater_than
            tier_to_archive_after_days_since_modification_greater_than     = base_blob.value.tier_to_archive_after_days_since_modification_greater_than
            tier_to_archive_after_days_since_last_access_time_greater_than = base_blob.value.tier_to_archive_after_days_since_last_access_time_greater_than
            tier_to_archive_after_days_since_creation_greater_than         = base_blob.value.tier_to_archive_after_days_since_creation_greater_than
            tier_to_archive_after_days_since_last_tier_change_greater_than = base_blob.value.tier_to_archive_after_days_since_last_tier_change_greater_than
            delete_after_days_since_modification_greater_than              = base_blob.value.delete_after_days_since_modification_greater_than
            delete_after_days_since_last_access_time_greater_than          = base_blob.value.delete_after_days_since_last_access_time_greater_than
            delete_after_days_since_creation_greater_than                  = base_blob.value.delete_after_days_since_creation_greater_than
          }
        }

        dynamic "snapshot" {
          for_each = rule.value.actions.snapshot != null ? [rule.value.actions.snapshot] : []
          content {
            change_tier_to_cool_after_days_since_creation                  = snapshot.value.change_tier_to_cool_after_days_since_creation
            tier_to_cold_after_days_since_creation_greater_than            = snapshot.value.tier_to_cold_after_days_since_creation_greater_than
            change_tier_to_archive_after_days_since_creation               = snapshot.value.change_tier_to_archive_after_days_since_creation
            tier_to_archive_after_days_since_last_tier_change_greater_than = snapshot.value.tier_to_archive_after_days_since_last_tier_change_greater_than
            delete_after_days_since_creation_greater_than                  = snapshot.value.delete_after_days_since_creation_greater_than
          }
        }

        dynamic "version" {
          for_each = rule.value.actions.version != null ? [rule.value.actions.version] : []
          content {
            change_tier_to_cool_after_days_since_creation                  = version.value.change_tier_to_cool_after_days_since_creation
            tier_to_cold_after_days_since_creation_greater_than            = version.value.tier_to_cold_after_days_since_creation_greater_than
            change_tier_to_archive_after_days_since_creation               = version.value.change_tier_to_archive_after_days_since_creation
            tier_to_archive_after_days_since_last_tier_change_greater_than = version.value.tier_to_archive_after_days_since_last_tier_change_greater_than
            delete_after_days_since_creation                               = version.value.delete_after_days_since_creation
          }
        }
      }
    }
  }
}

# Static website hosting, via the standalone resource (azurerm's inline
# static_website block is deprecated for removal in v5). The service
# auto-creates the $web container; upload site content there.
resource "azurerm_storage_account_static_website" "main" {
  count = var.spec.static_website != null ? 1 : 0

  storage_account_id = azurerm_storage_account.main.id
  index_document     = local.static_website_index_document
  error_404_document = local.static_website_error_404_document
}

# The Cosmos DB account -- the globally distributed database account that
# owns regions, consistency, network posture, encryption, and backup.
# Databases and containers are their own kinds (AzureCosmosdbSqlDatabase /
# AzureCosmosdbSqlContainer / AzureCosmosdbMongoDatabase /
# AzureCosmosdbMongoCollection) referencing this account, so this module
# creates exactly one resource.
resource "azurerm_cosmosdb_account" "main" {
  name                = var.spec.account_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Azure's only offer type. Not modeled in the spec because there is
  # nothing to choose.
  offer_type = "Standard"

  # The API the account speaks -- fixed at creation because the wire
  # protocol shapes how every byte is stored.
  kind = local.kind

  # The default consistency for every read in the account. The staleness
  # dials are meaningful only on BoundedStaleness (azurerm suppresses
  # them elsewhere); multi-region BoundedStaleness floors (>= 100000 /
  # >= 300) are enforced by the spec before the plan ever runs.
  consistency_policy {
    consistency_level       = local.consistency_level
    max_interval_in_seconds = local.max_interval_in_seconds
    max_staleness_prefix    = local.max_staleness_prefix
  }

  # The replicated regions. Adding/removing regions and re-prioritizing
  # failover are in-place updates -- the priority-0 region is the write
  # region.
  dynamic "geo_location" {
    for_each = var.spec.geo_locations
    content {
      location          = geo_location.value.location
      failover_priority = geo_location.value.failover_priority
      zone_redundant    = geo_location.value.zone_redundant
    }
  }

  # Capabilities are exactly what the spec declares -- never injected
  # silently. Most capability changes recreate the account (Azure allows
  # only a small add/remove-in-place set), which is why they live in the
  # spec, visibly, rather than being derived.
  dynamic "capabilities" {
    for_each = local.capabilities
    content {
      name = capabilities.value
    }
  }

  free_tier_enabled                = var.spec.free_tier_enabled
  automatic_failover_enabled       = var.spec.automatic_failover_enabled
  multiple_write_locations_enabled = var.spec.multiple_write_locations_enabled
  public_network_access_enabled    = var.spec.public_network_access_enabled

  # Network posture: the virtual-network filter admits the declared
  # subnets; the IP filter admits the declared addresses; the bypass
  # settings let trusted Azure services (or specific resource IDs)
  # through.
  is_virtual_network_filter_enabled = var.spec.is_virtual_network_filter_enabled

  dynamic "virtual_network_rule" {
    for_each = var.spec.virtual_network_rules
    content {
      id                                   = virtual_network_rule.value.subnet_id
      ignore_missing_vnet_service_endpoint = virtual_network_rule.value.ignore_missing_vnet_service_endpoint
    }
  }

  ip_range_filter = var.spec.ip_range_filter

  network_acl_bypass_for_azure_services = var.spec.network_acl_bypass_for_azure_services
  network_acl_bypass_ids                = var.spec.network_acl_bypass_ids

  # Key- and metadata-write posture: disabling local auth forces every
  # data-plane caller through Entra ID (the account keys stop working);
  # disabling key metadata writes restricts database/container/throughput
  # changes to ARM callers.
  #
  # PARITY-EXCEPTION: this module uses azurerm's non-deprecated
  # local_authentication_enabled input; pulumi-azure v6.38 has bridged
  # only the deprecated localAuthenticationDisabled form, so the Pulumi
  # module sends the negation. Both set the same DisableLocalAuth wire
  # property -- the created account is identical. Re-align the Pulumi
  # module when the bridge exposes localAuthenticationEnabled.
  access_key_metadata_writes_enabled = var.spec.access_key_metadata_writes_enabled
  local_authentication_enabled       = var.spec.local_authentication_enabled

  # The TLS floor for every endpoint. Unset materializes Tls12 --
  # Azure's own default since April 2023; 1.0/1.1 exist only to keep
  # legacy clients connecting during a migration.
  minimal_tls_version = local.minimal_tls_version

  # Backup: PERIODIC -> CONTINUOUS upgrades in place; the reverse
  # recreates the account. The per-mode field pairings are enforced by
  # the spec's validation rules.
  dynamic "backup" {
    for_each = var.spec.backup != null ? [var.spec.backup] : []
    content {
      type                = local.backup_type_map[backup.value.type]
      tier                = backup.value.tier != null && backup.value.tier != "" ? local.continuous_tier_map[backup.value.tier] : null
      interval_in_minutes = backup.value.interval_in_minutes
      retention_in_hours  = backup.value.retention_in_hours
      storage_redundancy  = backup.value.storage_redundancy != null && backup.value.storage_redundancy != "" ? local.backup_storage_redundancy_map[backup.value.storage_redundancy] : null
    }
  }

  mongo_server_version = local.mongo_server_version

  # The managed identity, and which identity the account acts AS by
  # default when reaching into other services (CMK unwrapping): the
  # composed "UserAssignedIdentity=<id>" form makes CMK ride an identity
  # that exists -- and holds Key Vault grants -- before the account does.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  default_identity_type = local.default_identity_type

  # CMK encryption rides the key's VERSIONLESS Key Vault identifier so
  # rotation propagates without touching the account. Fixed at creation.
  key_vault_key_id = var.spec.key_vault_key_id != "" ? var.spec.key_vault_key_id : null

  # The analytical store: enabling is in-place, DISABLING recreates the
  # account -- the spec documents "on" as effectively permanent.
  analytical_storage_enabled = var.spec.analytical_storage_enabled

  dynamic "analytical_storage" {
    for_each = var.spec.analytical_storage != null ? [var.spec.analytical_storage] : []
    content {
      schema_type = local.analytical_schema_type_map[analytical_storage.value.schema_type]
    }
  }

  # The account-wide provisioned-throughput cap (-1 = unlimited) -- the
  # guardrail against runaway RU provisioning cost.
  dynamic "capacity" {
    for_each = var.spec.capacity != null ? [var.spec.capacity] : []
    content {
      total_throughput_limit = capacity.value.total_throughput_limit
    }
  }

  burst_capacity_enabled  = var.spec.burst_capacity_enabled
  partition_merge_enabled = var.spec.partition_merge_enabled

  dynamic "cors_rule" {
    for_each = var.spec.cors_rule != null ? [var.spec.cors_rule] : []
    content {
      allowed_origins    = cors_rule.value.allowed_origins
      allowed_methods    = cors_rule.value.allowed_methods
      allowed_headers    = cors_rule.value.allowed_headers
      exposed_headers    = cors_rule.value.exposed_headers
      max_age_in_seconds = cors_rule.value.max_age_in_seconds
    }
  }

  # RESTORE creates the account FROM a continuous-backup restore point;
  # every restore field is fixed at creation (a restore happens exactly
  # once, into a new account).
  create_mode = local.create_mode

  dynamic "restore" {
    for_each = var.spec.restore != null ? [var.spec.restore] : []
    content {
      source_cosmosdb_account_id = restore.value.source_cosmosdb_account_id
      restore_timestamp_in_utc   = restore.value.restore_timestamp_in_utc

      dynamic "database" {
        for_each = restore.value.databases
        content {
          name             = database.value.name
          collection_names = length(database.value.collection_names) > 0 ? database.value.collection_names : null
        }
      }

      dynamic "gremlin_database" {
        for_each = restore.value.gremlin_databases
        content {
          name        = gremlin_database.value.name
          graph_names = length(gremlin_database.value.graph_names) > 0 ? gremlin_database.value.graph_names : null
        }
      }

      tables_to_restore = length(restore.value.tables_to_restore) > 0 ? restore.value.tables_to_restore : null
    }
  }

  tags = local.final_tags
}

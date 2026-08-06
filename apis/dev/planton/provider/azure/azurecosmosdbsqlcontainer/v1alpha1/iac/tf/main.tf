# The SQL (NoSQL) API container -- the unit of storage, indexing, and
# scale-out. Addressed by the (resource group, account, database, name)
# tuple azurerm requires, parsed from the parent database's ARM ID. No
# Azure tags: ARM does not support tags on Cosmos child resources, so
# the platform's identity tags live on the account.
resource "azurerm_cosmosdb_sql_container" "main" {
  name                = var.spec.container_name
  resource_group_name = local.resource_group_name
  account_name        = local.cosmosdb_account_name
  database_name       = local.sql_database_name

  # The partition key -- the single most consequential design decision
  # for the container. Fixed at creation; MULTI_HASH hierarchical keys
  # require version 2 (the spec enforces the pairing).
  partition_key_paths   = var.spec.partition_key_paths
  partition_key_kind    = local.partition_key_kind
  partition_key_version = var.spec.partition_key_version

  # Dedicated throughput. Sent only when set: serverless accounts
  # reject provisioned throughput, and unset means the container shares
  # the database's throughput. The spec enforces the fixed-XOR-autoscale
  # contract before the plan ever runs.
  throughput = var.spec.throughput

  dynamic "autoscale_settings" {
    for_each = var.spec.autoscale_max_throughput != null ? [var.spec.autoscale_max_throughput] : []
    content {
      max_throughput = autoscale_settings.value
    }
  }

  # Document TTL: -1 turns TTL on with per-document expiry only; a
  # positive value expires documents after their last write.
  default_ttl = var.spec.default_ttl

  # Analytical-store TTL (requires analytical storage on the account):
  # -1 keeps analytical data forever. Disabling on an existing container
  # forces a replacement -- ARM's contract, documented on the spec.
  analytical_storage_ttl = var.spec.analytical_storage_ttl

  # Unique key constraints, scoped to the logical partition. Fixed at
  # creation.
  dynamic "unique_key" {
    for_each = var.spec.unique_keys
    content {
      paths = unique_key.value.paths
    }
  }

  # The indexing policy -- the main lever for write RU cost and query
  # performance, updatable in place. When any included/excluded path is
  # declared the policy replaces Azure's index-everything default
  # wholesale, so tuned policies include "/*" and exclude exceptions.
  dynamic "indexing_policy" {
    for_each = var.spec.indexing_policy != null ? [var.spec.indexing_policy] : []
    content {
      indexing_mode = (
        indexing_policy.value.indexing_mode == null || indexing_policy.value.indexing_mode == ""
        ? "consistent"
        : local.indexing_mode_map[indexing_policy.value.indexing_mode]
      )

      dynamic "included_path" {
        for_each = indexing_policy.value.included_paths
        content {
          path = included_path.value.path
        }
      }

      dynamic "excluded_path" {
        for_each = indexing_policy.value.excluded_paths
        content {
          path = excluded_path.value.path
        }
      }

      dynamic "composite_index" {
        for_each = indexing_policy.value.composite_indexes
        content {
          dynamic "index" {
            for_each = composite_index.value.entries
            content {
              path = index.value.path
              order = (
                index.value.order == null || index.value.order == ""
                ? "Ascending"
                : local.composite_index_order_map[index.value.order]
              )
            }
          }
        }
      }

      dynamic "spatial_index" {
        for_each = indexing_policy.value.spatial_indexes
        content {
          path = spatial_index.value.path
        }
      }
    }
  }

  # Conflict resolution for multi-region-write accounts. Fixed at
  # creation; the per-mode field pairing is enforced by the spec.
  dynamic "conflict_resolution_policy" {
    for_each = var.spec.conflict_resolution_policy != null ? [var.spec.conflict_resolution_policy] : []
    content {
      mode                          = local.conflict_resolution_mode_map[conflict_resolution_policy.value.mode]
      conflict_resolution_path      = conflict_resolution_policy.value.conflict_resolution_path != "" ? conflict_resolution_policy.value.conflict_resolution_path : null
      conflict_resolution_procedure = conflict_resolution_policy.value.conflict_resolution_procedure != "" ? conflict_resolution_policy.value.conflict_resolution_procedure : null
    }
  }
}

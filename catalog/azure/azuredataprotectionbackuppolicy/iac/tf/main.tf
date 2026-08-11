# Create the Data Protection backup policy. Exactly one variant block
# is set in the spec (validated at admission); each variant creates
# its own provider resource -- ONE resource exists per deployment.
#
# EVERY policy variant is immutable after create (the provider ships
# no update path -- near-total ForceNew): changing anything replaces
# the policy.

# --- Blob storage -----------------------------------------------------------
resource "azurerm_data_protection_backup_policy_blob_storage" "main" {
  count = var.spec.blob_storage != null ? 1 : 0

  name     = var.spec.name
  vault_id = var.spec.vault_id

  # Setting the operational duration enables the operational
  # (in-account) tier; setting the vault duration enables the vault
  # tier and requires the schedule intervals (the spec's CELs mirror
  # the provider's own AtLeastOneOf/RequiredWith lattice).
  operational_default_retention_duration = var.spec.blob_storage.operational_default_retention_duration != "" ? var.spec.blob_storage.operational_default_retention_duration : null
  vault_default_retention_duration       = var.spec.blob_storage.vault_default_retention_duration != "" ? var.spec.blob_storage.vault_default_retention_duration : null
  backup_repeating_time_intervals        = length(var.spec.blob_storage.backup_repeating_time_intervals) > 0 ? var.spec.blob_storage.backup_repeating_time_intervals : null
  time_zone                              = var.spec.blob_storage.time_zone != "" ? var.spec.blob_storage.time_zone : null

  dynamic "retention_rule" {
    for_each = var.spec.blob_storage.retention_rules
    content {
      name     = retention_rule.value.name
      priority = retention_rule.value.priority

      criteria {
        absolute_criteria      = retention_rule.value.criteria.absolute_criteria != "" ? retention_rule.value.criteria.absolute_criteria : null
        days_of_month          = length(retention_rule.value.criteria.days_of_month) > 0 ? retention_rule.value.criteria.days_of_month : null
        days_of_week           = length(retention_rule.value.criteria.days_of_week) > 0 ? retention_rule.value.criteria.days_of_week : null
        months_of_year         = length(retention_rule.value.criteria.months_of_year) > 0 ? retention_rule.value.criteria.months_of_year : null
        scheduled_backup_times = length(retention_rule.value.criteria.scheduled_backup_times) > 0 ? retention_rule.value.criteria.scheduled_backup_times : null
        weeks_of_month         = length(retention_rule.value.criteria.weeks_of_month) > 0 ? retention_rule.value.criteria.weeks_of_month : null
      }

      life_cycle {
        data_store_type = retention_rule.value.life_cycle.data_store_type
        duration        = retention_rule.value.life_cycle.duration
      }
    }
  }
}

# --- Managed disks ----------------------------------------------------------
resource "azurerm_data_protection_backup_policy_disk" "main" {
  count = var.spec.disk != null ? 1 : 0

  name     = var.spec.name
  vault_id = var.spec.vault_id

  backup_repeating_time_intervals = var.spec.disk.backup_repeating_time_intervals
  default_retention_duration      = var.spec.disk.default_retention_duration
  time_zone                       = var.spec.disk.time_zone != "" ? var.spec.disk.time_zone : null

  dynamic "retention_rule" {
    for_each = var.spec.disk.retention_rules
    content {
      name     = retention_rule.value.name
      duration = retention_rule.value.duration
      priority = retention_rule.value.priority

      criteria {
        absolute_criteria = retention_rule.value.criteria.absolute_criteria != "" ? retention_rule.value.criteria.absolute_criteria : null
      }
    }
  }
}

# --- Kubernetes (AKS) clusters ----------------------------------------------
# The one variant the provider addresses by vault NAME + resource
# group instead of vault ID -- both derived from the spec's vault_id
# (ARM IDs are structured; the parse lives in locals.tf).
resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "main" {
  count = var.spec.kubernetes_cluster != null ? 1 : 0

  name                = var.spec.name
  resource_group_name = local.vault_resource_group_name
  vault_name          = local.vault_name

  backup_repeating_time_intervals = var.spec.kubernetes_cluster.backup_repeating_time_intervals
  time_zone                       = var.spec.kubernetes_cluster.time_zone != "" ? var.spec.kubernetes_cluster.time_zone : null

  default_retention_rule {
    dynamic "life_cycle" {
      for_each = var.spec.kubernetes_cluster.default_retention_rule.life_cycles
      content {
        data_store_type = life_cycle.value.data_store_type
        duration        = life_cycle.value.duration
      }
    }
  }

  dynamic "retention_rule" {
    for_each = var.spec.kubernetes_cluster.retention_rules
    content {
      name     = retention_rule.value.name
      priority = retention_rule.value.priority

      criteria {
        absolute_criteria      = retention_rule.value.criteria.absolute_criteria != "" ? retention_rule.value.criteria.absolute_criteria : null
        days_of_week           = length(retention_rule.value.criteria.days_of_week) > 0 ? retention_rule.value.criteria.days_of_week : null
        months_of_year         = length(retention_rule.value.criteria.months_of_year) > 0 ? retention_rule.value.criteria.months_of_year : null
        scheduled_backup_times = length(retention_rule.value.criteria.scheduled_backup_times) > 0 ? retention_rule.value.criteria.scheduled_backup_times : null
        weeks_of_month         = length(retention_rule.value.criteria.weeks_of_month) > 0 ? retention_rule.value.criteria.weeks_of_month : null
      }

      dynamic "life_cycle" {
        for_each = retention_rule.value.life_cycles
        content {
          data_store_type = life_cycle.value.data_store_type
          duration        = life_cycle.value.duration
        }
      }
    }
  }
}

# --- MySQL flexible servers -------------------------------------------------
resource "azurerm_data_protection_backup_policy_mysql_flexible_server" "main" {
  count = var.spec.mysql_flexible_server != null ? 1 : 0

  name     = var.spec.name
  vault_id = var.spec.vault_id

  backup_repeating_time_intervals = var.spec.mysql_flexible_server.backup_repeating_time_intervals
  time_zone                       = var.spec.mysql_flexible_server.time_zone != "" ? var.spec.mysql_flexible_server.time_zone : null

  default_retention_rule {
    dynamic "life_cycle" {
      for_each = var.spec.mysql_flexible_server.default_retention_rule.life_cycles
      content {
        data_store_type = life_cycle.value.data_store_type
        duration        = life_cycle.value.duration
      }
    }
  }

  dynamic "retention_rule" {
    for_each = var.spec.mysql_flexible_server.retention_rules
    content {
      name     = retention_rule.value.name
      priority = retention_rule.value.priority

      criteria {
        absolute_criteria      = retention_rule.value.criteria.absolute_criteria != "" ? retention_rule.value.criteria.absolute_criteria : null
        days_of_week           = length(retention_rule.value.criteria.days_of_week) > 0 ? retention_rule.value.criteria.days_of_week : null
        months_of_year         = length(retention_rule.value.criteria.months_of_year) > 0 ? retention_rule.value.criteria.months_of_year : null
        scheduled_backup_times = length(retention_rule.value.criteria.scheduled_backup_times) > 0 ? retention_rule.value.criteria.scheduled_backup_times : null
        weeks_of_month         = length(retention_rule.value.criteria.weeks_of_month) > 0 ? retention_rule.value.criteria.weeks_of_month : null
      }

      dynamic "life_cycle" {
        for_each = retention_rule.value.life_cycles
        content {
          data_store_type = life_cycle.value.data_store_type
          duration        = life_cycle.value.duration
        }
      }
    }
  }
}

# --- PostgreSQL flexible servers --------------------------------------------
resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "main" {
  count = var.spec.postgresql_flexible_server != null ? 1 : 0

  name     = var.spec.name
  vault_id = var.spec.vault_id

  backup_repeating_time_intervals = var.spec.postgresql_flexible_server.backup_repeating_time_intervals
  time_zone                       = var.spec.postgresql_flexible_server.time_zone != "" ? var.spec.postgresql_flexible_server.time_zone : null

  default_retention_rule {
    dynamic "life_cycle" {
      for_each = var.spec.postgresql_flexible_server.default_retention_rule.life_cycles
      content {
        data_store_type = life_cycle.value.data_store_type
        duration        = life_cycle.value.duration
      }
    }
  }

  dynamic "retention_rule" {
    for_each = var.spec.postgresql_flexible_server.retention_rules
    content {
      name     = retention_rule.value.name
      priority = retention_rule.value.priority

      criteria {
        absolute_criteria      = retention_rule.value.criteria.absolute_criteria != "" ? retention_rule.value.criteria.absolute_criteria : null
        days_of_week           = length(retention_rule.value.criteria.days_of_week) > 0 ? retention_rule.value.criteria.days_of_week : null
        months_of_year         = length(retention_rule.value.criteria.months_of_year) > 0 ? retention_rule.value.criteria.months_of_year : null
        scheduled_backup_times = length(retention_rule.value.criteria.scheduled_backup_times) > 0 ? retention_rule.value.criteria.scheduled_backup_times : null
        weeks_of_month         = length(retention_rule.value.criteria.weeks_of_month) > 0 ? retention_rule.value.criteria.weeks_of_month : null
      }

      dynamic "life_cycle" {
        for_each = retention_rule.value.life_cycles
        content {
          data_store_type = life_cycle.value.data_store_type
          duration        = life_cycle.value.duration
        }
      }
    }
  }
}

# --- Data Lake storage ------------------------------------------------------
# Retention rules are FLAT here (criteria fields on the rule -- the
# provider's own shape) and priorities are assigned by ORDER: the
# provider stamps rule N with priority N+1 (there is no priority
# argument).
resource "azurerm_data_protection_backup_policy_data_lake_storage" "main" {
  count = var.spec.data_lake_storage != null ? 1 : 0

  name                            = var.spec.name
  data_protection_backup_vault_id = var.spec.vault_id

  backup_schedule            = var.spec.data_lake_storage.backup_schedule
  default_retention_duration = var.spec.data_lake_storage.default_retention_duration
  time_zone                  = var.spec.data_lake_storage.time_zone != "" ? var.spec.data_lake_storage.time_zone : null

  dynamic "retention_rule" {
    for_each = var.spec.data_lake_storage.retention_rules
    content {
      name              = retention_rule.value.name
      duration          = retention_rule.value.duration
      absolute_criteria = retention_rule.value.absolute_criteria != "" ? retention_rule.value.absolute_criteria : null

      days_of_week           = length(retention_rule.value.days_of_week) > 0 ? retention_rule.value.days_of_week : null
      months_of_year         = length(retention_rule.value.months_of_year) > 0 ? retention_rule.value.months_of_year : null
      scheduled_backup_times = length(retention_rule.value.scheduled_backup_times) > 0 ? retention_rule.value.scheduled_backup_times : null
      weeks_of_month         = length(retention_rule.value.weeks_of_month) > 0 ? retention_rule.value.weeks_of_month : null
    }
  }
}

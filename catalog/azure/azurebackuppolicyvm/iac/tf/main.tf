# Create the VM backup policy -- the schedule and retention rules that
# govern VM backups, as an ARM child of its vault
# (.../vaults/{vault}/backupPolicies/{name}). The policy is a free
# configuration object; ARM carries no tags on backup policies.
#
# The spec's CEL contracts mirror the provider's frequency/retention
# coupling, so a manifest that validates renders a schedule ARM
# accepts.
resource "azurerm_backup_policy_vm" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  recovery_vault_name = var.spec.recovery_vault_name

  # The platform default fills policy_type ("V1", the provider's own
  # default); the null guard keeps direct module invocations safe.
  # ForceNew on the provider -- changing the generation replaces the
  # policy.
  policy_type = var.spec.policy_type != null ? var.spec.policy_type : "V1"

  # The platform default fills timezone ("UTC").
  timezone = var.spec.timezone != null ? var.spec.timezone : "UTC"

  backup {
    frequency = var.spec.backup.frequency
    time      = var.spec.backup.time
    # Weekly-only (spec CEL pins it).
    weekdays = length(var.spec.backup.weekdays) > 0 ? var.spec.backup.weekdays : null
    # Hourly-only dials (spec CEL pins them to Hourly + V2).
    hour_interval = var.spec.backup.hour_interval
    hour_duration = var.spec.backup.hour_duration
  }

  # Omit when unset so the SERVICE default applies (2 days on V1, 7 on
  # V2 -- version-dependent, so the platform pins no default).
  instant_restore_retention_days = var.spec.instant_restore_retention_days

  dynamic "instant_restore_resource_group" {
    for_each = var.spec.instant_restore_resource_group != null ? [var.spec.instant_restore_resource_group] : []
    content {
      prefix = instant_restore_resource_group.value.prefix
      suffix = instant_restore_resource_group.value.suffix != "" ? instant_restore_resource_group.value.suffix : null
    }
  }

  dynamic "tiering_policy" {
    for_each = local.has_tiering_policy ? [var.spec.tiering_policy] : []
    content {
      archived_restore_point {
        mode = tiering_policy.value.archived_restore_point.mode
        # TierAfter-only age (spec CEL pairs them with the mode).
        duration      = tiering_policy.value.archived_restore_point.duration
        duration_type = tiering_policy.value.archived_restore_point.duration_type != "" ? tiering_policy.value.archived_restore_point.duration_type : null
      }
    }
  }

  # V2 only; omitted otherwise so the service keeps its consistency
  # default (application/file-system consistent when possible).
  consistency_type = var.spec.consistency_type != "" ? var.spec.consistency_type : null

  dynamic "retention_daily" {
    for_each = var.spec.retention_daily != null ? [var.spec.retention_daily] : []
    content {
      count = retention_daily.value.count
    }
  }

  dynamic "retention_weekly" {
    for_each = var.spec.retention_weekly != null ? [var.spec.retention_weekly] : []
    content {
      count    = retention_weekly.value.count
      weekdays = retention_weekly.value.weekdays
    }
  }

  # The two mutually-exclusive forms (spec CEL): week-of-month
  # (weeks + weekdays) or month days (days / include_last_days).
  dynamic "retention_monthly" {
    for_each = var.spec.retention_monthly != null ? [var.spec.retention_monthly] : []
    content {
      count             = retention_monthly.value.count
      weeks             = length(retention_monthly.value.weeks) > 0 ? retention_monthly.value.weeks : null
      weekdays          = length(retention_monthly.value.weekdays) > 0 ? retention_monthly.value.weekdays : null
      days              = length(retention_monthly.value.days) > 0 ? retention_monthly.value.days : null
      include_last_days = retention_monthly.value.include_last_days ? true : null
    }
  }

  dynamic "retention_yearly" {
    for_each = var.spec.retention_yearly != null ? [var.spec.retention_yearly] : []
    content {
      count             = retention_yearly.value.count
      months            = retention_yearly.value.months
      weeks             = length(retention_yearly.value.weeks) > 0 ? retention_yearly.value.weeks : null
      weekdays          = length(retention_yearly.value.weekdays) > 0 ? retention_yearly.value.weekdays : null
      days              = length(retention_yearly.value.days) > 0 ? retention_yearly.value.days : null
      include_last_days = retention_yearly.value.include_last_days ? true : null
    }
  }
}

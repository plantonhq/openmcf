# Create the file-share backup policy -- the schedule and retention
# rules that govern Azure Files share backups, as an ARM child of its
# vault (.../vaults/{vault}/backupPolicies/{name}). The policy is a
# free configuration object; ARM carries no tags on backup policies.
#
# The spec's CEL contracts mirror the provider's schedule-shape and
# vault-standard contracts, so a manifest that validates renders a
# schedule ARM accepts.
resource "azurerm_backup_policy_file_share" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  recovery_vault_name = var.spec.recovery_vault_name

  # The platform default fills timezone ("UTC").
  timezone = var.spec.timezone != null ? var.spec.timezone : "UTC"

  backup {
    frequency = var.spec.backup.frequency
    # Daily-only (spec CEL pins the shape to the frequency): the time
    # of day the backup runs. Empty renders null so the provider's
    # time/hourly exactly-one contract holds.
    time = var.spec.backup.time != "" ? var.spec.backup.time : null

    # Hourly-only (spec CEL): the backup window.
    dynamic "hourly" {
      for_each = var.spec.backup.hourly != null ? [var.spec.backup.hourly] : []
      content {
        interval        = hourly.value.interval
        start_time      = hourly.value.start_time
        window_duration = hourly.value.window_duration
      }
    }
  }

  # The platform default fills backup_tier ("snapshot", the provider's
  # own default); the null guard keeps direct module invocations safe.
  backup_tier = var.spec.backup_tier != null ? var.spec.backup_tier : "snapshot"

  # vault-standard only (spec CEL); omitted otherwise so the service
  # manages local snapshot retention.
  snapshot_retention_in_days = var.spec.snapshot_retention_in_days

  # ALWAYS required (the provider's own contract) -- the base
  # retention layer for both Daily and Hourly schedules.
  retention_daily {
    count = var.spec.retention_daily.count
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

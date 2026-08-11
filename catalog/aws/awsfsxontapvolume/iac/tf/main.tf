# ---------------------------------------------------------------------------
# AWS FSx for NetApp ONTAP volume
# ---------------------------------------------------------------------------
# One aws_fsx_ontap_volume resource carries the whole spec. ForceNew
# attributes (replace the volume when changed): storage_virtual_machine_id,
# name, ontap_volume_type, volume_style, the SnapLock type, and the aggregate
# layout. Size, junction path, security style, snapshot policy, storage
# efficiency, tiering, and the mutable SnapLock settings update in place.
#
# The delete-time controls (skip_final_backup, final_backup_tags,
# bypass_snaplock_enterprise_retention) never trigger an UpdateVolume call —
# they are read from state at destroy time, which is why they must be applied
# BEFORE the change that deletes the volume.
# ---------------------------------------------------------------------------

resource "aws_fsx_ontap_volume" "this" {
  # Parent SVM (ForceNew — a volume cannot move between SVMs).
  storage_virtual_machine_id = var.spec.storage_virtual_machine_id

  # The ONTAP volume name (ForceNew) — distinct from metadata.name, which
  # only becomes the Name tag.
  name = var.spec.name

  # Size flows through exactly one arm (spec-enforced XOR). size_in_bytes is
  # the byte-precise arm and the only way past 2 PiB — the provider types it
  # as a string (see locals). Both update in place.
  size_in_megabytes = var.spec.size_in_megabytes
  size_in_bytes     = local.size_in_bytes

  # Namespace mount point; omitted means the volume exists but is not
  # reachable over NFS/SMB until a junction path is set. Updatable in place.
  junction_path = local.junction_path

  # Fundamental architecture (ForceNew): RW/DP and FLEXVOL/FLEXGROUP. Sent
  # explicitly — the spec defaults (RW, FLEXVOL) are materialized before the
  # module runs.
  ontap_volume_type = var.spec.ontap_volume_type
  volume_style      = var.spec.volume_style

  # Per-volume security style; omitted inherits the SVM's root volume style.
  security_style = local.security_style

  # ONTAP snapshot policy; omitted keeps ONTAP's "default" policy.
  snapshot_policy = local.snapshot_policy

  # Tri-state efficiency switch: true enables dedup/compression/compaction,
  # false disables, omitted keeps ONTAP's per-volume-type default. The
  # provider distinguishes false from absent, so the value passes through
  # without coercion.
  storage_efficiency_enabled = var.spec.storage_efficiency_enabled

  copy_tags_to_backups = var.spec.copy_tags_to_backups

  # Delete-time controls (see the header comment).
  skip_final_backup                    = var.spec.skip_final_backup
  final_backup_tags                    = local.final_backup_tags
  bypass_snaplock_enterprise_retention = var.spec.bypass_snaplock_enterprise_retention

  # Capacity-pool tiering. cooling_period is only meaningful for AUTO /
  # SNAPSHOT_ONLY (CEL-enforced); zero means "let ONTAP default it" (31 days
  # for AUTO, 2 for SNAPSHOT_ONLY). The block is keyed on a non-empty name
  # (CEL-required, but tfvars bypass validation): a name-less block would send
  # an empty TieringPolicy struct to the AWS API.
  dynamic "tiering_policy" {
    for_each = var.spec.tiering_policy != null && try(var.spec.tiering_policy.name, "") != "" ? [var.spec.tiering_policy] : []

    content {
      name           = tiering_policy.value.name
      cooling_period = tiering_policy.value.cooling_period > 0 ? tiering_policy.value.cooling_period : null
    }
  }

  # SnapLock WORM configuration. snaplock_type is ForceNew; everything else
  # in the block updates in place.
  dynamic "snaplock_configuration" {
    for_each = var.spec.snaplock_configuration != null ? [var.spec.snaplock_configuration] : []

    content {
      snaplock_type              = snaplock_configuration.value.snaplock_type
      audit_log_volume           = snaplock_configuration.value.audit_log_volume
      privileged_delete          = snaplock_configuration.value.privileged_delete
      volume_append_mode_enabled = snaplock_configuration.value.volume_append_mode_enabled

      # Autocommit: files idle for the period transition to WORM
      # automatically. The block is keyed on a non-empty type (the AWS API
      # marks AutocommitPeriod.Type required; the CEL enforces it, but tfvars
      # bypass validation): a type-less block would send an empty
      # AutocommitPeriod struct to the AWS API. value is CEL-required whenever
      # type is a real unit, so the zero-elision below can never drop a
      # meaningful value.
      dynamic "autocommit_period" {
        for_each = snaplock_configuration.value.autocommit_period != null && try(snaplock_configuration.value.autocommit_period.type, "") != "" ? [snaplock_configuration.value.autocommit_period] : []

        content {
          type  = autocommit_period.value.type
          value = autocommit_period.value.value > 0 ? autocommit_period.value.value : null
        }
      }

      # Retention bounds. Each duration is keyed on a non-empty type (the AWS
      # API marks RetentionPeriod.Type required), and the container renders
      # only when at least one duration survives — an empty retention_period
      # block would send an empty struct to the AWS API. A value of 0 IS
      # meaningful for unit types (e.g. a 0-day minimum retention), so the
      # value is sent whenever the type is a unit type and only elided for
      # INFINITE/UNSPECIFIED, where AWS takes no value.
      dynamic "retention_period" {
        for_each = snaplock_configuration.value.retention_period != null && (
          try(snaplock_configuration.value.retention_period.default_retention.type, "") != "" ||
          try(snaplock_configuration.value.retention_period.minimum_retention.type, "") != "" ||
          try(snaplock_configuration.value.retention_period.maximum_retention.type, "") != ""
        ) ? [snaplock_configuration.value.retention_period] : []

        content {
          dynamic "default_retention" {
            for_each = try(retention_period.value.default_retention.type, "") != "" ? [retention_period.value.default_retention] : []

            content {
              type  = default_retention.value.type
              value = contains(["INFINITE", "UNSPECIFIED"], default_retention.value.type) ? null : default_retention.value.value
            }
          }

          dynamic "minimum_retention" {
            for_each = try(retention_period.value.minimum_retention.type, "") != "" ? [retention_period.value.minimum_retention] : []

            content {
              type  = minimum_retention.value.type
              value = contains(["INFINITE", "UNSPECIFIED"], minimum_retention.value.type) ? null : minimum_retention.value.value
            }
          }

          dynamic "maximum_retention" {
            for_each = try(retention_period.value.maximum_retention.type, "") != "" ? [retention_period.value.maximum_retention] : []

            content {
              type  = maximum_retention.value.type
              value = contains(["INFINITE", "UNSPECIFIED"], maximum_retention.value.type) ? null : maximum_retention.value.value
            }
          }
        }
      }
    }
  }

  # FLEXGROUP aggregate layout (ForceNew). Omitted aggregates let AWS spread
  # constituents across all of the file system's aggregates —
  # constituents_per_aggregate WITHOUT aggregates is a legitimate shape. The
  # block renders only when at least one leaf is set: an empty block would
  # send an empty AggregateConfiguration struct to the AWS API.
  dynamic "aggregate_configuration" {
    for_each = var.spec.aggregate_configuration != null && (
      length(try(var.spec.aggregate_configuration.aggregates, [])) > 0 ||
      try(var.spec.aggregate_configuration.constituents_per_aggregate, 0) > 0
    ) ? [var.spec.aggregate_configuration] : []

    content {
      aggregates                 = length(aggregate_configuration.value.aggregates) > 0 ? aggregate_configuration.value.aggregates : null
      constituents_per_aggregate = aggregate_configuration.value.constituents_per_aggregate > 0 ? aggregate_configuration.value.constituents_per_aggregate : null
    }
  }

  tags = local.aws_tags
}

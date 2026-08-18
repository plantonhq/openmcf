# One Data Lifecycle Manager policy - AWS's simplified default mode
# XOR the full custom engine (schedules or event-based).
#
# Lifecycle facts the render below depends on:
#   - the provider's default_policy top-level argument and the
#     policy_details.policy_language are DERIVED here from the
#     configured arm (SIMPLIFIED for default mode, STANDARD for
#     custom) - a spec field for either could contradict the arm;
#   - default mode expresses "no exclusions" and "AWS-default
#     cadence" by OMITTING arguments (the provider diff-suppresses
#     create_interval 0->1 and retain_interval 0->7);
#   - schedule.copy_tags is ForceNew (replaces the whole schedule);
#     everything else in a schedule updates in place;
#   - the policy references NO volume or snapshot - it targets by
#     tags at fire time, which is why this kind carries no
#     volume/snapshot edges.

locals {
  is_default = var.spec.default_policy != null

  description = var.spec.description != "" ? var.spec.description : var.metadata.name
  state       = var.spec.disabled ? "DISABLED" : "ENABLED"

  custom_policy_type = local.is_default ? null : (
    var.spec.custom_policy.policy_type != "" ? var.spec.custom_policy.policy_type : "EBS_SNAPSHOT_MANAGEMENT"
  )
  is_event_based = local.custom_policy_type == "EVENT_BASED_POLICY"
}

resource "aws_dlm_lifecycle_policy" "this" {
  description        = local.description
  execution_role_arn = var.spec.execution_role_arn
  state              = local.state

  # DEFAULT mode: the provider wants the resource type here and the
  # dials inside policy_details.
  default_policy = local.is_default ? var.spec.default_policy.resource_type : null

  policy_details {
    # ----- DEFAULT mode -------------------------------------------------
    resource_type   = local.is_default ? var.spec.default_policy.resource_type : null
    create_interval = local.is_default && var.spec.default_policy.create_interval_days > 0 ? var.spec.default_policy.create_interval_days : null
    retain_interval = local.is_default && var.spec.default_policy.retain_interval_days > 0 ? var.spec.default_policy.retain_interval_days : null
    copy_tags       = local.is_default && var.spec.default_policy.copy_tags ? true : null
    extend_deletion = local.is_default && var.spec.default_policy.extend_deletion ? true : null

    dynamic "exclusions" {
      for_each = local.is_default && var.spec.default_policy.exclusions != null ? [var.spec.default_policy.exclusions] : []
      content {
        exclude_boot_volumes = exclusions.value.exclude_boot_volumes ? true : null
        exclude_tags         = length(exclusions.value.exclude_tags) > 0 ? exclusions.value.exclude_tags : null
        exclude_volume_types = length(exclusions.value.exclude_volume_types) > 0 ? exclusions.value.exclude_volume_types : null
      }
    }

    # ----- CUSTOM mode --------------------------------------------------
    policy_type        = local.is_default ? null : local.custom_policy_type
    resource_types     = !local.is_default && length(var.spec.custom_policy.resource_types) > 0 ? var.spec.custom_policy.resource_types : null
    resource_locations = !local.is_default && length(var.spec.custom_policy.resource_locations) > 0 ? var.spec.custom_policy.resource_locations : null
    target_tags        = !local.is_default && length(var.spec.custom_policy.target_tags) > 0 ? var.spec.custom_policy.target_tags : null

    dynamic "parameters" {
      for_each = !local.is_default && var.spec.custom_policy.parameters != null ? [var.spec.custom_policy.parameters] : []
      content {
        exclude_boot_volume = parameters.value.exclude_boot_volume ? true : null
        no_reboot           = parameters.value.no_reboot ? true : null
      }
    }

    dynamic "event_source" {
      for_each = !local.is_default && var.spec.custom_policy.event_source != null ? [var.spec.custom_policy.event_source] : []
      content {
        type = "MANAGED_CWE"
        parameters {
          description_regex = event_source.value.description_regex
          event_type        = event_source.value.event_type
          snapshot_owner    = event_source.value.snapshot_owners
        }
      }
    }

    dynamic "action" {
      for_each = !local.is_default && var.spec.custom_policy.action != null ? [var.spec.custom_policy.action] : []
      content {
        name = action.value.name
        dynamic "cross_region_copy" {
          for_each = action.value.cross_region_copies
          content {
            target = cross_region_copy.value.target
            encryption_configuration {
              encrypted = cross_region_copy.value.encrypted
              cmk_arn   = cross_region_copy.value.cmk_arn != "" ? cross_region_copy.value.cmk_arn : null
            }
            dynamic "retain_rule" {
              for_each = cross_region_copy.value.retain_rule != null ? [cross_region_copy.value.retain_rule] : []
              content {
                interval      = retain_rule.value.interval
                interval_unit = retain_rule.value.interval_unit
              }
            }
          }
        }
      }
    }

    dynamic "schedule" {
      for_each = local.is_default ? [] : var.spec.custom_policy.schedules
      content {
        name          = schedule.value.name
        copy_tags     = schedule.value.copy_tags ? true : null
        tags_to_add   = length(schedule.value.tags_to_add) > 0 ? schedule.value.tags_to_add : null
        variable_tags = length(schedule.value.variable_tags) > 0 ? schedule.value.variable_tags : null

        create_rule {
          interval        = schedule.value.create_rule.interval_hours > 0 ? schedule.value.create_rule.interval_hours : null
          interval_unit   = schedule.value.create_rule.interval_hours > 0 ? "HOURS" : null
          times           = length(schedule.value.create_rule.times) > 0 ? schedule.value.create_rule.times : null
          cron_expression = schedule.value.create_rule.cron_expression != "" ? schedule.value.create_rule.cron_expression : null
          location        = schedule.value.create_rule.location != "" ? schedule.value.create_rule.location : null

          dynamic "scripts" {
            for_each = schedule.value.create_rule.scripts != null ? [schedule.value.create_rule.scripts] : []
            content {
              execution_handler                   = scripts.value.execution_handler
              stages                              = length(scripts.value.stages) > 0 ? scripts.value.stages : null
              execution_handler_service           = scripts.value.execution_handler_service != "" ? scripts.value.execution_handler_service : null
              execute_operation_on_script_failure = scripts.value.execute_operation_on_script_failure ? true : null
              execution_timeout                   = scripts.value.execution_timeout_seconds > 0 ? scripts.value.execution_timeout_seconds : null
              maximum_retry_count                 = scripts.value.maximum_retry_count > 0 ? scripts.value.maximum_retry_count : null
            }
          }
        }

        retain_rule {
          count         = schedule.value.retain_rule.count > 0 ? schedule.value.retain_rule.count : null
          interval      = schedule.value.retain_rule.interval > 0 ? schedule.value.retain_rule.interval : null
          interval_unit = schedule.value.retain_rule.interval_unit != "" ? schedule.value.retain_rule.interval_unit : null
        }

        dynamic "archive_rule" {
          for_each = schedule.value.archive_rule != null ? [schedule.value.archive_rule] : []
          content {
            archive_retain_rule {
              retention_archive_tier {
                count         = archive_rule.value.count > 0 ? archive_rule.value.count : null
                interval      = archive_rule.value.interval > 0 ? archive_rule.value.interval : null
                interval_unit = archive_rule.value.interval_unit != "" ? archive_rule.value.interval_unit : null
              }
            }
          }
        }

        dynamic "cross_region_copy_rule" {
          for_each = schedule.value.cross_region_copy_rules
          content {
            target_region = cross_region_copy_rule.value.target_region
            encrypted     = cross_region_copy_rule.value.encrypted
            cmk_arn       = cross_region_copy_rule.value.cmk_arn != "" ? cross_region_copy_rule.value.cmk_arn : null
            copy_tags     = cross_region_copy_rule.value.copy_tags ? true : null

            dynamic "retain_rule" {
              for_each = cross_region_copy_rule.value.retain_rule != null ? [cross_region_copy_rule.value.retain_rule] : []
              content {
                interval      = retain_rule.value.interval
                interval_unit = retain_rule.value.interval_unit
              }
            }

            dynamic "deprecate_rule" {
              for_each = cross_region_copy_rule.value.deprecate_rule != null ? [cross_region_copy_rule.value.deprecate_rule] : []
              content {
                interval      = deprecate_rule.value.interval
                interval_unit = deprecate_rule.value.interval_unit
              }
            }
          }
        }

        dynamic "deprecate_rule" {
          for_each = schedule.value.deprecate_rule != null ? [schedule.value.deprecate_rule] : []
          content {
            count         = deprecate_rule.value.count > 0 ? deprecate_rule.value.count : null
            interval      = deprecate_rule.value.interval > 0 ? deprecate_rule.value.interval : null
            interval_unit = deprecate_rule.value.interval_unit != "" ? deprecate_rule.value.interval_unit : null
          }
        }

        dynamic "fast_restore_rule" {
          for_each = schedule.value.fast_restore_rule != null ? [schedule.value.fast_restore_rule] : []
          content {
            availability_zones = fast_restore_rule.value.availability_zones
            count              = fast_restore_rule.value.count > 0 ? fast_restore_rule.value.count : null
            interval           = fast_restore_rule.value.interval > 0 ? fast_restore_rule.value.interval : null
            interval_unit      = fast_restore_rule.value.interval_unit != "" ? fast_restore_rule.value.interval_unit : null
          }
        }

        dynamic "share_rule" {
          for_each = schedule.value.share_rule != null ? [schedule.value.share_rule] : []
          content {
            target_accounts       = share_rule.value.target_accounts
            unshare_interval      = share_rule.value.unshare_interval > 0 ? share_rule.value.unshare_interval : null
            unshare_interval_unit = share_rule.value.unshare_interval_unit != "" ? share_rule.value.unshare_interval_unit : null
          }
        }
      }
    }
  }

  tags = local.aws_tags
}

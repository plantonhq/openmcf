# An AWS Backup plan: scheduled rules creating recovery points, plus
# the folded resource selections that assign resources to the plan.
#
# Lifecycle facts the render below depends on:
#   - the plan's identity at AWS is a generated UUID, not the name;
#     the name forces replacement;
#   - the provider CANNOT send an explicit zero for the lifecycle day
#     counts (zero is dropped as unset) - the spec presence-types them
#     so the truth is explicit;
#   - opt_in_to_archive_for_supported_resources is transmitted only
#     when true (the provider never sends an explicit false);
#   - selections are fully ForceNew (no update path) and AWS refuses
#     to delete a plan while selections exist - the provider retries
#     the plan delete while the folded selections drain.

resource "aws_backup_plan" "this" {
  # metadata.name is the plan name on both engines. Changing it
  # replaces the plan (and re-parents every selection).
  name = var.metadata.name

  dynamic "rule" {
    for_each = { for r in var.spec.rules : r.name => r }
    content {
      rule_name         = rule.value.name
      target_vault_name = rule.value.target_vault_name

      # Rendered only on an explicit choice so the module never fights
      # the provider defaults (timezone Etc/UTC, start window 60,
      # completion window 180).
      schedule                     = rule.value.schedule != "" ? rule.value.schedule : null
      schedule_expression_timezone = rule.value.schedule_expression_timezone != "" ? rule.value.schedule_expression_timezone : null
      start_window                 = rule.value.start_window_minutes != 0 ? rule.value.start_window_minutes : null
      completion_window            = rule.value.completion_window_minutes != 0 ? rule.value.completion_window_minutes : null

      enable_continuous_backup = rule.value.enable_continuous_backup

      recovery_point_tags = length(rule.value.recovery_point_tags) > 0 ? rule.value.recovery_point_tags : null

      target_logically_air_gapped_backup_vault_arn = rule.value.target_logically_air_gapped_backup_vault_arn != "" ? rule.value.target_logically_air_gapped_backup_vault_arn : null

      dynamic "lifecycle" {
        for_each = rule.value.lifecycle != null ? [rule.value.lifecycle] : []
        content {
          cold_storage_after                        = lifecycle.value.cold_storage_after_days
          delete_after                              = lifecycle.value.delete_after_days
          opt_in_to_archive_for_supported_resources = lifecycle.value.opt_in_to_archive_for_supported_resources
        }
      }

      dynamic "copy_action" {
        for_each = rule.value.copy_actions
        content {
          destination_vault_arn = copy_action.value.destination_vault_arn

          dynamic "lifecycle" {
            for_each = copy_action.value.lifecycle != null ? [copy_action.value.lifecycle] : []
            content {
              cold_storage_after                        = lifecycle.value.cold_storage_after_days
              delete_after                              = lifecycle.value.delete_after_days
              opt_in_to_archive_for_supported_resources = lifecycle.value.opt_in_to_archive_for_supported_resources
            }
          }
        }
      }

      dynamic "scan_action" {
        for_each = rule.value.scan_actions
        content {
          malware_scanner = scan_action.value.malware_scanner
          scan_mode       = scan_action.value.scan_mode
        }
      }
    }
  }

  dynamic "advanced_backup_setting" {
    for_each = var.spec.advanced_backup_settings
    content {
      resource_type  = advanced_backup_setting.value.resource_type
      backup_options = advanced_backup_setting.value.backup_options
    }
  }

  dynamic "scan_setting" {
    for_each = var.spec.scan_setting != null ? [var.spec.scan_setting] : []
    content {
      malware_scanner  = scan_setting.value.malware_scanner
      resource_types   = scan_setting.value.resource_types
      scanner_role_arn = scan_setting.value.scanner_role_arn
    }
  }

  tags = local.aws_tags
}

# Folded selections, keyed by name (fully ForceNew at the provider -
# any change replaces the selection, never edits it).
resource "aws_backup_selection" "this" {
  for_each = { for s in var.spec.selections : s.name => s }

  name         = each.value.name
  plan_id      = aws_backup_plan.this.id
  iam_role_arn = each.value.iam_role_arn

  resources     = length(each.value.resources) > 0 ? each.value.resources : null
  not_resources = length(each.value.not_resources) > 0 ? each.value.not_resources : null

  dynamic "selection_tag" {
    for_each = each.value.selection_tags
    content {
      type  = selection_tag.value.type
      key   = selection_tag.value.key
      value = selection_tag.value.value
    }
  }

  dynamic "condition" {
    for_each = each.value.conditions != null ? [each.value.conditions] : []
    content {
      dynamic "string_equals" {
        for_each = condition.value.string_equals
        content {
          key   = string_equals.value.key
          value = string_equals.value.value
        }
      }
      dynamic "string_not_equals" {
        for_each = condition.value.string_not_equals
        content {
          key   = string_not_equals.value.key
          value = string_not_equals.value.value
        }
      }
      dynamic "string_like" {
        for_each = condition.value.string_like
        content {
          key   = string_like.value.key
          value = string_like.value.value
        }
      }
      dynamic "string_not_like" {
        for_each = condition.value.string_not_like
        content {
          key   = string_not_like.value.key
          value = string_not_like.value.value
        }
      }
    }
  }
}

# An AWS Backup restore testing plan with its folded per-type
# selections: scheduled automated restore tests proving recovery
# points actually restore.
#
# Lifecycle facts the render below depends on:
#   - AWS restore testing names forbid hyphens and periods, so the
#     names are spec.plan_name / selection.name (explicit fields),
#     never metadata.name;
#   - a selection covers resources by EXACTLY ONE of explicit ARNs or
#     tag conditions (the provider enforces the same rule
#     resource-wide; the spec's CEL matches it);
#   - AWS returns empty condition lists as present-but-empty; the
#     provider collapses that to absent on both create and read - the
#     module renders the block only when the spec carries conditions;
#   - several Optional+Computed knobs (timezone, start window,
#     selection window, validation window, exclude vaults) keep an
#     AWS-side value once set - they cannot be cleared back to unset.

resource "aws_backup_restore_testing_plan" "this" {
  name = var.spec.plan_name

  schedule_expression          = var.spec.schedule_expression
  schedule_expression_timezone = var.spec.schedule_expression_timezone != "" ? var.spec.schedule_expression_timezone : null
  start_window_hours           = var.spec.start_window_hours != 0 ? var.spec.start_window_hours : null

  recovery_point_selection {
    algorithm             = var.spec.recovery_point_selection.algorithm
    include_vaults        = var.spec.recovery_point_selection.include_vaults
    recovery_point_types  = var.spec.recovery_point_selection.recovery_point_types
    exclude_vaults        = length(var.spec.recovery_point_selection.exclude_vaults) > 0 ? var.spec.recovery_point_selection.exclude_vaults : null
    selection_window_days = var.spec.recovery_point_selection.selection_window_days != 0 ? var.spec.recovery_point_selection.selection_window_days : null
  }

  tags = local.aws_tags
}

# Folded per-type selections, keyed by name.
resource "aws_backup_restore_testing_selection" "this" {
  for_each = { for s in var.spec.selections : s.name => s }

  name                      = each.value.name
  restore_testing_plan_name = aws_backup_restore_testing_plan.this.name
  protected_resource_type   = each.value.protected_resource_type
  iam_role_arn              = each.value.iam_role_arn

  protected_resource_arns = length(each.value.protected_resource_arns) > 0 ? each.value.protected_resource_arns : null

  dynamic "protected_resource_conditions" {
    for_each = each.value.protected_resource_conditions != null ? [each.value.protected_resource_conditions] : []
    content {
      dynamic "string_equals" {
        for_each = protected_resource_conditions.value.string_equals
        content {
          key   = string_equals.value.key
          value = string_equals.value.value
        }
      }
      dynamic "string_not_equals" {
        for_each = protected_resource_conditions.value.string_not_equals
        content {
          key   = string_not_equals.value.key
          value = string_not_equals.value.value
        }
      }
    }
  }

  restore_metadata_overrides = length(each.value.restore_metadata_overrides) > 0 ? each.value.restore_metadata_overrides : null
  validation_window_hours    = each.value.validation_window_hours != 0 ? each.value.validation_window_hours : null
}

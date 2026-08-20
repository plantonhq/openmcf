# CloudWatch Synthetics: one optional canary, owned groups, and the
# canary's group joins.
#
# Lifecycle facts the render below depends on:
#   - the canary's name is metadata.name (the canary charset is
#     lowercase letters, digits, hyphens, underscores) and renaming
#     replaces the canary;
#   - start_canary drives Start/StopCanary around updates - false
#     creates the canary READY but never running (no run costs);
#   - a canary that lands in CREATE_FAILED is deleted and recreated by
#     the provider (AWS offers no other repair);
#   - run_config.environment_variables are WRITE-ONLY at AWS (reads
#     never return them) - never put secrets there;
#   - the association joins the canary by ARN and the group by NAME
#     (create/delete only - both sides replace on change);
#   - groups and the canary are tagged; the association join is
#     untaggable at AWS.

resource "aws_synthetics_canary" "this" {
  count = var.spec.canary != null ? 1 : 0

  name                 = var.metadata.name
  artifact_s3_location = local.artifact_s3_location
  execution_role_arn   = var.spec.canary.execution_role_arn
  handler              = var.spec.canary.handler
  runtime_version      = var.spec.canary.runtime_version

  s3_bucket  = var.spec.canary.code.s3_bucket
  s3_key     = var.spec.canary.code.s3_key
  s3_version = var.spec.canary.code.s3_version != "" ? var.spec.canary.code.s3_version : null

  schedule {
    expression          = var.spec.canary.schedule.expression
    duration_in_seconds = var.spec.canary.schedule.duration_in_seconds != 0 ? var.spec.canary.schedule.duration_in_seconds : null

    dynamic "retry_config" {
      for_each = var.spec.canary.schedule.max_retries != null ? [var.spec.canary.schedule.max_retries] : []
      content {
        max_retries = retry_config.value
      }
    }
  }

  dynamic "run_config" {
    for_each = var.spec.canary.run_config != null ? [var.spec.canary.run_config] : []
    content {
      active_tracing        = run_config.value.active_tracing
      environment_variables = length(run_config.value.environment_variables) > 0 ? run_config.value.environment_variables : null
      ephemeral_storage     = run_config.value.ephemeral_storage != null ? run_config.value.ephemeral_storage : null
      memory_in_mb          = run_config.value.memory_in_mb != null ? run_config.value.memory_in_mb : null
      timeout_in_seconds    = run_config.value.timeout_in_seconds != null ? run_config.value.timeout_in_seconds : null
    }
  }

  dynamic "vpc_config" {
    for_each = var.spec.canary.vpc_config != null ? [var.spec.canary.vpc_config] : []
    content {
      subnet_ids                  = vpc_config.value.subnet_ids
      security_group_ids          = length(vpc_config.value.security_group_ids) > 0 ? vpc_config.value.security_group_ids : null
      ipv6_allowed_for_dual_stack = vpc_config.value.ipv6_allowed_for_dual_stack
    }
  }

  dynamic "artifact_config" {
    for_each = (var.spec.canary.artifact_encryption_mode != "" || var.spec.canary.artifact_encryption_kms_key_arn != "") ? [1] : []
    content {
      s3_encryption {
        encryption_mode = var.spec.canary.artifact_encryption_mode != "" ? var.spec.canary.artifact_encryption_mode : null
        kms_key_arn     = var.spec.canary.artifact_encryption_kms_key_arn != "" ? var.spec.canary.artifact_encryption_kms_key_arn : null
      }
    }
  }

  failure_retention_period = var.spec.canary.failure_retention_period != null ? var.spec.canary.failure_retention_period : null
  success_retention_period = var.spec.canary.success_retention_period != null ? var.spec.canary.success_retention_period : null

  start_canary  = var.spec.canary.start_canary
  delete_lambda = var.spec.canary.delete_lambda

  tags = local.aws_tags
}

resource "aws_synthetics_group" "this" {
  for_each = local.groups_by_name

  name = each.value.name

  tags = local.aws_tags
}

resource "aws_synthetics_group_association" "this" {
  for_each = local.group_names

  canary_arn = aws_synthetics_canary.this[0].arn
  group_name = each.value

  # Joins to owned groups wait for their group; external names resolve
  # as-is.
  depends_on = [aws_synthetics_group.this]
}

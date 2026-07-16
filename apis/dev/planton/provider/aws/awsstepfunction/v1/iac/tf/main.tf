# AWS Step Functions state machine.
#
# One resource carries the whole surface: the ASL definition, the execution
# role, observability (CloudWatch logging + X-Ray tracing), customer-managed
# encryption, and version publishing. The heavy lifting is in locals.tf, which
# normalizes the spec's presence semantics into provider arguments.
resource "aws_sfn_state_machine" "this" {
  name       = local.resource_name
  definition = local.definition
  role_arn   = var.spec.role_arn
  type       = local.sm_type

  # Publish an immutable version on create and on every configuration change.
  # The latest version's ARN is exported as a stack output so consumers can
  # pin executions to a snapshot instead of the mutable state machine.
  publish = var.spec.publish

  tags = local.aws_tags

  # X-Ray tracing is a single toggle; the role must be able to put trace
  # segments (xray:PutTraceSegments / PutTelemetryRecords).
  dynamic "tracing_configuration" {
    for_each = var.spec.tracing_enabled ? [1] : []
    content {
      enabled = true
    }
  }

  # Execution-history logging. Only rendered for a real level -- AWS treats
  # level OFF and an absent block identically, and rendering an OFF block
  # with no destination would fail the provider's plan-time validation.
  dynamic "logging_configuration" {
    for_each = local.logging_enabled ? [1] : []
    content {
      level                  = local.logging_level
      include_execution_data = try(var.spec.logging.include_execution_data, false)
      log_destination        = local.log_destination
    }
  }

  # Customer-managed KMS encryption. AWS's other arm (AWS_OWNED_KEY) is the
  # no-block default, so the spec models exactly one honest shape: a block
  # with a key means customer-managed.
  dynamic "encryption_configuration" {
    for_each = local.has_encryption ? [1] : []
    content {
      type                              = "CUSTOMER_MANAGED_KMS_KEY"
      kms_key_id                        = var.spec.encryption.kms_key_id
      kms_data_key_reuse_period_seconds = local.kms_data_key_reuse_period
    }
  }
}

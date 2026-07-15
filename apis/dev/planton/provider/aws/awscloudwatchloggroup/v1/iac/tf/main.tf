# The log group container. All satellite resources below key off its name and
# share its lifecycle — that is why they are folded into this module instead of
# being separate kinds.
resource "aws_cloudwatch_log_group" "this" {
  name = local.resource_name

  # Retention: null leaves the retention policy unmanaged (AWS default: never
  # expire). The provider issues a separate PutRetentionPolicy call for this.
  retention_in_days = local.retention_in_days

  # KMS encryption: customer-managed key for log data at rest. Associating or
  # clearing the key is an in-place update (AssociateKmsKey/DisassociateKmsKey).
  kms_key_id = local.kms_key_id

  # Log group class is create-time (ForceNew): STANDARD, INFREQUENT_ACCESS, or
  # DELIVERY. Null lets AWS default to STANDARD.
  log_group_class = local.log_group_class

  # Deletion protection blocks every delete (including IaC destroy) until the
  # flag is cleared; applied via a separate PutLogGroupDeletionProtection call.
  deletion_protection_enabled = local.deletion_protection_enabled

  tags = local.aws_tags
}

# Metric filters: one resource per named filter (many-per-group). Each turns
# matching log events into a CloudWatch metric — the raw material for alarms
# on signals that live only in logs.
resource "aws_cloudwatch_log_metric_filter" "this" {
  for_each = { for filter in var.spec.metric_filters : filter.name => filter }

  name           = each.value.name
  log_group_name = aws_cloudwatch_log_group.this.name

  # The provider requires pattern; empty string matches every event.
  pattern = each.value.pattern

  # Only applies when the group has an active log transformer.
  apply_on_transformed_logs = each.value.apply_on_transformed_logs ? true : null

  metric_transformation {
    name      = each.value.transformation.metric_name
    namespace = each.value.transformation.metric_namespace
    value     = each.value.transformation.metric_value

    # default_value is a tri-state: null publishes nothing for non-matching
    # periods. The provider types it as a nullable-float STRING, so format the
    # number only when present. (AWS forbids combining it with dimensions —
    # the spec CEL rejects that before plan time.)
    default_value = each.value.transformation.default_value != null ? tostring(each.value.transformation.default_value) : null

    dimensions = length(each.value.transformation.dimensions) > 0 ? each.value.transformation.dimensions : null

    # The provider defaults unit to "None" when omitted.
    unit = each.value.transformation.unit != "" ? each.value.transformation.unit : null
  }
}

# Subscription filters: real-time delivery of matching events to a Kinesis
# stream, Firehose delivery stream, or Lambda function. AWS allows at most two
# per group (CEL-enforced in the spec).
resource "aws_cloudwatch_log_subscription_filter" "this" {
  for_each = { for filter in var.spec.subscription_filters : filter.name => filter }

  name           = each.value.name
  log_group_name = aws_cloudwatch_log_group.this.name

  # References arrive pre-resolved to plain ARNs by the orchestrator.
  destination_arn = each.value.destination_arn

  # The provider requires filter_pattern; empty string delivers everything.
  filter_pattern = each.value.filter_pattern

  # Required for Kinesis/Firehose destinations (CloudWatch Logs assumes this
  # role to put records); Lambda destinations authorize via a Lambda
  # permission instead, so null is correct there.
  role_arn = each.value.role_arn != "" ? each.value.role_arn : null

  # Only meaningful for Kinesis stream destinations; the provider defaults to
  # ByLogStream (per-stream ordering preserved).
  distribution = each.value.distribution != "" ? each.value.distribution : null

  # Source-account/region enrichment for centralized log destinations.
  emit_system_fields = length(each.value.emit_system_fields) > 0 ? each.value.emit_system_fields : null

  apply_on_transformed_logs = each.value.apply_on_transformed_logs ? true : null
}

# Data protection policy: a single group-scoped policy document (PII audit +
# masking). The Struct spec field arrives as a NESTED OBJECT in tfvars, so it
# is jsonencode'd for the provider's document argument.
resource "aws_cloudwatch_log_data_protection_policy" "this" {
  count = var.spec.data_protection_policy != null ? 1 : 0

  log_group_name  = aws_cloudwatch_log_group.this.name
  policy_document = jsonencode(var.spec.data_protection_policy)
}

# Field index policy: a single group-scoped policy listing the log fields to
# index for faster, cheaper Logs Insights queries.
resource "aws_cloudwatch_log_index_policy" "this" {
  count = var.spec.field_index_policy != null ? 1 : 0

  log_group_name  = aws_cloudwatch_log_group.this.name
  policy_document = jsonencode(var.spec.field_index_policy)
}

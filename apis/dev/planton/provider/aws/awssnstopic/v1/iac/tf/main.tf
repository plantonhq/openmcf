resource "aws_sns_topic" "this" {
  name       = local.topic_name
  fifo_topic = local.is_fifo

  # FIFO-specific settings — AWS rejects these on standard topics, so they are
  # nulled unless fifo_topic is set (CEL already blocks the manifest shape;
  # the guard here keeps the module safe even against hand-built tfvars).
  content_based_deduplication = local.is_fifo ? var.spec.content_based_deduplication : null
  fifo_throughput_scope       = local.is_fifo && var.spec.fifo_throughput_scope != "" ? var.spec.fifo_throughput_scope : null

  # Message archiving (FIFO-only). Subscriptions opt into replay individually
  # via their own replay_policy; the topic only defines the retention window.
  archive_policy = local.archive_policy

  display_name = var.spec.display_name != "" ? var.spec.display_name : null

  # Encryption — customer-managed KMS only (SNS has no managed-SSE option).
  kms_master_key_id = local.kms_master_key_id

  # Resource-based access policy. AWS always keeps a policy on a topic: unset
  # here means the AWS default owner-only policy applies.
  policy = local.policy

  # HTTP/S delivery retry policy.
  delivery_policy = local.delivery_policy

  # Observability.
  tracing_config    = local.tracing_config
  signature_version = local.signature_version

  # Per-protocol delivery-status logging. Failure logging has no sample rate —
  # failures are always logged when the failure role is set. The locals are
  # normalized to zero-value objects, so absent blocks read as empty strings.
  application_success_feedback_role_arn    = local.application_feedback.success_feedback_role != "" ? local.application_feedback.success_feedback_role : null
  application_failure_feedback_role_arn    = local.application_feedback.failure_feedback_role != "" ? local.application_feedback.failure_feedback_role : null
  application_success_feedback_sample_rate = local.application_feedback.success_feedback_sample_rate != 0 ? local.application_feedback.success_feedback_sample_rate : null

  firehose_success_feedback_role_arn    = local.firehose_feedback.success_feedback_role != "" ? local.firehose_feedback.success_feedback_role : null
  firehose_failure_feedback_role_arn    = local.firehose_feedback.failure_feedback_role != "" ? local.firehose_feedback.failure_feedback_role : null
  firehose_success_feedback_sample_rate = local.firehose_feedback.success_feedback_sample_rate != 0 ? local.firehose_feedback.success_feedback_sample_rate : null

  http_success_feedback_role_arn    = local.http_feedback.success_feedback_role != "" ? local.http_feedback.success_feedback_role : null
  http_failure_feedback_role_arn    = local.http_feedback.failure_feedback_role != "" ? local.http_feedback.failure_feedback_role : null
  http_success_feedback_sample_rate = local.http_feedback.success_feedback_sample_rate != 0 ? local.http_feedback.success_feedback_sample_rate : null

  lambda_success_feedback_role_arn    = local.lambda_feedback.success_feedback_role != "" ? local.lambda_feedback.success_feedback_role : null
  lambda_failure_feedback_role_arn    = local.lambda_feedback.failure_feedback_role != "" ? local.lambda_feedback.failure_feedback_role : null
  lambda_success_feedback_sample_rate = local.lambda_feedback.success_feedback_sample_rate != 0 ? local.lambda_feedback.success_feedback_sample_rate : null

  sqs_success_feedback_role_arn    = local.sqs_feedback.success_feedback_role != "" ? local.sqs_feedback.success_feedback_role : null
  sqs_failure_feedback_role_arn    = local.sqs_feedback.failure_feedback_role != "" ? local.sqs_feedback.failure_feedback_role : null
  sqs_success_feedback_sample_rate = local.sqs_feedback.success_feedback_sample_rate != 0 ? local.sqs_feedback.success_feedback_sample_rate : null

  tags = local.aws_tags
}

# Data protection policy is a single-per-topic satellite setting (PII/PHI
# detection with audit/mask/deny operations) keyed by the topic ARN. AWS
# models it as its own API call, so it materializes as its own resource;
# deleting it clears the policy from the topic.
resource "aws_sns_topic_data_protection_policy" "this" {
  count = local.data_protection_policy != null ? 1 : 0

  arn    = aws_sns_topic.this.arn
  policy = local.data_protection_policy
}

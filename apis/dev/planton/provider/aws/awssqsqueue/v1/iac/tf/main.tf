resource "aws_sqs_queue" "this" {
  name       = local.queue_name
  fifo_queue = local.is_fifo

  # Delivery settings — null lets AWS apply its own defaults (see locals.tf).
  visibility_timeout_seconds = local.visibility_timeout_seconds
  message_retention_seconds  = local.message_retention_seconds
  max_message_size           = local.max_message_size
  delay_seconds              = local.delay_seconds
  receive_wait_time_seconds  = local.receive_wait_time_seconds

  # FIFO-specific settings — AWS rejects these on standard queues, so they are
  # nulled unless fifo_queue is set (CEL already blocks the manifest shape;
  # the guard here keeps the module safe even against hand-built tfvars).
  content_based_deduplication = local.is_fifo ? var.spec.content_based_deduplication : null
  deduplication_scope         = local.is_fifo && var.spec.deduplication_scope != "" ? var.spec.deduplication_scope : null
  fifo_throughput_limit       = local.is_fifo && var.spec.fifo_throughput_limit != "" ? var.spec.fifo_throughput_limit : null

  # Dead-letter wiring: redrive_policy points THIS queue's failures at a DLQ;
  # redrive_allow_policy governs who may point at THIS queue as their DLQ.
  redrive_policy       = local.redrive_policy
  redrive_allow_policy = local.redrive_allow_policy

  # Encryption — customer-managed KMS XOR SQS-managed SSE.
  kms_master_key_id                 = local.kms_master_key_id
  kms_data_key_reuse_period_seconds = local.kms_data_key_reuse_period_seconds
  sqs_managed_sse_enabled           = local.sqs_managed_sse_enabled

  # Resource-based access policy (who can SendMessage/ReceiveMessage/etc).
  policy = local.policy

  tags = local.aws_tags
}

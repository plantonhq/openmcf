locals {
  # metadata.name is the queue's cloud name basis. FIFO queues must carry the
  # ".fifo" suffix, so it is appended when absent — manifests stay suffix-free
  # and flipping fifo_queue never requires renaming the resource in YAML.
  resource_name = var.metadata.name
  is_fifo       = var.spec.fifo_queue
  queue_name    = local.is_fifo && !endswith(local.resource_name, ".fifo") ? "${local.resource_name}.fifo" : local.resource_name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSqsQueue"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Delivery settings — 0 means "unset" in the spec, so null is passed to let
  # AWS apply its own defaults (30s visibility, 4-day retention, 256 KiB max
  # message size, no delay, short polling) instead of freezing them here.
  visibility_timeout_seconds = var.spec.visibility_timeout_seconds != 0 ? var.spec.visibility_timeout_seconds : null
  message_retention_seconds  = var.spec.message_retention_seconds != 0 ? var.spec.message_retention_seconds : null
  max_message_size           = var.spec.max_message_size_bytes != 0 ? var.spec.max_message_size_bytes : null
  delay_seconds              = var.spec.delay_seconds != 0 ? var.spec.delay_seconds : null
  receive_wait_time_seconds  = var.spec.receive_wait_time_seconds != 0 ? var.spec.receive_wait_time_seconds : null

  # Dead letter queue — AWS models the DLQ pointer as a JSON redrive policy on
  # the source queue rather than a standalone resource.
  has_dlq = var.spec.dead_letter_config != null
  redrive_policy = local.has_dlq ? jsonencode({
    deadLetterTargetArn = var.spec.dead_letter_config.target_arn
    maxReceiveCount     = var.spec.dead_letter_config.max_receive_count
  }) : null

  # Redrive ALLOW policy — the permission side of the dead-letter relationship:
  # which source queues may point at THIS queue as their DLQ. sourceQueueArns is
  # only accepted alongside the byQueue mode, so it is emitted conditionally
  # rather than as an empty list (AWS rejects allowAll/denyAll documents that
  # carry the key).
  has_redrive_allow = var.spec.redrive_allow_policy != null
  redrive_allow_policy = local.has_redrive_allow ? jsonencode(merge(
    { redrivePermission = var.spec.redrive_allow_policy.redrive_permission },
    var.spec.redrive_allow_policy.redrive_permission == "byQueue" ? {
      sourceQueueArns = var.spec.redrive_allow_policy.source_queue_arns
    } : {}
  )) : null

  # Encryption — customer-managed KMS XOR SQS-managed SSE (CEL enforces the
  # exclusivity; both null lets AWS leave the queue unencrypted).
  kms_master_key_id                 = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  kms_data_key_reuse_period_seconds = var.spec.kms_data_key_reuse_period_seconds != 0 ? var.spec.kms_data_key_reuse_period_seconds : null
  sqs_managed_sse_enabled           = var.spec.sqs_managed_sse_enabled ? true : null

  # Access policy — the Struct arrives from the tfvars layer as a nested
  # object; the provider wants the document as a JSON string.
  policy = var.spec.policy != null ? jsonencode(var.spec.policy) : null
}

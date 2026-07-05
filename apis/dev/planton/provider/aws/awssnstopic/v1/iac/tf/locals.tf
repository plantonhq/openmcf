locals {
  # metadata.name is the topic's cloud name basis. FIFO topics must carry the
  # ".fifo" suffix, so it is appended when absent — manifests stay suffix-free
  # and flipping fifo_topic never requires renaming the resource in YAML.
  resource_name = var.metadata.name
  is_fifo       = var.spec.fifo_topic
  topic_name    = local.is_fifo && !endswith(local.resource_name, ".fifo") ? "${local.resource_name}.fifo" : local.resource_name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSnsTopic"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Encryption — SNS has no managed-SSE option; encryption requires an
  # explicit customer-managed KMS key.
  kms_master_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Policy documents arrive from the tfvars layer as nested objects (the
  # specs model them as Structs); the provider wants JSON strings.
  policy                 = var.spec.policy != null ? jsonencode(var.spec.policy) : null
  archive_policy         = var.spec.archive_policy != null ? jsonencode(var.spec.archive_policy) : null
  data_protection_policy = var.spec.data_protection_policy != null ? jsonencode(var.spec.data_protection_policy) : null

  # Delivery policy (HTTP/S retry behavior) is already a raw JSON string in
  # the spec.
  delivery_policy = var.spec.delivery_policy != "" ? var.spec.delivery_policy : null

  # Observability — zero values mean "let AWS default" (PassThrough tracing,
  # signature version 1).
  tracing_config    = var.spec.tracing_config != "" ? var.spec.tracing_config : null
  signature_version = var.spec.signature_version != 0 ? var.spec.signature_version : null

  # Per-protocol delivery-status logging. Each protocol block is an
  # independent opt-in; an absent block leaves that protocol's feedback
  # attributes unset. Sample rates of 0 are passed as null so AWS applies its
  # own default rather than freezing "log nothing". Each block is normalized
  # to a zero-value object here (via try, since HCL's && does not
  # short-circuit on null) so main.tf can read attributes unconditionally.
  feedback_zero = {
    success_feedback_role        = ""
    failure_feedback_role        = ""
    success_feedback_sample_rate = 0
  }

  application_feedback = try(coalesce(var.spec.delivery_feedback.application), local.feedback_zero)
  firehose_feedback    = try(coalesce(var.spec.delivery_feedback.firehose), local.feedback_zero)
  http_feedback        = try(coalesce(var.spec.delivery_feedback.http), local.feedback_zero)
  lambda_feedback      = try(coalesce(var.spec.delivery_feedback.lambda), local.feedback_zero)
  sqs_feedback         = try(coalesce(var.spec.delivery_feedback.sqs), local.feedback_zero)
}

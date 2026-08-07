locals {
  # metadata.name is the stream's cloud name (ForceNew — Kinesis streams
  # cannot be renamed).
  stream_name = var.metadata.name

  # Resource-identity tags follow the catalog convention.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsKinesisStream"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Capacity mode — shard_count only applies to PROVISIONED streams (CEL
  # enforces the coupling); ON_DEMAND streams must not send one.
  stream_mode = var.spec.stream_mode
  shard_count = local.stream_mode == "PROVISIONED" ? var.spec.shard_count : null

  # Warm throughput — pre-provisioned burst capacity for ON_DEMAND streams.
  # 0 means "unset": on-demand scaling manages capacity reactively.
  warm_throughput_mib_ps = var.spec.warm_throughput_mib_ps != 0 ? var.spec.warm_throughput_mib_ps : null

  # Data retention — null lets AWS apply its default (24h) instead of
  # freezing it here.
  retention_period = var.spec.retention_period_hours != 0 ? var.spec.retention_period_hours : null

  # Encryption — presence of a KMS key implies KMS encryption; absence means
  # the stream is unencrypted (Kinesis has no account-level default the way
  # S3 does).
  kms_key_id      = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  encryption_type = local.kms_key_id != null ? "KMS" : "NONE"

  # Maximum record size — 0 means "unset" (the AWS default of 1 MiB).
  # Requires the kinesis:UpdateMaxRecordSize IAM permission when set.
  max_record_size_in_kib = var.spec.max_record_size_in_kib != 0 ? var.spec.max_record_size_in_kib : null

  # Enhanced monitoring — null (not an empty set) when no shard-level
  # metrics are requested, so the provider leaves monitoring untouched.
  shard_level_metrics = length(var.spec.shard_level_metrics) > 0 ? var.spec.shard_level_metrics : null

  # Resource policy — the Struct arrives from the tfvars layer as a nested
  # object; the provider wants the document as a JSON string.
  resource_policy = var.spec.resource_policy != null ? jsonencode(var.spec.resource_policy) : null
}

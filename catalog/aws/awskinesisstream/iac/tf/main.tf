resource "aws_kinesis_stream" "this" {
  name = local.stream_name

  # Capacity mode. AWS defaults a bare stream to PROVISIONED; the spec makes
  # the choice explicit and required, so the block is always emitted.
  stream_mode_details {
    stream_mode = local.stream_mode
  }

  # Only PROVISIONED streams carry a shard count; ON_DEMAND streams instead
  # may pre-provision warm burst capacity (the two are mutually exclusive —
  # the provider's own ConflictsWith, mirrored as a spec CEL rule).
  shard_count            = local.shard_count
  warm_throughput_mib_ps = local.warm_throughput_mib_ps

  retention_period = local.retention_period

  # Encryption is derived: a KMS key means KMS encryption, absence means
  # NONE. Kinesis accepts key id/ARN/alias, including "alias/aws/kinesis"
  # for the AWS-owned key.
  encryption_type = local.encryption_type
  kms_key_id      = local.kms_key_id

  # Kinesis 10 MiB large-record support. In-place update; needs the
  # kinesis:UpdateMaxRecordSize IAM permission on the deploying principal.
  max_record_size_in_kib = local.max_record_size_in_kib

  shard_level_metrics = local.shard_level_metrics

  # Deletion-time behavior only: deregister enhanced fan-out consumers
  # instead of failing the delete.
  enforce_consumer_deletion = var.spec.enforce_consumer_deletion

  tags = local.aws_tags
}

# Resource-based access policy — AWS models this as a separate API keyed by
# the stream ARN (one policy per stream), folded into the spec because it has
# no identity of its own. The primary use is cross-account producer/consumer
# grants without role assumption.
resource "aws_kinesis_resource_policy" "this" {
  count = local.resource_policy != null ? 1 : 0

  resource_arn = aws_kinesis_stream.this.arn
  policy       = local.resource_policy
}

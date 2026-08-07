# An enhanced fan-out consumer is a replace-only resource: both its name and
# its stream ARN are ForceNew, and AWS manages everything else about it.
resource "aws_kinesis_stream_consumer" "this" {
  name       = local.consumer_name
  stream_arn = var.spec.stream_arn

  tags = local.aws_tags
}

# Resource-based access policy — AWS models this as a separate API keyed by
# the consumer ARN (one policy per consumer), folded into the spec because it
# has no identity of its own. The primary use is cross-account enhanced
# fan-out: SubscribeToShard grants without role assumption.
resource "aws_kinesis_resource_policy" "this" {
  count = local.resource_policy != null ? 1 : 0

  resource_arn = aws_kinesis_stream_consumer.this.arn
  policy       = local.resource_policy
}

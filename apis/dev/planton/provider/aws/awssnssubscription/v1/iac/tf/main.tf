resource "aws_sns_topic_subscription" "this" {
  # The immutable identity trio: changing any of these replaces the
  # subscription (AWS has no repoint operation for topic/protocol/endpoint).
  topic_arn = var.spec.topic_arn
  protocol  = var.spec.protocol
  endpoint  = var.spec.endpoint

  # Message filtering. filter_policy_scope is only sent alongside a filter
  # policy — AWS rejects a scope without a policy, and CEL blocks the
  # manifest shape; the guard here keeps the module safe against hand-built
  # tfvars.
  filter_policy       = local.filter_policy
  filter_policy_scope = local.filter_policy != null && var.spec.filter_policy_scope != "" ? var.spec.filter_policy_scope : null

  # Delivery behavior.
  raw_message_delivery = var.spec.raw_message_delivery
  redrive_policy       = local.redrive_policy
  delivery_policy      = local.delivery_policy

  # Replay of archived messages (FIFO topics with an archive policy) — the
  # backfill mechanism for a consumer added after messages were published.
  replay_policy = local.replay_policy

  # Firehose delivery requires the role SNS assumes to write to the stream.
  subscription_role_arn = var.spec.subscription_role_arn != "" ? var.spec.subscription_role_arn : null

  # HTTP/S confirmation handshake. Both attributes are meaningless for
  # auto-confirming protocols (SQS/Lambda/Firehose/application) and CEL keeps
  # them off those manifests.
  endpoint_auto_confirms          = var.spec.endpoint_auto_confirms
  confirmation_timeout_in_minutes = var.spec.confirmation_timeout_minutes != 0 ? var.spec.confirmation_timeout_minutes : null
}

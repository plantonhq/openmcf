locals {
  # A subscription has no cloud-side name — AWS identifies it by a
  # server-assigned ARN. metadata.name only drives the Planton resource
  # identity; no naming basis is rendered to AWS. Subscriptions are also
  # untaggable in AWS, so the module carries no tag block (the identity tags
  # live on the topic and the endpoint resources).

  # Dead letter queue — AWS models the subscription DLQ as a JSON redrive
  # policy document holding just the target ARN (retry exhaustion is governed
  # by the delivery policy, not a receive count).
  redrive_policy = var.spec.dead_letter_config != null ? jsonencode({
    deadLetterTargetArn = var.spec.dead_letter_config.dead_letter_target_arn
  }) : null

  # Policy documents arrive from the tfvars layer as nested objects (the
  # specs model them as Structs); the provider wants JSON strings.
  filter_policy = var.spec.filter_policy != null ? jsonencode(var.spec.filter_policy) : null
  replay_policy = var.spec.replay_policy != null ? jsonencode(var.spec.replay_policy) : null

  # delivery_policy is already a raw JSON string in the spec.
  delivery_policy = var.spec.delivery_policy != "" ? var.spec.delivery_policy : null
}

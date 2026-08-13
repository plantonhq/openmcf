locals {
  # The guardrail name is metadata.name -- the naming basis both engines
  # share so a manifest deploys identically on either. AWS allows 1-50
  # characters of alphanumeric plus - and _ (no spaces or dots).
  guardrail_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockGuardrail"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Policy families render only when the manifest declares them.
  has_content_policy              = var.spec.content_policy != null
  has_topic_policy                = var.spec.topic_policy != null
  has_word_policy                 = var.spec.word_policy != null
  has_sensitive_information       = var.spec.sensitive_information_policy != null
  has_contextual_grounding_policy = var.spec.contextual_grounding_policy != null
  has_cross_region                = var.spec.cross_region_profile_arn != ""

  # Published versions keyed by their stable entry name (the for_each key
  # both engines share; AWS assigns the actual version number).
  versions = { for v in var.spec.versions : v.name => v }
}

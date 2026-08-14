# Amazon Bedrock application inference profile: a named, taggable,
# IAM-scopeable handle over a foundation model (or an AWS system-defined
# cross-region profile) for per-application usage tracking and cost
# allocation.
#
# Every argument is create-time-immutable -- changing one replaces the
# profile and its ARN.

resource "aws_bedrock_inference_profile" "this" {
  name = local.profile_name

  # Sent only when set (1-200 characters; ForceNew at the provider).
  description = var.spec.description != "" ? var.spec.description : null

  # The model source this profile routes to. AWS never echoes it back
  # (GetInferenceProfile reports the resolved models list instead), so the
  # provider pins the configured value in state.
  model_source {
    copy_from = var.spec.source_arn
  }

  tags = local.aws_tags
}

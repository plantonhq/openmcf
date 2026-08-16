# Bedrock model invocation logging for one region -- a settings
# singleton: AWS keeps exactly one logging configuration per
# account+region (the resource id IS the region), and this module
# manages it. metadata.name never reaches AWS.
#
# Lifecycle facts the render below depends on:
#   - the four data-type toggles default to TRUE upstream; the spec's
#     presence-typed optional bools pass through as-is (null lets the
#     provider apply AWS's enabled default, explicit false is sent);
#   - at least one delivery destination is required by the spec's CEL
#     (a destination-less configuration delivers nothing), so the two
#     dynamic blocks can never BOTH be empty;
#   - AWS validates the CloudWatch role's permission chain at apply
#     ("Failed to validate permissions for log group") and the
#     provider retries through IAM propagation lag;
#   - destroy DELETES the configuration -- the region reverts to no
#     invocation logging.

resource "aws_bedrock_model_invocation_logging_configuration" "this" {
  logging_config {
    embedding_data_delivery_enabled = var.spec.embedding_data_delivery_enabled
    image_data_delivery_enabled     = var.spec.image_data_delivery_enabled
    text_data_delivery_enabled      = var.spec.text_data_delivery_enabled
    video_data_delivery_enabled     = var.spec.video_data_delivery_enabled

    dynamic "cloudwatch_config" {
      for_each = var.spec.cloudwatch != null ? [var.spec.cloudwatch] : []
      content {
        log_group_name = cloudwatch_config.value.log_group_name
        role_arn       = cloudwatch_config.value.role_arn

        # Where CloudWatch delivery spills payloads larger than a log
        # event (256 KB); without it oversized bodies are truncated.
        dynamic "large_data_delivery_s3_config" {
          for_each = cloudwatch_config.value.large_data_delivery_s3 != null ? [cloudwatch_config.value.large_data_delivery_s3] : []
          content {
            bucket_name = large_data_delivery_s3_config.value.bucket_name
            key_prefix  = large_data_delivery_s3_config.value.key_prefix != "" ? large_data_delivery_s3_config.value.key_prefix : null
          }
        }
      }
    }

    dynamic "s3_config" {
      for_each = var.spec.s3 != null ? [var.spec.s3] : []
      content {
        bucket_name = s3_config.value.bucket_name
        key_prefix  = s3_config.value.key_prefix != "" ? s3_config.value.key_prefix : null
      }
    }
  }
}

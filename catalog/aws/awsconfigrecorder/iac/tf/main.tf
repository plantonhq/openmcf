# The region's AWS Config recording posture: one configuration
# recorder, its delivery channel, the recorder's on/off state, and the
# configuration-item retention window.
#
# Lifecycle facts the renders below depend on:
#   - AWS allows ONE recorder and ONE delivery channel per region, both
#     named "default" by convention -- the names are hardcoded here and
#     metadata.name never reaches AWS (the settings-singleton
#     contract);
#   - the delivery channel cannot be created before a recorder exists,
#     and cannot be DELETED while the recorder is running -- the
#     depends_on below fixes create order, and the provider retries
#     channel deletion for 30s while the recorder stop lands;
#   - the recorder-status resource is the folded start/stop toggle:
#     Create/Update/Delete all map to Start/StopConfigurationRecorder;
#   - starting a recorder without a delivery channel fails
#     (NoAvailableDeliveryChannelException) -- the spec CEL guarantees
#     a channel arrives whenever recording_enabled is not false;
#   - AWS validates the delivery bucket's POLICY ("insufficient
#     delivery policy") -- the bucket policy is the consumer's
#     contract (AwsS3Bucket spec.policy), never this module's.

resource "aws_config_configuration_recorder" "this" {
  # The regional singleton's conventional name -- an AWS-side
  # constant, deliberately not configurable.
  name     = "default"
  role_arn = var.spec.role_arn

  dynamic "recording_group" {
    for_each = var.spec.recording_group != null ? [var.spec.recording_group] : []
    content {
      all_supported                 = recording_group.value.all_supported
      include_global_resource_types = recording_group.value.include_global_resource_types
      resource_types                = length(recording_group.value.resource_types) > 0 ? recording_group.value.resource_types : null

      dynamic "exclusion_by_resource_types" {
        for_each = length(recording_group.value.exclusion_by_resource_types) > 0 ? [recording_group.value.exclusion_by_resource_types] : []
        content {
          resource_types = exclusion_by_resource_types.value
        }
      }

      dynamic "recording_strategy" {
        for_each = recording_group.value.recording_strategy != "" ? [recording_group.value.recording_strategy] : []
        content {
          use_only = recording_strategy.value
        }
      }
    }
  }

  dynamic "recording_mode" {
    for_each = var.spec.recording_mode != null ? [var.spec.recording_mode] : []
    content {
      recording_frequency = recording_mode.value.recording_frequency != "" ? recording_mode.value.recording_frequency : "CONTINUOUS"

      dynamic "recording_mode_override" {
        for_each = recording_mode.value.override != null ? [recording_mode.value.override] : []
        content {
          description         = recording_mode_override.value.description != "" ? recording_mode_override.value.description : null
          recording_frequency = recording_mode_override.value.recording_frequency
          resource_types      = recording_mode_override.value.resource_types
        }
      }
    }
  }
}

resource "aws_config_delivery_channel" "this" {
  count = var.spec.delivery_channel != null ? 1 : 0

  # The regional singleton's conventional name (see the recorder).
  name           = "default"
  s3_bucket_name = var.spec.delivery_channel.s3_bucket_name
  s3_key_prefix  = var.spec.delivery_channel.s3_key_prefix != "" ? var.spec.delivery_channel.s3_key_prefix : null
  s3_kms_key_arn = var.spec.delivery_channel.s3_kms_key_arn != "" ? var.spec.delivery_channel.s3_kms_key_arn : null
  sns_topic_arn  = var.spec.delivery_channel.sns_topic_arn != "" ? var.spec.delivery_channel.sns_topic_arn : null

  dynamic "snapshot_delivery_properties" {
    for_each = var.spec.delivery_channel.snapshot_delivery_frequency != "" ? [var.spec.delivery_channel.snapshot_delivery_frequency] : []
    content {
      delivery_frequency = snapshot_delivery_properties.value
    }
  }

  # AWS refuses a channel without a recorder.
  depends_on = [aws_config_configuration_recorder.this]
}

# The folded recorder toggle: unset recording_enabled means RUNNING
# (the reason this component exists).
resource "aws_config_configuration_recorder_status" "this" {
  name       = aws_config_configuration_recorder.this.name
  is_enabled = var.spec.recording_enabled == null ? true : var.spec.recording_enabled

  # Starting requires the delivery channel to exist.
  depends_on = [aws_config_delivery_channel.this]
}

# The retention singleton (AWS names it "default"; the name is
# API-computed and cannot be chosen). Managed only when the spec sets
# a window.
resource "aws_config_retention_configuration" "this" {
  count = var.spec.retention_period_in_days > 0 ? 1 : 0

  retention_period_in_days = var.spec.retention_period_in_days
}

# The ECS cluster itself: Container Insights, ECS Exec auditing, Fargate
# storage encryption, and the Service Connect default namespace. Only the
# name is create-time; every posture block edits in place via
# UpdateCluster.
#
# Nested optional OBJECTS from the tfvars pipeline are guarded with
# ternaries / dynamic-block for_each, never `!= null &&` -- HCL's && does
# not short-circuit and errors on the null dereference.
resource "aws_ecs_cluster" "this" {
  name = local.cluster_name

  # Unset keeps the account's default setting (AWS computes the
  # effective value); the spec values are the exact AWS API strings.
  dynamic "setting" {
    for_each = var.spec.container_insights != "" ? [1] : []
    content {
      name  = "containerInsights"
      value = var.spec.container_insights
    }
  }

  # Exec auditing and managed-storage encryption share the provider's
  # configuration block; build it only when either is present so an
  # empty block never overwrites AWS defaults.
  dynamic "configuration" {
    for_each = var.spec.execute_command_configuration != null || var.spec.managed_storage_configuration != null ? [1] : []
    content {
      dynamic "execute_command_configuration" {
        for_each = var.spec.execute_command_configuration != null ? [var.spec.execute_command_configuration] : []
        content {
          logging    = execute_command_configuration.value.logging != "" ? execute_command_configuration.value.logging : null
          kms_key_id = execute_command_configuration.value.kms_key_id != "" ? execute_command_configuration.value.kms_key_id : null

          # Present only with OVERRIDE logging (CEL enforces the
          # coupling both directions).
          dynamic "log_configuration" {
            for_each = execute_command_configuration.value.log_configuration != null ? [execute_command_configuration.value.log_configuration] : []
            content {
              cloud_watch_log_group_name     = log_configuration.value.cloud_watch_log_group_name != "" ? log_configuration.value.cloud_watch_log_group_name : null
              cloud_watch_encryption_enabled = log_configuration.value.cloud_watch_encryption_enabled
              s3_bucket_name                 = log_configuration.value.s3_bucket_name != "" ? log_configuration.value.s3_bucket_name : null
              s3_key_prefix                  = log_configuration.value.s3_key_prefix != "" ? log_configuration.value.s3_key_prefix : null
              s3_bucket_encryption_enabled   = log_configuration.value.s3_bucket_encryption_enabled
            }
          }
        }
      }

      dynamic "managed_storage_configuration" {
        for_each = var.spec.managed_storage_configuration != null ? [var.spec.managed_storage_configuration] : []
        content {
          fargate_ephemeral_storage_kms_key_id = managed_storage_configuration.value.fargate_ephemeral_storage_kms_key_id != "" ? managed_storage_configuration.value.fargate_ephemeral_storage_kms_key_id : null
          kms_key_id                           = managed_storage_configuration.value.kms_key_id != "" ? managed_storage_configuration.value.kms_key_id : null
        }
      }
    }
  }

  # One environment-wide Service Connect namespace, overridable per
  # service -- set here so services join the mesh with zero wiring.
  dynamic "service_connect_defaults" {
    for_each = var.spec.service_connect_namespace_arn != "" ? [1] : []
    content {
      namespace = var.spec.service_connect_namespace_arn
    }
  }

  tags = local.aws_tags
}

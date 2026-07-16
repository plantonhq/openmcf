# AWS ECS Task Definition Terraform Module
#
# Registers one immutable revision of the task-definition family per apply.
# The revision-carrying ARN output is what lets a referencing ECS service
# roll on each change -- "new image tag in, service rolls" is the composed
# behavior, not module magic.

# The shared CloudWatch log group for the task-level logging default. Only
# created when logging is enabled AND no existing group is referenced; a
# referenced group owns its own retention and lifecycle.
resource "aws_cloudwatch_log_group" "this" {
  count             = local.create_log_group ? 1 : 0
  name              = local.log_group_name
  retention_in_days = local.log_retention_days
  tags              = local.aws_tags
}

resource "aws_ecs_task_definition" "this" {
  family                   = local.family
  container_definitions    = jsonencode(local.container_definitions)
  requires_compatibilities = local.requires_compatibilities
  network_mode             = local.network_mode

  # AWS takes task-level sizing as strings; null omits them (EC2/EXTERNAL
  # tasks may size per container instead -- the spec's CEL guarantees
  # Fargate tasks always carry both).
  cpu    = local.task_cpu
  memory = local.task_memory

  # Two roles by design: the agent's setup identity (pull images, fetch
  # secrets, write logs) and the application's runtime identity stay
  # separate so neither accumulates the other's permissions.
  execution_role_arn = local.execution_role_arn
  task_role_arn      = local.task_role_arn

  # runtime_platform is a pruned optional message: guard with a ternary
  # (HCL's && does not short-circuit on null).
  dynamic "runtime_platform" {
    for_each = var.spec.runtime_platform != null ? [var.spec.runtime_platform] : []
    content {
      cpu_architecture        = runtime_platform.value.cpu_architecture != "" ? runtime_platform.value.cpu_architecture : null
      operating_system_family = runtime_platform.value.operating_system_family != "" ? runtime_platform.value.operating_system_family : null
    }
  }

  # Fargate includes 20 GiB at no charge; the block is only sent when the
  # workload needs more.
  dynamic "ephemeral_storage" {
    for_each = var.spec.ephemeral_storage_gib > 0 ? [var.spec.ephemeral_storage_gib] : []
    content {
      size_in_gib = ephemeral_storage.value
    }
  }

  dynamic "volume" {
    for_each = var.spec.volumes
    content {
      name = volume.value.name
      # host_path and efs are mutually exclusive (CEL-enforced); an empty
      # host_path must be null, not "", or AWS rejects the registration.
      host_path = volume.value.host_path != "" ? volume.value.host_path : null

      dynamic "efs_volume_configuration" {
        for_each = volume.value.efs != null ? [volume.value.efs] : []
        content {
          file_system_id = efs_volume_configuration.value.file_system_id
          root_directory = efs_volume_configuration.value.root_directory != "" ? efs_volume_configuration.value.root_directory : null
          # Transit encryption is always on: AWS requires it with access
          # points or IAM auth, and there is no good reason to mount EFS
          # unencrypted in transit without them.
          transit_encryption = "ENABLED"

          dynamic "authorization_config" {
            for_each = (efs_volume_configuration.value.access_point_id != "" || efs_volume_configuration.value.iam_authorization) ? [efs_volume_configuration.value] : []
            content {
              access_point_id = authorization_config.value.access_point_id != "" ? authorization_config.value.access_point_id : null
              iam             = authorization_config.value.iam_authorization ? "ENABLED" : null
            }
          }
        }
      }
    }
  }

  # Keep old revisions registered on destroy when other consumers (a
  # scheduled task, a manual RunTask) may still reference them.
  skip_destroy = var.spec.skip_destroy

  tags = local.aws_tags

  # The group must exist before any task launches from this revision -- the
  # awslogs driver fails at task START (not registration) when it is
  # missing, which no offline check can catch.
  depends_on = [aws_cloudwatch_log_group.this]
}

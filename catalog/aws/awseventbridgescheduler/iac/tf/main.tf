# One EventBridge Scheduler schedule, with an optionally owned
# schedule group.
#
# Lifecycle facts the render below depends on:
#   - the schedule's name (metadata.name) and its group are fixed for
#     life (replace-on-change);
#   - the group is a name-and-tags container - its own update path is
#     tags-only; an owned group carries the identity tags (the
#     schedule itself is UNTAGGABLE at AWS - the deliberate
#     tag-convention absence);
#   - a first deploy with a freshly created execution role is
#     eventually consistent - the provider retries "must allow AWS
#     EventBridge Scheduler to assume the role" for up to two minutes;
#   - schedule_expression_timezone defaults to UTC at the provider;
#   - with action_after_completion = DELETE, AWS deletes a completed
#     one-time schedule out from under IaC state - the next deploy
#     recreates it (fire-and-forget only).

# --- owned group (optional) ------------------------------------------------

resource "aws_scheduler_schedule_group" "this" {
  count = var.spec.group != null ? 1 : 0

  name = var.spec.group.name
  tags = local.aws_tags
}

# --- schedule ----------------------------------------------------------------

resource "aws_scheduler_schedule" "this" {
  name = var.metadata.name

  # Owned group -> its name; else the joined group_name; else AWS's
  # "default" group (null).
  group_name = var.spec.group != null ? aws_scheduler_schedule_group.this[0].name : (var.spec.group_name != "" ? var.spec.group_name : null)

  description                  = var.spec.description != "" ? var.spec.description : null
  schedule_expression          = var.spec.schedule_expression
  schedule_expression_timezone = var.spec.schedule_expression_timezone != "" ? var.spec.schedule_expression_timezone : null
  start_date                   = var.spec.start_date != "" ? var.spec.start_date : null
  end_date                     = var.spec.end_date != "" ? var.spec.end_date : null
  state                        = var.spec.state != "" ? var.spec.state : null
  action_after_completion      = var.spec.action_after_completion != "" ? var.spec.action_after_completion : null
  kms_key_arn                  = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  flexible_time_window {
    mode                      = var.spec.flexible_time_window.mode
    maximum_window_in_minutes = var.spec.flexible_time_window.maximum_window_in_minutes != null ? var.spec.flexible_time_window.maximum_window_in_minutes : null
  }

  target {
    arn      = var.spec.target.arn
    role_arn = var.spec.target.role_arn
    input    = var.spec.target.input != "" ? var.spec.target.input : null

    dynamic "dead_letter_config" {
      for_each = var.spec.target.dead_letter_queue_arn != "" ? [1] : []
      content {
        arn = var.spec.target.dead_letter_queue_arn
      }
    }

    dynamic "retry_policy" {
      for_each = var.spec.target.retry_policy != null ? [var.spec.target.retry_policy] : []
      content {
        maximum_event_age_in_seconds = retry_policy.value.maximum_event_age_in_seconds != null ? retry_policy.value.maximum_event_age_in_seconds : null
        maximum_retry_attempts       = retry_policy.value.maximum_retry_attempts != null ? retry_policy.value.maximum_retry_attempts : null
      }
    }

    dynamic "ecs_parameters" {
      for_each = var.spec.target.ecs_parameters != null ? [var.spec.target.ecs_parameters] : []
      content {
        task_definition_arn     = ecs_parameters.value.task_definition_arn
        task_count              = ecs_parameters.value.task_count != null ? ecs_parameters.value.task_count : null
        launch_type             = ecs_parameters.value.launch_type != "" ? ecs_parameters.value.launch_type : null
        group                   = ecs_parameters.value.group != "" ? ecs_parameters.value.group : null
        platform_version        = ecs_parameters.value.platform_version != "" ? ecs_parameters.value.platform_version : null
        propagate_tags          = ecs_parameters.value.propagate_tags != "" ? ecs_parameters.value.propagate_tags : null
        reference_id            = ecs_parameters.value.reference_id != "" ? ecs_parameters.value.reference_id : null
        enable_ecs_managed_tags = ecs_parameters.value.enable_ecs_managed_tags
        enable_execute_command  = ecs_parameters.value.enable_execute_command
        tags                    = length(ecs_parameters.value.tags) > 0 ? ecs_parameters.value.tags : null

        dynamic "capacity_provider_strategy" {
          for_each = ecs_parameters.value.capacity_provider_strategy != null ? ecs_parameters.value.capacity_provider_strategy : []
          content {
            capacity_provider = capacity_provider_strategy.value.capacity_provider
            base              = capacity_provider_strategy.value.base
            weight            = capacity_provider_strategy.value.weight
          }
        }

        dynamic "network_configuration" {
          for_each = ecs_parameters.value.network_configuration != null ? [ecs_parameters.value.network_configuration] : []
          content {
            subnets          = network_configuration.value.subnets
            security_groups  = network_configuration.value.security_groups != null && length(network_configuration.value.security_groups) > 0 ? network_configuration.value.security_groups : null
            assign_public_ip = network_configuration.value.assign_public_ip
          }
        }

        dynamic "placement_constraints" {
          for_each = ecs_parameters.value.placement_constraints != null ? ecs_parameters.value.placement_constraints : []
          content {
            type       = placement_constraints.value.type
            expression = placement_constraints.value.expression != "" ? placement_constraints.value.expression : null
          }
        }

        dynamic "placement_strategy" {
          for_each = ecs_parameters.value.placement_strategy != null ? ecs_parameters.value.placement_strategy : []
          content {
            type  = placement_strategy.value.type
            field = placement_strategy.value.field != "" ? placement_strategy.value.field : null
          }
        }
      }
    }

    dynamic "eventbridge_parameters" {
      for_each = var.spec.target.eventbridge_parameters != null ? [var.spec.target.eventbridge_parameters] : []
      content {
        detail_type = eventbridge_parameters.value.detail_type
        source      = eventbridge_parameters.value.source
      }
    }

    dynamic "kinesis_parameters" {
      for_each = var.spec.target.kinesis_parameters != null ? [var.spec.target.kinesis_parameters] : []
      content {
        partition_key = kinesis_parameters.value.partition_key
      }
    }

    dynamic "sagemaker_pipeline_parameters" {
      for_each = var.spec.target.sagemaker_pipeline_parameters != null ? [var.spec.target.sagemaker_pipeline_parameters] : []
      content {
        dynamic "pipeline_parameter" {
          for_each = sagemaker_pipeline_parameters.value.pipeline_parameters != null ? sagemaker_pipeline_parameters.value.pipeline_parameters : []
          content {
            name  = pipeline_parameter.value.name
            value = pipeline_parameter.value.value
          }
        }
      }
    }

    dynamic "sqs_parameters" {
      for_each = var.spec.target.sqs_parameters != null ? [var.spec.target.sqs_parameters] : []
      content {
        message_group_id = sqs_parameters.value.message_group_id != "" ? sqs_parameters.value.message_group_id : null
      }
    }
  }
}

# ---- EventBridge Rule ----

resource "aws_cloudwatch_event_rule" "this" {
  name                = local.resource_name
  description         = var.spec.description != "" ? var.spec.description : null
  event_bus_name      = local.event_bus_name
  event_pattern       = local.event_pattern
  schedule_expression = local.schedule_expression
  state               = local.state

  # Rule-level invocation role (per-target role_arn takes precedence for its
  # own target).
  role_arn = var.spec.role_arn != "" ? var.spec.role_arn : null

  # AWS refuses to delete a rule that still has targets unless forced. The
  # module removes its own targets first, so force only matters when an
  # out-of-band consumer attached extra targets to this rule.
  force_destroy = var.spec.force_destroy

  tags = local.aws_tags
}

# ---- EventBridge Targets ----
# Many-per-rule, keyed by target name (the EventBridge target_id). Exactly one
# service-typed block may accompany each target (CEL-enforced); each dynamic
# block below renders only for its own type.

resource "aws_cloudwatch_event_target" "this" {
  for_each = local.targets

  rule           = aws_cloudwatch_event_rule.this.name
  event_bus_name = local.event_bus_name
  target_id      = each.key
  arn            = each.value.arn

  role_arn = each.value.role_arn != "" ? each.value.role_arn : null

  # Input transformation (mutually exclusive — CEL-enforced).
  input      = each.value.input != "" ? each.value.input : null
  input_path = each.value.input_path != "" ? each.value.input_path : null

  dynamic "input_transformer" {
    for_each = each.value.input_transformer != null ? [each.value.input_transformer] : []
    content {
      input_paths    = input_transformer.value.input_paths
      input_template = input_transformer.value.input_template
    }
  }

  # Per-target DLQ for events that exhaust the retry policy.
  dynamic "dead_letter_config" {
    for_each = each.value.dead_letter_config != null ? [each.value.dead_letter_config] : []
    content {
      arn = dead_letter_config.value.arn
    }
  }

  # Retry policy: both fields are presence-aware (null = AWS default) because
  # zero retry attempts is a meaningful setting (fail straight to the DLQ).
  dynamic "retry_policy" {
    for_each = each.value.retry_policy != null ? [each.value.retry_policy] : []
    content {
      maximum_event_age_in_seconds = retry_policy.value.maximum_event_age_in_seconds
      maximum_retry_attempts       = retry_policy.value.maximum_retry_attempts
    }
  }

  # SQS: message group for FIFO queues.
  dynamic "sqs_target" {
    for_each = each.value.sqs_target != null ? [each.value.sqs_target] : []
    content {
      message_group_id = sqs_target.value.message_group_id != "" ? sqs_target.value.message_group_id : null
    }
  }

  # Kinesis: shard routing via a partition-key JSONPath.
  dynamic "kinesis_target" {
    for_each = each.value.kinesis_target != null ? [each.value.kinesis_target] : []
    content {
      partition_key_path = kinesis_target.value.partition_key_path != "" ? kinesis_target.value.partition_key_path : null
    }
  }

  # API destination: path/query/header parameters for the HTTP invocation.
  dynamic "http_target" {
    for_each = each.value.http_target != null ? [each.value.http_target] : []
    content {
      path_parameter_values   = http_target.value.path_parameter_values
      query_string_parameters = http_target.value.query_string_parameters
      header_parameters       = http_target.value.header_parameters
    }
  }

  # AWS Batch: job submission parameters.
  dynamic "batch_target" {
    for_each = each.value.batch_target != null ? [each.value.batch_target] : []
    content {
      job_definition = batch_target.value.job_definition
      job_name       = batch_target.value.job_name
      array_size     = batch_target.value.array_size != 0 ? batch_target.value.array_size : null
      job_attempts   = batch_target.value.job_attempts != 0 ? batch_target.value.job_attempts : null
    }
  }

  # ECS RunTask: the target arn is the CLUSTER; the task definition, sizing,
  # networking, and placement live in this block.
  dynamic "ecs_target" {
    for_each = each.value.ecs_target != null ? [each.value.ecs_target] : []
    content {
      task_definition_arn     = ecs_target.value.task_definition_arn
      task_count              = ecs_target.value.task_count != 0 ? ecs_target.value.task_count : null
      launch_type             = ecs_target.value.launch_type != "" ? ecs_target.value.launch_type : null
      platform_version        = ecs_target.value.platform_version != "" ? ecs_target.value.platform_version : null
      group                   = ecs_target.value.group != "" ? ecs_target.value.group : null
      propagate_tags          = ecs_target.value.propagate_tags != "" ? ecs_target.value.propagate_tags : null
      enable_ecs_managed_tags = ecs_target.value.enable_ecs_managed_tags
      enable_execute_command  = ecs_target.value.enable_execute_command

      dynamic "capacity_provider_strategy" {
        for_each = ecs_target.value.capacity_provider_strategy
        content {
          capacity_provider = capacity_provider_strategy.value.capacity_provider
          base              = capacity_provider_strategy.value.base
          weight            = capacity_provider_strategy.value.weight
        }
      }

      dynamic "network_configuration" {
        for_each = ecs_target.value.network_configuration != null ? [ecs_target.value.network_configuration] : []
        content {
          subnets          = network_configuration.value.subnets
          security_groups  = network_configuration.value.security_groups
          assign_public_ip = network_configuration.value.assign_public_ip
        }
      }

      dynamic "ordered_placement_strategy" {
        for_each = ecs_target.value.ordered_placement_strategy
        content {
          type  = ordered_placement_strategy.value.type
          field = ordered_placement_strategy.value.field != "" ? ordered_placement_strategy.value.field : null
        }
      }

      dynamic "placement_constraint" {
        for_each = ecs_target.value.placement_constraints
        content {
          type       = placement_constraint.value.type
          expression = placement_constraint.value.expression != "" ? placement_constraint.value.expression : null
        }
      }
    }
  }
}

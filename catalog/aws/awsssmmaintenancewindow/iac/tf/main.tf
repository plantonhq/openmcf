# An SSM maintenance window with its folded target registrations and
# tasks (Run Command / Automation / Lambda / Step Functions).
#
# Lifecycle facts the render below depends on:
#   - enabled is a provider-default-TRUE toggle - rendered only on an
#     explicit choice so the module never fights the default (unset =
#     enabled; the provider needs a second API call to create a window
#     paused);
#   - targets and tasks are true window satellites: window_id (and the
#     target's name/description/resource_type, and the task's
#     task_type) force replacement of the registration;
#   - rate controls (max_concurrency/max_errors) are only legal on a
#     task WITH targets - AWS rejects them on untargeted tasks, so the
#     module renders them only when targets exist;
#   - the invocation union renders exactly the one arm the spec set
#     (the spec's CELs guarantee one arm, matching task_type).

resource "aws_ssm_maintenance_window" "this" {
  # metadata.name is the window name on both engines.
  name = var.metadata.name

  schedule = var.spec.schedule
  duration = var.spec.duration
  cutoff   = var.spec.cutoff

  description = var.spec.description != "" ? var.spec.description : null

  enabled = var.spec.enabled

  allow_unassociated_targets = var.spec.allow_unassociated_targets

  schedule_timezone = var.spec.schedule_timezone != "" ? var.spec.schedule_timezone : null
  schedule_offset   = var.spec.schedule_offset != 0 ? var.spec.schedule_offset : null
  start_date        = var.spec.start_date != "" ? var.spec.start_date : null
  end_date          = var.spec.end_date != "" ? var.spec.end_date : null

  tags = local.aws_tags
}

# Folded target registrations, keyed by name (name, description, and
# resource_type force replacement of the registration at the
# provider).
resource "aws_ssm_maintenance_window_target" "this" {
  for_each = { for t in var.spec.targets : t.name => t }

  window_id     = aws_ssm_maintenance_window.this.id
  name          = each.value.name
  resource_type = each.value.resource_type

  dynamic "targets" {
    for_each = each.value.targets
    content {
      key    = targets.value.key
      values = targets.value.values
    }
  }

  description       = each.value.description != "" ? each.value.description : null
  owner_information = each.value.owner_information != "" ? each.value.owner_information : null
}

# Folded tasks, keyed by name.
resource "aws_ssm_maintenance_window_task" "this" {
  for_each = { for t in var.spec.tasks : t.name => t }

  window_id = aws_ssm_maintenance_window.this.id
  name      = each.value.name
  task_type = each.value.task_type
  task_arn  = each.value.task_arn

  description      = each.value.description != "" ? each.value.description : null
  service_role_arn = each.value.service_role_arn != "" ? each.value.service_role_arn : null
  priority         = each.value.priority
  cutoff_behavior  = each.value.cutoff_behavior != "" ? each.value.cutoff_behavior : null

  # Rate controls are only legal on a task WITH targets.
  max_concurrency = length(each.value.targets) > 0 && each.value.max_concurrency != "" ? each.value.max_concurrency : null
  max_errors      = length(each.value.targets) > 0 && each.value.max_errors != "" ? each.value.max_errors : null

  dynamic "targets" {
    for_each = each.value.targets
    content {
      key    = targets.value.key
      values = targets.value.values
    }
  }

  dynamic "task_invocation_parameters" {
    for_each = each.value.invocation != null ? [each.value.invocation] : []
    content {
      dynamic "run_command_parameters" {
        for_each = task_invocation_parameters.value.run_command != null ? [task_invocation_parameters.value.run_command] : []
        content {
          comment              = run_command_parameters.value.comment != "" ? run_command_parameters.value.comment : null
          document_hash        = run_command_parameters.value.document_hash != "" ? run_command_parameters.value.document_hash : null
          document_hash_type   = run_command_parameters.value.document_hash_type != "" ? run_command_parameters.value.document_hash_type : null
          document_version     = run_command_parameters.value.document_version != "" ? run_command_parameters.value.document_version : null
          output_s3_bucket     = run_command_parameters.value.output_s3_bucket != "" ? run_command_parameters.value.output_s3_bucket : null
          output_s3_key_prefix = run_command_parameters.value.output_s3_key_prefix != "" ? run_command_parameters.value.output_s3_key_prefix : null
          service_role_arn     = run_command_parameters.value.service_role_arn != "" ? run_command_parameters.value.service_role_arn : null
          timeout_seconds      = run_command_parameters.value.timeout_seconds != 0 ? run_command_parameters.value.timeout_seconds : null

          dynamic "parameter" {
            for_each = run_command_parameters.value.parameters
            content {
              name   = parameter.value.name
              values = parameter.value.values
            }
          }

          dynamic "cloudwatch_config" {
            for_each = run_command_parameters.value.cloudwatch_config != null ? [run_command_parameters.value.cloudwatch_config] : []
            content {
              cloudwatch_log_group_name = cloudwatch_config.value.cloudwatch_log_group_name != "" ? cloudwatch_config.value.cloudwatch_log_group_name : null
              cloudwatch_output_enabled = cloudwatch_config.value.cloudwatch_output_enabled
            }
          }

          dynamic "notification_config" {
            for_each = run_command_parameters.value.notification_config != null ? [run_command_parameters.value.notification_config] : []
            content {
              notification_arn    = notification_config.value.notification_arn != "" ? notification_config.value.notification_arn : null
              notification_events = length(notification_config.value.notification_events) > 0 ? notification_config.value.notification_events : null
              notification_type   = notification_config.value.notification_type != "" ? notification_config.value.notification_type : null
            }
          }
        }
      }

      dynamic "automation_parameters" {
        for_each = task_invocation_parameters.value.automation != null ? [task_invocation_parameters.value.automation] : []
        content {
          document_version = automation_parameters.value.document_version != "" ? automation_parameters.value.document_version : null

          dynamic "parameter" {
            for_each = automation_parameters.value.parameters
            content {
              name   = parameter.value.name
              values = parameter.value.values
            }
          }
        }
      }

      dynamic "lambda_parameters" {
        for_each = task_invocation_parameters.value.lambda != null ? [task_invocation_parameters.value.lambda] : []
        content {
          client_context = lambda_parameters.value.client_context != "" ? lambda_parameters.value.client_context : null
          payload        = lambda_parameters.value.payload != "" ? lambda_parameters.value.payload : null
          qualifier      = lambda_parameters.value.qualifier != "" ? lambda_parameters.value.qualifier : null
        }
      }

      dynamic "step_functions_parameters" {
        for_each = task_invocation_parameters.value.step_functions != null ? [task_invocation_parameters.value.step_functions] : []
        content {
          input = step_functions_parameters.value.input != "" ? step_functions_parameters.value.input : null
          name  = step_functions_parameters.value.name != "" ? step_functions_parameters.value.name : null
        }
      }
    }
  }
}

# The composite alarm has no metrics, periods, or thresholds of its own — it
# re-evaluates its boolean alarm_rule whenever any referenced alarm changes
# state, which is why this module is a 1:1 mapping with no evaluation
# plumbing.
resource "aws_cloudwatch_composite_alarm" "this" {
  alarm_name = local.resource_name

  # The boolean expression over other alarms' states. Referenced alarms are
  # addressed by NAME (compose from AwsCloudwatchAlarm's exported alarm_name
  # output); the rule text passes through verbatim.
  alarm_rule = var.spec.alarm_rule

  alarm_description = local.alarm_description

  # Tri-state: null = AWS default (true). ForceNew on this resource — a
  # change replaces the alarm (unlike on metric alarms).
  actions_enabled = local.actions_enabled

  alarm_actions             = local.alarm_actions
  ok_actions                = local.ok_actions
  insufficient_data_actions = local.insufficient_data_actions

  # The actions suppressor silences actions (never state transitions) while
  # the designated suppressor alarm is in ALARM — the maintenance-window
  # mechanism. The suppressor is addressed by alarm NAME per the CloudWatch
  # API contract.
  dynamic "actions_suppressor" {
    for_each = var.spec.actions_suppressor != null ? [var.spec.actions_suppressor] : []
    content {
      alarm            = actions_suppressor.value.alarm
      wait_period      = actions_suppressor.value.wait_period
      extension_period = actions_suppressor.value.extension_period
    }
  }

  tags = local.aws_tags
}

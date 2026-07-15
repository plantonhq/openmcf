resource "aws_cloudwatch_metric_alarm" "this" {
  alarm_name = local.resource_name

  # Threshold-mode evaluation contract. Null for PromQL alarms — the query
  # itself expresses the condition, and the provider rejects these arguments
  # in that mode.
  comparison_operator = local.comparison_operator
  evaluation_periods  = local.evaluation_periods
  threshold           = local.threshold
  threshold_metric_id = local.threshold_metric_id
  datapoints_to_alarm = local.datapoints_to_alarm
  treat_missing_data  = local.treat_missing_data

  # Simple metric mode (mutually exclusive with metric_query and
  # evaluation_criteria — CEL-enforced in the spec).
  metric_name        = local.metric_name
  namespace          = local.namespace
  period             = local.period
  statistic          = local.statistic
  extended_statistic = local.extended_statistic
  dimensions         = local.dimensions
  unit               = local.unit

  # Actions. actions_enabled is tri-state: null = AWS default (true); an
  # explicit false suppresses actions without stopping evaluation.
  actions_enabled           = local.actions_enabled
  alarm_actions             = local.alarm_actions
  ok_actions                = local.ok_actions
  insufficient_data_actions = local.insufficient_data_actions

  alarm_description = local.alarm_description

  # Percentile low-sample behavior — only meaningful for percentile stats.
  evaluate_low_sample_count_percentiles = local.evaluate_low_sample_count_percentiles

  # Metric query mode: math expressions, anomaly detection bands, and
  # multi-metric alarms. Each query is either an expression or a raw metric
  # retrieval (never both — CEL-enforced).
  dynamic "metric_query" {
    for_each = var.spec.metric_queries
    content {
      id          = metric_query.value.id
      expression  = metric_query.value.expression != "" ? metric_query.value.expression : null
      label       = metric_query.value.label != "" ? metric_query.value.label : null
      period      = metric_query.value.period != 0 ? metric_query.value.period : null
      return_data = metric_query.value.return_data
      account_id  = metric_query.value.account_id != "" ? metric_query.value.account_id : null

      dynamic "metric" {
        for_each = metric_query.value.metric != null ? [metric_query.value.metric] : []
        content {
          metric_name = metric.value.metric_name
          namespace   = metric.value.namespace != "" ? metric.value.namespace : null
          period      = metric.value.period
          stat        = metric.value.stat
          dimensions  = length(metric.value.dimensions) > 0 ? metric.value.dimensions : null
          unit        = metric.value.unit != "" ? metric.value.unit : null
        }
      }
    }
  }

  # PromQL mode: alarm on a PromQL query against an Amazon Managed Service
  # for Prometheus workspace. pending/recovery periods are tri-states — null
  # lets AWS apply its default firing behavior, while an explicit 0 means
  # "transition immediately".
  dynamic "evaluation_criteria" {
    for_each = local.is_promql ? [var.spec.evaluation_criteria] : []
    content {
      promql_criteria {
        query           = evaluation_criteria.value.promql_criteria.query
        pending_period  = evaluation_criteria.value.promql_criteria.pending_period
        recovery_period = evaluation_criteria.value.promql_criteria.recovery_period
      }
    }
  }

  evaluation_interval = local.evaluation_interval

  tags = local.aws_tags
}

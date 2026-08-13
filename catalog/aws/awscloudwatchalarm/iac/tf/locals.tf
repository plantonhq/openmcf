locals {
  # The alarm's cloud name is the resource's metadata.name — the same basis
  # the Pulumi module uses. The name is also the alarm's identity in composite
  # alarm rule expressions, so it must be stable and predictable.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map -- user labels never merge into cloud tags).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCloudwatchAlarm"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # PromQL alarms carry their alerting contract inside the query; the
  # threshold-evaluation arguments only apply to the other two modes. The spec
  # CELs guarantee the fields are populated exactly when the mode needs them.
  is_promql = var.spec.evaluation_criteria != null

  # Threshold vs. anomaly detection are mutually exclusive: when
  # threshold_metric_id points at an ANOMALY_DETECTION_BAND query, the band is
  # the (dynamic) threshold and the static value must not be sent.
  threshold           = local.is_promql || var.spec.threshold_metric_id != "" ? null : var.spec.threshold
  threshold_metric_id = var.spec.threshold_metric_id != "" ? var.spec.threshold_metric_id : null

  # Null-when-unset scalars so the provider applies its own defaults instead
  # of this module freezing them.
  comparison_operator = local.is_promql ? null : var.spec.comparison_operator
  evaluation_periods  = local.is_promql ? null : var.spec.evaluation_periods
  datapoints_to_alarm = var.spec.datapoints_to_alarm != 0 ? var.spec.datapoints_to_alarm : null
  treat_missing_data  = var.spec.treat_missing_data != "" ? var.spec.treat_missing_data : null

  # actions_enabled is a genuine tri-state (optional bool without a platform
  # default): null lets AWS default to true; an explicit false is a real user
  # choice (suppressing actions during maintenance) that must reach the
  # provider.
  actions_enabled = var.spec.actions_enabled

  # Simple metric mode fields — null when the alarm uses another mode.
  metric_name        = var.spec.metric_name != "" ? var.spec.metric_name : null
  namespace          = var.spec.namespace != "" ? var.spec.namespace : null
  period             = var.spec.period != 0 ? var.spec.period : null
  statistic          = var.spec.statistic != "" ? var.spec.statistic : null
  extended_statistic = var.spec.extended_statistic != "" ? var.spec.extended_statistic : null
  dimensions         = length(var.spec.dimensions) > 0 ? var.spec.dimensions : null
  unit               = var.spec.unit != "" ? var.spec.unit : null

  alarm_description = var.spec.alarm_description != "" ? var.spec.alarm_description : null

  evaluate_low_sample_count_percentiles = var.spec.evaluate_low_sample_count_percentiles != "" ? var.spec.evaluate_low_sample_count_percentiles : null

  # PromQL evaluation cadence — only valid alongside evaluation_criteria.
  evaluation_interval = var.spec.evaluation_interval != 0 ? var.spec.evaluation_interval : null

  # Action ARNs arrive pre-resolved to plain strings by the orchestrator.
  alarm_actions             = var.spec.alarm_actions
  ok_actions                = var.spec.ok_actions
  insufficient_data_actions = var.spec.insufficient_data_actions
}

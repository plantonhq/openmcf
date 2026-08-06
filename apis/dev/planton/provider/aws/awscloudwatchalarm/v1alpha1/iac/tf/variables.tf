variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsCloudwatchAlarm specification"
  type = object({
    region = string
    comparison_operator = optional(string, "")
    evaluation_periods = optional(number, 0)
    datapoints_to_alarm = optional(number, 0)
    threshold = optional(number, 0)
    threshold_metric_id = optional(string, "")
    treat_missing_data = optional(string, "")
    actions_enabled = optional(bool)
    metric_name = optional(string, "")
    namespace = optional(string, "")
    period = optional(number, 0)
    statistic = optional(string, "")
    extended_statistic = optional(string, "")
    dimensions = optional(map(string), {})
    unit = optional(string, "")
    metric_queries = optional(list(object({
      id = string
      expression = optional(string, "")
      metric = optional(object({
        metric_name = string
        namespace = optional(string, "")
        period = optional(number, 0)
        stat = string
        dimensions = optional(map(string), {})
        unit = optional(string, "")
      }))
      label = optional(string, "")
      period = optional(number, 0)
      return_data = optional(bool, false)
      account_id = optional(string, "")
    })), [])
    alarm_actions = optional(list(string), [])
    ok_actions = optional(list(string), [])
    insufficient_data_actions = optional(list(string), [])
    alarm_description = optional(string, "")
    evaluate_low_sample_count_percentiles = optional(string, "")
    evaluation_criteria = optional(object({
      promql_criteria = object({
        query = string
        pending_period = optional(number)
        recovery_period = optional(number)
      })
    }))
    evaluation_interval = optional(number, 0)
  })
}

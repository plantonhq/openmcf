variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name        = string
    id          = optional(string, "")
    org         = optional(string, "")
    env         = optional(string, "")
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})
    tags        = optional(list(string), [])
  })
}

variable "spec" {
  description = "GcpMonitoringAlertPolicy specification"
  type = object({
    project_id   = optional(string, "")
    display_name = optional(string, "")
    combiner     = string
    severity     = optional(string, "")
    enabled      = optional(bool)
    conditions = list(object({
      display_name = string
      condition_threshold = optional(object({
        filter          = string
        comparison      = string
        threshold_value = optional(number, 0)
        duration        = string
        aggregations = optional(list(object({
          alignment_period     = optional(string, "")
          per_series_aligner   = optional(string, "")
          cross_series_reducer = optional(string, "")
          group_by_fields      = optional(list(string), [])
        })), [])
        denominator_filter = optional(string, "")
        denominator_aggregations = optional(list(object({
          alignment_period     = optional(string, "")
          per_series_aligner   = optional(string, "")
          cross_series_reducer = optional(string, "")
          group_by_fields      = optional(list(string), [])
        })), [])
        forecast_options = optional(object({
          forecast_horizon = string
        }))
        trigger = optional(object({
          count   = optional(number, 0)
          percent = optional(number, 0)
        }))
        evaluation_missing_data = optional(string, "")
      }))
      condition_absent = optional(object({
        filter   = string
        duration = string
        aggregations = optional(list(object({
          alignment_period     = optional(string, "")
          per_series_aligner   = optional(string, "")
          cross_series_reducer = optional(string, "")
          group_by_fields      = optional(list(string), [])
        })), [])
        trigger = optional(object({
          count   = optional(number, 0)
          percent = optional(number, 0)
        }))
      }))
      condition_matched_log = optional(object({
        filter           = string
        label_extractors = optional(map(string), {})
      }))
      condition_monitoring_query_language = optional(object({
        query    = string
        duration = string
        trigger = optional(object({
          count   = optional(number, 0)
          percent = optional(number, 0)
        }))
        evaluation_missing_data = optional(string, "")
      }))
      condition_prometheus_query_language = optional(object({
        query                     = string
        duration                  = optional(string, "")
        evaluation_interval       = optional(string, "")
        labels                    = optional(map(string), {})
        rule_group                = optional(string, "")
        alert_rule                = optional(string, "")
        disable_metric_validation = optional(bool, false)
      }))
      condition_sql = optional(object({
        query = string
        minutes = optional(object({
          periodicity = optional(number, 0)
        }))
        hourly = optional(object({
          periodicity   = optional(number, 0)
          minute_offset = optional(number)
        }))
        daily = optional(object({
          periodicity = optional(number, 0)
          execution_time = optional(object({
            hours   = optional(number, 0)
            minutes = optional(number, 0)
            seconds = optional(number, 0)
            nanos   = optional(number, 0)
          }))
        }))
        row_count_test = optional(object({
          comparison = string
          threshold  = optional(number, 0)
        }))
        boolean_test = optional(object({
          column = string
        }))
      }))
    }))
    notification_channels = optional(list(string), [])
    alert_strategy = optional(object({
      auto_close = optional(string, "")
      notification_rate_limit = optional(object({
        period = optional(string, "")
      }))
      notification_channel_strategy = optional(list(object({
        notification_channel_names = optional(list(string), [])
        renotify_interval          = optional(string, "")
      })), [])
      notification_prompts = optional(list(string), [])
    }))
    documentation = optional(object({
      content   = optional(string, "")
      mime_type = optional(string, "")
      subject   = optional(string, "")
      links = optional(list(object({
        display_name = optional(string, "")
        url          = optional(string, "")
      })), [])
    }))
    labels          = optional(map(string), {})
    deletion_policy = optional(string, "")
  })
}
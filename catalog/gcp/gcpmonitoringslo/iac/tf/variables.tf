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
  description = "GcpMonitoringSlo specification"
  type = object({
    project_id = optional(string, "")
    service = object({
      service_id = optional(string, "")
      custom_service = optional(object({
        service_id              = optional(string, "")
        display_name            = optional(string, "")
        telemetry_resource_name = optional(string, "")
      }))
      basic_service = optional(object({
        service_id     = optional(string, "")
        service_type   = optional(string, "")
        service_labels = optional(map(string), {})
      }))
    })
    slo_id              = optional(string, "")
    display_name        = optional(string, "")
    goal                = number
    calendar_period     = optional(string, "")
    rolling_period_days = optional(number, 0)
    sli = object({
      basic_sli = optional(object({
        location = optional(list(string), [])
        method   = optional(list(string), [])
        version  = optional(list(string), [])
        availability = optional(object({
          enabled = optional(bool)
        }))
        latency = optional(object({
          threshold = string
        }))
      }))
      request_based_sli = optional(object({
        distribution_cut = optional(object({
          distribution_filter = string
          range = optional(object({
            min = optional(number)
            max = optional(number)
          }))
        }))
        good_total_ratio = optional(object({
          good_service_filter  = optional(string, "")
          bad_service_filter   = optional(string, "")
          total_service_filter = optional(string, "")
        }))
      }))
      windows_based_sli = optional(object({
        window_period          = optional(string, "")
        good_bad_metric_filter = optional(string, "")
        good_total_ratio_threshold = optional(object({
          threshold = optional(number, 0)
          basic_sli_performance = optional(object({
            location = optional(list(string), [])
            method   = optional(list(string), [])
            version  = optional(list(string), [])
            availability = optional(object({
              enabled = optional(bool)
            }))
            latency = optional(object({
              threshold = string
            }))
          }))
          performance = optional(object({
            distribution_cut = optional(object({
              distribution_filter = string
              range = optional(object({
                min = optional(number)
                max = optional(number)
              }))
            }))
            good_total_ratio = optional(object({
              good_service_filter  = optional(string, "")
              bad_service_filter   = optional(string, "")
              total_service_filter = optional(string, "")
            }))
          }))
        }))
        metric_mean_in_range = optional(object({
          time_series = string
          range = optional(object({
            min = optional(number)
            max = optional(number)
          }))
        }))
        metric_sum_in_range = optional(object({
          time_series = string
          range = optional(object({
            min = optional(number)
            max = optional(number)
          }))
        }))
      }))
    })
    labels          = optional(map(string), {})
    deletion_policy = optional(string, "")
  })
}
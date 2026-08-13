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
  description = "GcpLogMetric specification"
  type = object({
    project_id  = optional(string, "")
    metric_name = optional(string, "")
    filter      = string
    bucket_name = optional(string, "")
    description = optional(string, "")
    disabled    = optional(bool, false)
    metric_descriptor = optional(object({
      metric_kind  = string
      value_type   = string
      unit         = optional(string, "")
      display_name = optional(string, "")
      labels = optional(list(object({
        key         = string
        description = optional(string, "")
        value_type  = optional(string, "")
      })), [])
    }))
    value_extractor  = optional(string, "")
    label_extractors = optional(map(string), {})
    bucket_options = optional(object({
      explicit_buckets = optional(object({
        bounds = list(number)
      }))
      exponential_buckets = optional(object({
        num_finite_buckets = optional(number, 0)
        growth_factor      = optional(number, 0)
        scale              = optional(number, 0)
      }))
      linear_buckets = optional(object({
        num_finite_buckets = optional(number, 0)
        offset             = optional(number, 0)
        width              = optional(number, 0)
      }))
    }))
    deletion_policy = optional(string, "")
  })
}
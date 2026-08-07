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
  description = "AwsCloudwatchLogGroup specification"
  type = object({
    region = string
    retention_in_days = optional(number, 0)
    kms_key_id = optional(string, "")
    log_group_class = optional(string, "")
    deletion_protection_enabled = optional(bool, false)
    metric_filters = optional(list(object({
      name = string
      pattern = optional(string, "")
      apply_on_transformed_logs = optional(bool, false)
      transformation = object({
        metric_name = string
        metric_namespace = string
        metric_value = string
        default_value = optional(number)
        dimensions = optional(map(string), {})
        unit = optional(string, "")
      })
    })), [])
    subscription_filters = optional(list(object({
      name = string
      destination_arn = string
      filter_pattern = optional(string, "")
      role_arn = optional(string, "")
      distribution = optional(string, "")
      emit_system_fields = optional(list(string), [])
      apply_on_transformed_logs = optional(bool, false)
    })), [])
    data_protection_policy = optional(any)
    field_index_policy = optional(any)
  })
}

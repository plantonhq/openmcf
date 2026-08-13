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
  description = "AwsStepFunction specification"
  type = object({
    region = string
    type = optional(string, "")
    definition = any
    role_arn = string
    publish = optional(bool, false)
    aliases = optional(list(object({
      name = optional(string, "")
      description = optional(string, "")
    })), [])
    logging = optional(object({
      level = optional(string, "")
      include_execution_data = optional(bool, false)
      log_destination = optional(string, "")
    }))
    tracing_enabled = optional(bool)
    encryption = optional(object({
      kms_key_id = string
      kms_data_key_reuse_period_seconds = optional(number, 0)
    }))
  })
}
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
  description = "AwsEventBridgeBus specification"
  type = object({
    region = string
    description = optional(string, "")
    kms_key_identifier = optional(string, "")
    event_source_name = optional(string, "")
    dead_letter_config = optional(object({
      arn = string
    }))
    log_config = optional(object({
      level = string
      include_detail = optional(string, "")
    }))
    resource_policy = optional(any)
  })
}

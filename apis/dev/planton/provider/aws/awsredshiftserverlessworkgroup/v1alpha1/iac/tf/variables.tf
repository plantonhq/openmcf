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
  description = "AwsRedshiftServerlessWorkgroup specification"
  type = object({
    region = string
    namespace_name = string
    base_capacity = optional(number, 0)
    max_capacity = optional(number, 0)
    price_performance_target = optional(object({
      enabled = optional(bool, false)
      level = optional(number, 0)
    }))
    subnet_ids = optional(list(string), [])
    security_group_ids = optional(list(string), [])
    enhanced_vpc_routing = optional(bool, false)
    publicly_accessible = optional(bool, false)
    port = optional(number, 0)
    config_parameters = optional(list(object({
      name = string
      value = string
    })), [])
    track_name = optional(string, "")
  })
}

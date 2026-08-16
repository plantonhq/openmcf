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
  description = "AwsGlobalAccelerator specification"
  type = object({
    region = string
    enabled = optional(bool)
    ip_address_type = optional(string)
    ip_addresses = optional(list(string), [])
    flow_logs = optional(object({
      enabled = optional(bool, false)
      s3_bucket = optional(string, "")
      s3_prefix = optional(string, "")
    }))
    listeners = list(object({
      name = string
      protocol = string
      client_affinity = optional(string)
      port_ranges = list(object({
        from_port = number
        to_port = number
      }))
      endpoint_groups = list(object({
        name = string
        endpoint_group_region = optional(string, "")
        health_check_port = optional(number)
        health_check_protocol = optional(string)
        health_check_path = optional(string, "")
        health_check_interval_seconds = optional(number)
        threshold_count = optional(number)
        traffic_dial_percentage = optional(number)
        endpoints = optional(list(object({
          endpoint_id = string
          weight = optional(number)
          client_ip_preservation_enabled = optional(bool)
          attachment_arn = optional(string, "")
        })), [])
        port_overrides = optional(list(object({
          listener_port = number
          endpoint_port = number
        })), [])
      }))
    }))
  })
}

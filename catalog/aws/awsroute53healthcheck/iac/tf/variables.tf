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
  description = "AwsRoute53HealthCheck specification"
  type = object({
    region = string
    check_type = string
    fqdn = optional(string, "")
    ip_address = optional(string, "")
    port = optional(number, 0)
    resource_path = optional(string, "")
    search_string = optional(string, "")
    request_interval = optional(number)
    failure_threshold = optional(number)
    measure_latency = optional(bool, false)
    enable_sni = optional(bool)
    regions = optional(list(string), [])
    invert_healthcheck = optional(bool, false)
    disabled = optional(bool, false)
    child_health_checks = optional(list(string), [])
    child_health_threshold = optional(number)
    cloudwatch_alarm_name = optional(string, "")
    cloudwatch_alarm_region = optional(string, "")
    insufficient_data_health_status = optional(string, "")
    routing_control_arn = optional(string, "")
  })
}

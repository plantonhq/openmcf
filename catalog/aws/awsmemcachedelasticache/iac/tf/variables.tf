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
  description = "AwsMemcachedElasticache specification"
  type = object({
    region = string
    engine_version = optional(string, "")
    node_type = string
    num_cache_nodes = optional(number, 0)
    az_mode = optional(string, "")
    port = optional(number)
    transit_encryption_enabled = optional(bool, false)
    subnet_ids = optional(list(string), [])
    subnet_group_name = optional(string, "")
    security_group_ids = optional(list(string), [])
    network_type = optional(string, "")
    ip_discovery = optional(string, "")
    parameter_group_family = optional(string, "")
    parameters = optional(list(object({
      name = string
      value = string
    })), [])
    parameter_group_name = optional(string, "")
    maintenance_window = optional(string, "")
    apply_immediately = optional(bool, false)
    auto_minor_version_upgrade = optional(bool)
    notification_topic_arn = optional(string, "")
    preferred_availability_zones = optional(list(string), [])
    availability_zone = optional(string, "")
  })
}
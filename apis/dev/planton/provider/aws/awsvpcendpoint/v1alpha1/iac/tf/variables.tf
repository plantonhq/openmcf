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
  description = "AwsVpcEndpoint specification"
  type = object({
    region = string
    vpc_id = string
    endpoint_type = optional(string, "")
    service_name = optional(string, "")
    resource_configuration_arn = optional(string, "")
    service_network_arn = optional(string, "")
    route_table_ids = optional(list(string), [])
    subnet_ids = optional(list(string), [])
    security_group_ids = optional(list(string), [])
    private_dns_enabled = optional(bool, false)
    dns_options = optional(object({
      dns_record_ip_type = optional(string, "")
      private_dns_only_for_inbound_resolver_endpoint = optional(bool, false)
      private_dns_preference = optional(string, "")
      private_dns_specified_domains = optional(list(string), [])
    }))
    ip_address_type = optional(string, "")
    policy = optional(string, "")
    subnet_configurations = optional(list(object({
      subnet_id = string
      ipv4 = optional(string, "")
      ipv6 = optional(string, "")
    })), [])
    service_region = optional(string, "")
    auto_accept = optional(bool, false)
  })
}

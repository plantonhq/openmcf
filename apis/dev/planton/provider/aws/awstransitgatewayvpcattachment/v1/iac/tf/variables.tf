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
  description = "AwsTransitGatewayVpcAttachment specification"
  type = object({
    region = string
    transit_gateway_id = string
    vpc_id = string
    subnet_ids = list(string)
    dns_support = optional(bool)
    ipv6_support = optional(bool, false)
    appliance_mode_support = optional(bool, false)
    security_group_referencing_support = optional(bool)
    default_route_table_association = optional(bool)
    default_route_table_propagation = optional(bool)
  })
}

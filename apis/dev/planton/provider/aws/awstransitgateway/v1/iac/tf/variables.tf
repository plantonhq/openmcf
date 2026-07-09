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
  description = "AwsTransitGateway specification"
  type = object({
    region = string
    description = optional(string, "")
    amazon_side_asn = optional(number, 0)
    default_route_table_association = optional(bool)
    default_route_table_propagation = optional(bool)
    dns_support = optional(bool)
    vpn_ecmp_support = optional(bool)
    auto_accept_shared_attachments = optional(bool, false)
    security_group_referencing_support = optional(bool, false)
    multicast_support = optional(bool, false)
    encryption_support = optional(bool)
    transit_gateway_cidr_blocks = optional(list(string), [])
  })
}

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
  description = "AwsTransitGatewayRouteTable specification"
  type = object({
    region = string
    transit_gateway_id = string
    associations = optional(list(object({
      attachment_id = string
      replace_existing_association = optional(bool, false)
    })), [])
    propagations = optional(list(string), [])
    routes = optional(list(object({
      destination_cidr_block = string
      attachment_id = optional(string, "")
      blackhole = optional(bool, false)
    })), [])
    prefix_list_references = optional(list(object({
      prefix_list_id = string
      attachment_id = optional(string, "")
      blackhole = optional(bool, false)
    })), [])
    set_as_default_association_table = optional(bool, false)
    set_as_default_propagation_table = optional(bool, false)
  })
}

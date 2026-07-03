variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Route Table specification"
  type = object({
    # The Azure region the table is created in; must match the region of
    # the networks whose subnets attach it.
    region = string

    # The resource group the table lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The table's name, unique within the resource group. Renaming replaces
    # the table and detaches it from every subnet.
    name = string

    # The user-defined routes. Each steers traffic bound for its
    # address_prefix (CIDR or Azure service tag) to its next hop.
    # next_hop_type carries the spec enum's name string; VIRTUAL_APPLIANCE
    # routes additionally carry the appliance IP in next_hop_in_ip_address
    # (spec-level validation enforces the pairing).
    routes = optional(list(object({
      name                   = string
      address_prefix         = string
      next_hop_type          = string
      next_hop_in_ip_address = optional(string)
    })), [])

    # Whether BGP-learned routes (ExpressRoute/VPN) propagate into attached
    # subnets. Azure defaults to true; disabling is the forced-tunneling
    # hardening.
    bgp_route_propagation_enabled = optional(bool, true)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

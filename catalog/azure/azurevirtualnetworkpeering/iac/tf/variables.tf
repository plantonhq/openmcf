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
  description = "Azure Virtual Network Peering specification"
  type = object({
    # The peering's name, unique within the local virtual network.
    # Renaming replaces the peering (a brief connectivity gap for this
    # direction).
    name = string

    # The LOCAL virtual network's ARM ID -- the side this peering is
    # written on. The resource group and network name are derived from it.
    # References are resolved to a literal ID by the platform before the
    # module runs.
    virtual_network_id = string

    # The REMOTE virtual network's ARM ID. Works across subscriptions and
    # regions unchanged.
    remote_virtual_network_id = string

    # The four connectivity flags; defaults mirror Azure's (access on;
    # forwarding, gateway transit, and remote gateways off).
    allow_virtual_network_access = optional(bool, true)
    allow_forwarded_traffic      = optional(bool, false)
    allow_gateway_transit        = optional(bool, false)
    use_remote_gateways          = optional(bool, false)

    # Whether the complete address spaces are peered (Azure's default) or
    # only the subnets listed below (subnet-scoped peering).
    peer_complete_virtual_networks_enabled = optional(bool, true)

    # Subnet-scoped peering: the local/remote subnets included, by name.
    # Only meaningful when peer_complete_virtual_networks_enabled = false
    # (spec-level validation enforces the pairing).
    local_subnet_names  = optional(list(string), [])
    remote_subnet_names = optional(list(string), [])

    # Whether only the IPv6 address space is peered (dual-stack networks).
    only_ipv6_peering_enabled = optional(bool, false)
  })
}

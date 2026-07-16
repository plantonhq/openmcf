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
  description = "Azure Network Security Group specification"
  type = object({
    # The Azure region the NSG is created in; must match the region of the
    # subnets and NICs it guards.
    region = string

    # The resource group the NSG lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The NSG's name, unique within the resource group. Renaming replaces
    # the group and detaches it from every subnet and NIC.
    name = string

    # The security rules. direction/access/protocol carry the spec enums'
    # name strings (INBOUND/OUTBOUND, ALLOW/DENY, ANY/TCP/UDP/ICMP/AH/ESP).
    # Ports and addressing each take exactly one form (spec-level
    # validation enforces the pairings); unset source ports and unset
    # addressing mean any ("*").
    security_rules = optional(list(object({
      name        = string
      description = optional(string)
      priority    = number
      direction   = string
      access      = string
      protocol    = string

      source_port_range       = optional(string)
      source_port_ranges      = optional(list(string), [])
      destination_port_range  = optional(string)
      destination_port_ranges = optional(list(string), [])

      source_address_prefix                  = optional(string)
      source_address_prefixes                = optional(list(string), [])
      source_application_security_group_ids  = optional(list(string), [])
      destination_address_prefix             = optional(string)
      destination_address_prefixes           = optional(list(string), [])
      destination_application_security_group_ids = optional(list(string), [])
    })), [])

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

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
  description = "Azure Virtual Network specification"
  type = object({
    # The Azure region the network is created in (a regional resource).
    region = string

    # The resource group the network lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The network's name, unique within the resource group. Renaming
    # replaces the network and everything inside it.
    name = string

    # Self-managed CIDR blocks forming the address space. Exactly one of
    # address_spaces or ip_address_pools is set (spec-level validation).
    address_spaces = optional(list(string), [])

    # Delegated allocation from Azure Network Manager IPAM pools (at most
    # two: one per IP version). The provisioned CIDRs surface in outputs.
    ip_address_pools = optional(list(object({
      id                     = string
      number_of_ip_addresses = string
    })), [])

    # Custom DNS servers; empty means Azure's default resolver
    # (168.63.129.16), which private DNS zone resolution requires.
    dns_servers = optional(list(string), [])

    # BGP community advertised with the network's routes over ExpressRoute,
    # "asn:community" notation (ASN segment is always 12076 today).
    bgp_community = optional(string)

    # Attachment of an existing DDoS Protection Plan (a separate, billed,
    # shared resource). ARM keeps attachment and activation distinct so a
    # plan can stay attached with protection toggled off.
    ddos_protection_plan = optional(object({
      id     = string
      enable = bool
    }))

    # Virtual network encryption enforcement: the spec enum's name string
    # ("ALLOW_UNENCRYPTED"/"DROP_UNENCRYPTED"), or unset for ARM's default
    # (encryption off).
    encryption = optional(string)

    # Connection-tracking flow timeout in minutes (4-30); unset applies
    # Azure's 4-minute default.
    flow_timeout_in_minutes = optional(number)

    # Network-wide private endpoint policy: the spec enum's name string
    # ("BASIC"), or unset for ARM's default ("Disabled").
    private_endpoint_vnet_policies = optional(string)

    # Azure Edge Zone placement; unset deploys to the standard region.
    edge_zone = optional(string)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

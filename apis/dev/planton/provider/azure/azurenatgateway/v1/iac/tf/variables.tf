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
  description = "Azure NAT Gateway specification"
  type = object({
    # The Azure region the gateway is created in; it only serves subnets in
    # its own region.
    region = string

    # The resource group the gateway lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The gateway's name, unique within the resource group. Renaming
    # replaces the gateway, briefly interrupting egress for attached
    # subnets.
    name = string

    # The SKU, as the spec enum's name string (STANDARD/STANDARD_V2).
    # Unset lets Azure apply its default (Standard).
    sku_name = optional(string)

    # How long an idle outbound TCP connection's SNAT port stays reserved,
    # in minutes (Azure defaults to 4).
    idle_timeout_in_minutes = optional(number, 4)

    # The availability zone pinning a STANDARD gateway (STANDARD_V2 is
    # zone-redundant and forbids zones; spec-level validation enforces the
    # pairing).
    zones = optional(list(string), [])

    # The public IPs and prefixes the gateway SNATs through, as resolved
    # ARM IDs. Each drives one association resource.
    public_ip_ids        = optional(list(string), [])
    public_ip_prefix_ids = optional(list(string), [])

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

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
  description = "Azure Public IP Prefix specification"
  type = object({
    # The Azure region the prefix is reserved in; addresses can only be
    # allocated by resources in the same region.
    region = string

    # The resource group the prefix lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The prefix's name, unique within the resource group. Renaming
    # replaces the prefix -- and with it the actual reserved range.
    name = string

    # The CIDR length to reserve (Azure defaults to 28 = 16 addresses).
    prefix_length = optional(number)

    # The IP version, SKU, and tier, as the spec enums' name strings
    # (IPV4/IPV6, STANDARD/STANDARD_V2, REGIONAL/GLOBAL). Unset lets Azure
    # apply its defaults (IPv4 / Standard / Regional).
    ip_version = optional(string)
    sku        = optional(string)
    sku_tier   = optional(string)

    # Availability zones anchoring the range; multiple zones make it
    # zone-redundant.
    zones = optional(list(string), [])

    # The ARM ID of a Custom IP Prefix (bring-your-own range) to carve this
    # prefix out of. Unset allocates from Microsoft's pool.
    custom_ip_prefix_id = optional(string)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision) -- the only thing on a prefix that updates in
    # place.
    tags = optional(map(string), {})
  })
}

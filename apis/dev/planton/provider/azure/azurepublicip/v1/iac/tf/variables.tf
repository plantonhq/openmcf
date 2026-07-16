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
  description = "Azure Public IP specification"
  type = object({
    # The Azure region the address is created in; must match the region of
    # the resource it will attach to.
    region = string

    # The resource group the address lives in. References are resolved to
    # a literal name by the platform before the module runs.
    resource_group = string

    # The address's name, unique within the resource group. Renaming
    # replaces the resource -- and with it the actual address.
    name = string

    # The SKU, tier, and IP version, as the spec enums' name strings
    # (STANDARD/STANDARD_V2, REGIONAL/GLOBAL, IPV4/IPV6). Unset lets Azure
    # apply its defaults (Standard / Regional / IPv4).
    sku        = optional(string)
    sku_tier   = optional(string)
    ip_version = optional(string)

    # Availability zones anchoring the address; multiple zones make it
    # zone-redundant.
    zones = optional(list(string), [])

    # The ARM ID of the public IP prefix to allocate from. References are
    # resolved to a literal ID by the platform before the module runs.
    public_ip_prefix_id = optional(string)

    # The Azure-managed DNS label ({label}.{region}.cloudapp.azure.com) and
    # its scope-based reuse policy (the spec enum's name string, e.g.
    # TENANT_REUSE). Unset scope keeps the classic region-unique behavior.
    domain_name_label       = optional(string)
    domain_name_label_scope = optional(string)

    # The reverse-DNS (PTR) name resolving TO this address; the forward
    # record must exist first.
    reverse_fqdn = optional(string)

    # Idle TCP connection timeout in minutes (Azure defaults to 4).
    idle_timeout_in_minutes = optional(number, 4)

    # Azure IP tags (routing metadata like RoutingPreference), NOT
    # governance tags.
    ip_tags = optional(map(string), {})

    # DDoS stance (the spec enum's name string: DISABLED/ENABLED) and the
    # plan backing ENABLED. Unset mode inherits from the network.
    ddos_protection_mode    = optional(string)
    ddos_protection_plan_id = optional(string)

    # Azure Edge Zone deployment; unset for the standard region.
    edge_zone = optional(string)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

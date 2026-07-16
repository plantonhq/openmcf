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
  description = "Azure Private DNS Zone Virtual Network Link specification"
  type = object({
    # The link's ARM resource name under the parent zone, unique per zone.
    name = string

    # The parent private DNS zone's full ARM ID. References are resolved to
    # a literal ID by the platform before the module runs; the zone's name
    # and resource group are derived from it.
    private_dns_zone_id = string

    # The full ARM ID of the virtual network the zone becomes resolvable
    # from. References are resolved to a literal ID before the module runs.
    virtual_network_id = string

    # Whether Azure auto-registers VM DNS records from the linked network.
    # Azure defaults to false; only ONE registration-enabled link is
    # allowed per network.
    registration_enabled = optional(bool, false)

    # Resolution behavior for names the private zone cannot answer: the
    # spec enum's name string ("DEFAULT"/"NX_DOMAIN_REDIRECT"), or unset to
    # let Azure apply its per-zone-type default.
    resolution_policy = optional(string)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

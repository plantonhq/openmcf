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
  description = "Azure Private DNS Zone specification"
  type = object({
    # The resource group the zone lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The zone's DNS domain name (e.g. "privatelink.postgres.database.azure.com"
    # or "corp.internal"). Renaming replaces the zone and every record in it.
    name = string

    # Optional Start of Authority customization. Unset timers fall back to
    # Azure's defaults; the SOA record is written at creation and cannot be
    # customized afterwards.
    soa_record = optional(object({
      email        = string
      expire_time  = optional(number)
      minimum_ttl  = optional(number)
      refresh_time = optional(number)
      retry_time   = optional(number)
      ttl          = optional(number)
      tags         = optional(map(string), {})
    }))

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

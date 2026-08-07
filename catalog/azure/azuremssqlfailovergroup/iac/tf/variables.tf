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
  description = "Azure SQL Failover Group specification"
  type = object({
    # The failover group name -- also the listener DNS label.
    name = string

    # The primary logical server, resolved to a literal ARM ID by the
    # platform before the module runs.
    server_id = string

    # The partner (secondary) servers, each an object with a resolved
    # server_id ARM ID.
    partner_servers = list(object({
      server_id = string
    }))

    # The databases on the primary to replicate, resolved to literal ARM
    # IDs. Empty is allowed.
    database_ids = optional(list(string), [])

    # The read-write listener failover policy. mode is the spec enum name
    # string (AUTOMATIC / MANUAL); grace_minutes applies only to AUTOMATIC.
    read_write_endpoint_failover_policy = object({
      mode          = string
      grace_minutes = optional(number, 0)
    })

    # Whether the read-only listener also fails over (Azure default true).
    readonly_endpoint_failover_policy_enabled = optional(bool)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

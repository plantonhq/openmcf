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
  description = "Azure User-Assigned Managed Identity specification"
  type = object({
    # The Azure region the identity is created in (a regional resource).
    region = string

    # The resource group the identity lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The identity's name, unique within the resource group. Renaming
    # replaces the identity and mints a new principal.
    name = string

    # Opt-in regional isolation: the spec enum's name string ("REGIONAL"),
    # or unset for ARM's default (usable from any region).
    isolation_scope = optional(string)

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

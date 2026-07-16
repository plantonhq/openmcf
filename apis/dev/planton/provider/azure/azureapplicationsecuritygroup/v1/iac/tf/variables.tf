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
  description = "Azure Application Security Group specification"
  type = object({
    # The Azure region the group lives in; it can only be referenced by
    # network interfaces in the same region.
    region = string

    # The resource group the group is created in. References are resolved to
    # a literal name by the platform before the module runs.
    resource_group = string

    # The group's name, unique within the resource group. Renaming replaces
    # the group (and every rule/NIC that referenced it must be re-pointed).
    name = string

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision) -- the only thing on an ASG that updates in
    # place.
    tags = optional(map(string), {})
  })
}

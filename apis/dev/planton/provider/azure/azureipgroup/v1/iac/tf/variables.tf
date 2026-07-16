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
  description = "Azure IP Group specification"
  type = object({
    # The Azure region the group lives in. IP Groups are regional but can
    # be referenced by firewall policies in any region.
    region = string

    # The resource group the group is created in. References are resolved
    # to a literal name by the platform before the module runs.
    resource_group = string

    # The group's name, unique within the resource group. Renaming replaces
    # the group (and every rule that referenced it must be re-pointed).
    name = string

    # IP addresses and CIDR ranges in the group. An empty group is legal
    # (a placeholder anchor rules can reference before the address plan is
    # final); entries update in place and every referencing rule follows.
    cidrs = optional(list(string), [])

    # Free-form user tags, merged over the metadata-derived tags (user tags
    # win on key collision).
    tags = optional(map(string), {})
  })
}

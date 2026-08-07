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
  description = "Azure RBAC custom Role Definition specification"
  type = object({
    # The role's display name -- unique within the Azure AD tenant.
    name = string

    # The ARM scope the definition is created at (management group,
    # subscription, or resource group ID). References are resolved to a
    # literal ID by the platform before the module runs.
    scope = string

    # Free-text intent shown in the portal's role picker.
    description = optional(string)

    # Permission blocks. Azure evaluates them as a union; one block is the
    # norm. An empty list is legal (a placeholder role granting nothing).
    permissions = optional(list(object({
      actions          = optional(list(string), [])
      not_actions      = optional(list(string), [])
      data_actions     = optional(list(string), [])
      not_data_actions = optional(list(string), [])
    })), [])

    # Scopes at which the role may be assigned. References are resolved to
    # literal IDs before the module runs. Empty means azurerm defaults it to
    # [scope].
    assignable_scopes = optional(list(string), [])

    # Optional pinned GUID for the definition's ARM resource name; Azure
    # generates one when omitted.
    role_definition_id = optional(string)
  })
}

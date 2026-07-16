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
  description = "Azure Cosmos DB SQL role assignment specification"
  type = object({
    # The Cosmos DB account the assignment lives in. References are
    # resolved to a literal ARM ID by the platform before the module
    # runs.
    cosmosdb_account_id = string

    # The fully-scoped resource ID of the role definition to bind --
    # a built-in role's well-known ID or a custom definition's ID.
    # Rebinding is the assignment's one in-place update.
    role_definition_id = string

    # The Entra OBJECT ID of the principal receiving the role.
    principal_id = string

    # The data-plane path the grant applies at: the account's ARM ID,
    # or a database/container path under it.
    scope = string

    # Optional pinned GUID for the assignment's ARM resource name; a
    # random one is generated at deploy time when unset.
    name = optional(string)
  })
}

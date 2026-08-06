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
  description = "Azure Cosmos DB SQL role definition specification"
  type = object({
    # The Cosmos DB account the definition lives in. References are
    # resolved to a literal ARM ID by the platform before the module
    # runs.
    cosmosdb_account_id = string

    # The role's display name -- unique among the account's role
    # definitions; renaming is an in-place update.
    role_name = string

    # The definition's type as the proto enum value name
    # (CUSTOM_ROLE / BUILT_IN_ROLE). Unspecified deploys azurerm's own
    # CustomRole default.
    type = optional(string)

    # The fully-qualified scopes (account, database, or container
    # paths) at or below which this role may be assigned. References
    # are resolved to literal paths before the module runs.
    assignable_scopes = list(string)

    # Additive permission blocks; each carries the Cosmos data actions
    # it allows.
    permissions = list(object({
      data_actions = list(string)
    }))

    # Optional pinned GUID for the definition's ARM resource name; a
    # random one is generated at deploy time when unset.
    role_definition_id = optional(string)
  })
}

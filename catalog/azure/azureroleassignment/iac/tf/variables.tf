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
  description = "Azure RBAC Role Assignment specification"
  type = object({
    # The ARM scope the role applies at (management group, subscription,
    # resource group, or individual resource ID). References are resolved to a
    # literal ID by the platform before the module runs.
    scope = string

    # Exactly one of role_definition_name / role_definition_id is set
    # (enforced by the spec's proto validation before the module runs).
    role_definition_name = optional(string)
    role_definition_id   = optional(string)

    # The Azure AD OBJECT ID of the principal receiving the role
    # (not the application/client ID).
    principal_id = string

    # Enum name from the spec: "SERVICE_PRINCIPAL", "USER", or "GROUP".
    # Omitted when unspecified so Azure infers the type from the directory.
    principal_type = optional(string)

    # Free-text audit note recorded on the assignment.
    description = optional(string)

    # Azure ABAC condition expression and its syntax version ("1.0" or "2.0").
    condition         = optional(string)
    condition_version = optional(string)

    # Cross-tenant (Azure Lighthouse) delegation only.
    delegated_managed_identity_resource_id = optional(string)

    # Skip the Azure AD existence check for freshly created service
    # principals (replication-lag escape hatch).
    skip_service_principal_aad_check = optional(bool, false)

    # Optional pinned GUID for the assignment's ARM resource name; Azure
    # generates one when omitted.
    name = optional(string)
  })
}

# Role assignments carry no ARM tags (Microsoft.Authorization resources do not
# support them), so the usual metadata-derived tag locals are intentionally
# absent from this module.
locals {
  # Map the spec's principal-type enum name onto ARM's PrincipalType strings.
  # null when unspecified so Azure infers the type from the directory object --
  # passing a wrong explicit type is what breaks assignments under
  # ABAC-constrained creators, so we never guess.
  principal_type = (
    var.spec.principal_type == null || var.spec.principal_type == "" ? null :
    var.spec.principal_type == "SERVICE_PRINCIPAL" ? "ServicePrincipal" :
    var.spec.principal_type == "USER" ? "User" :
    var.spec.principal_type == "GROUP" ? "Group" : null
  )

  # Normalize empty strings to null: azurerm validates several of these fields
  # as non-empty WHEN PRESENT, so an empty string must mean "unset", exactly as
  # the Pulumi module omits empty optionals.
  role_definition_name                   = var.spec.role_definition_name == null || var.spec.role_definition_name == "" ? null : var.spec.role_definition_name
  role_definition_id                     = var.spec.role_definition_id == null || var.spec.role_definition_id == "" ? null : var.spec.role_definition_id
  description                            = var.spec.description == null || var.spec.description == "" ? null : var.spec.description
  condition                              = var.spec.condition == null || var.spec.condition == "" ? null : var.spec.condition
  condition_version                      = var.spec.condition_version == null || var.spec.condition_version == "" ? null : var.spec.condition_version
  delegated_managed_identity_resource_id = var.spec.delegated_managed_identity_resource_id == null || var.spec.delegated_managed_identity_resource_id == "" ? null : var.spec.delegated_managed_identity_resource_id
  name                                   = var.spec.name == null || var.spec.name == "" ? null : var.spec.name
}

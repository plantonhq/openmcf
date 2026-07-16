# Create the Azure RBAC role assignment.
#
# Every argument is immutable in Azure (ForceNew): a role assignment is an
# atomic grant record, so any change replaces it (delete + create). That is
# the correct, expected lifecycle -- no ignore_changes games needed.
#
# When condition is set without condition_version, azurerm applies Azure's
# default version ("2.0"); both engines inherit that same defaulting, so the
# deployed result is identical either way.
resource "azurerm_role_assignment" "main" {
  scope        = var.spec.scope
  principal_id = var.spec.principal_id

  # Exactly one of these two is non-null (enforced by the spec's proto
  # validation); azurerm resolves a role name to its definition ID at the
  # target scope.
  role_definition_name = local.role_definition_name
  role_definition_id   = local.role_definition_id

  principal_type                         = local.principal_type
  description                            = local.description
  condition                              = local.condition
  condition_version                      = local.condition_version
  delegated_managed_identity_resource_id = local.delegated_managed_identity_resource_id

  # Skips the Azure AD existence check for freshly created service principals
  # -- new principals replicate asynchronously, and an assignment racing that
  # replication fails with PrincipalNotFound.
  skip_service_principal_aad_check = var.spec.skip_service_principal_aad_check

  # The ARM resource name is a GUID; null means Azure generates one at create
  # time. Pinning it keeps the assignment's full ARM ID stable across
  # replacements.
  name = local.name
}

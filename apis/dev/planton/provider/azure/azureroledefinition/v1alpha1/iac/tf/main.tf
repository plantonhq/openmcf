# Create the custom Azure RBAC role definition.
#
# Lifecycle notes worth knowing before operating this resource:
# - name, description, permissions, and assignable_scopes update IN PLACE;
#   scope and role_definition_id are the definition's ARM identity, so
#   changing either replaces it (delete + create).
# - Updates are eventually consistent: Azure swaps in a new definition record
#   and consolidates it over the following seconds, and azurerm deliberately
#   polls until the created/updated timestamps settle -- expect an update
#   apply to take a few minutes. The deployed permissions are correct as soon
#   as the apply finishes.
# - Deletion also waits for propagation (consecutive not-found polls), so a
#   destroy takes a few minutes as well. Azure refuses to delete a definition
#   that still has role assignments -- destroy the assignments first (in a
#   composed environment the DAG's reverse order does this naturally).
resource "azurerm_role_definition" "main" {
  # The tenant-unique display name and the creation scope. The scope anchors
  # the definition's fully-scoped ARM ID and is the default assignable scope.
  name  = var.spec.name
  scope = var.spec.scope

  description = local.description

  # Permission blocks map 1:1 from the spec. Azure treats multiple blocks as
  # a union of grants; the carve-out lists (not_*) trim only THIS role's
  # grant -- they are not deny rules.
  dynamic "permissions" {
    for_each = var.spec.permissions
    content {
      actions          = permissions.value.actions
      not_actions      = permissions.value.not_actions
      data_actions     = permissions.value.data_actions
      not_data_actions = permissions.value.not_data_actions
    }
  }

  # null lets azurerm default the assignable scopes to [scope] -- the same
  # server-side defaulting the Pulumi engine inherits, so both engines
  # deploy identical definitions for an identical spec.
  assignable_scopes = local.assignable_scopes

  # The ARM resource name is a GUID; null means Azure generates one at
  # create time. Pinning it keeps the definition's full ARM ID stable across
  # replacements.
  role_definition_id = local.role_definition_id
}

# Create the user-assigned managed identity -- a standalone Azure AD identity
# workloads authenticate as, with no credential to store or rotate.
#
# Lifecycle notes worth knowing before operating this resource:
# - tags and isolation_scope update IN PLACE; name, region, and resource
#   group are the identity's ARM identity, so changing any of them replaces
#   it -- which mints a NEW principal and client ID, silently invalidating
#   every existing grant and federated trust rule that referenced the old
#   ones. Composed environments recover automatically (references
#   re-resolve); externally-wired consumers do not.
# - The identity is deliberately just the identity: grants live in
#   AzureRoleAssignment and keyless trust rules in
#   AzureFederatedIdentityCredential, both referencing this identity's
#   outputs.
resource "azurerm_user_assigned_identity" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Omitted (null) means ARM's default: the identity is usable by resources
  # in any region. Only the opt-in "Regional" isolation mode is ever sent.
  isolation_scope = local.isolation_scope

  tags = local.final_tags
}

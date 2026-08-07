# A Cosmos DB SQL (NoSQL) API role definition -- a named bundle of
# DATA-PLANE permissions the account's role assignments bind to Entra
# principals. This is Cosmos DB's own RBAC system, separate from ARM
# RBAC: ARM roles manage the account, these roles govern the data
# inside it. Addressed by the (resource group, account, GUID) trio
# azurerm requires, parsed from the parent account's ARM ID. No Azure
# tags: ARM does not support tags on Cosmos child resources, so the
# platform's identity tags live on the account.
resource "azurerm_cosmosdb_sql_role_definition" "main" {
  name                = var.spec.role_name
  resource_group_name = local.resource_group_name
  account_name        = local.cosmosdb_account_name

  # Pinned GUID identity. Sent only when the spec pins one -- unset
  # lets the provider generate a random GUID at create time (the
  # recommended posture; assignments consume the full resource ID, not
  # the GUID).
  role_definition_id = var.spec.role_definition_id

  # Unset deploys azurerm's own CustomRole default -- the only type
  # organizations author (built-in definitions already exist in every
  # account).
  type = local.role_type

  # WHERE assignments of this role may be created: the account itself
  # or database/container paths under it. Scopes above the account are
  # not enforceable in Cosmos data-plane RBAC. The referenced paths
  # need not exist yet.
  assignable_scopes = var.spec.assignable_scopes

  # WHAT the role allows. Blocks are additive (a union): an operation
  # is permitted if any block's data actions match it. Cosmos supports
  # ALLOW rules only -- no not_data_actions carve-out exists in this
  # RBAC system.
  dynamic "permissions" {
    for_each = var.spec.permissions
    content {
      data_actions = permissions.value.data_actions
    }
  }
}

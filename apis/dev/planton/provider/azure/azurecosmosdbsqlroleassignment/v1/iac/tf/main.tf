# A Cosmos DB SQL (NoSQL) API role assignment -- the grant record
# binding a Cosmos data-plane role to a Microsoft Entra principal at an
# account, database, or container scope. This is Cosmos DB's own RBAC
# system, separate from ARM RBAC: with the account's key authentication
# disabled, these grants are the ONLY way clients connect. Addressed by
# the (resource group, account, GUID) trio azurerm requires, parsed
# from the parent account's ARM ID. No Azure tags: ARM does not support
# tags on Cosmos child resources.
resource "azurerm_cosmosdb_sql_role_assignment" "main" {
  resource_group_name = local.resource_group_name
  account_name        = local.cosmosdb_account_name

  # Pinned GUID identity. Sent only when the spec pins one -- unset
  # lets the provider generate a random GUID at create time (the
  # recommended posture).
  name = var.spec.name

  # WHAT is permitted. The provider validates this is a full
  # sqlRoleDefinitions resource ID at plan time; the definition must
  # live in the same account. Rebinding is the one in-place update --
  # every other field replaces the grant record.
  role_definition_id = var.spec.role_definition_id

  # WHO receives it -- the Entra OBJECT ID (not the client ID; a
  # client ID is accepted by ARM but grants nothing).
  principal_id = var.spec.principal_id

  # WHERE it applies -- permissions inherit downward from the scope.
  # Must sit at or below one of the definition's assignable scopes;
  # Azure enforces that pairing at apply time.
  scope = var.spec.scope
}

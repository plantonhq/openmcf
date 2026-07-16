locals {
  # azurerm addresses Cosmos children by the (resource group, account,
  # name) trio rather than an ARM ID, so both are parsed from the
  # resolved account ID -- the spec models a single parent reference
  # and the module derives the rest (no redundant, contradictable
  # state). The named-group regexes fail the plan loudly if the ID is
  # not a Cosmos DB account ARM ID.
  cosmosdb_account_name = regex("/databaseAccounts/(?P<name>[^/]+)$", var.spec.cosmosdb_account_id)["name"]
  resource_group_name   = regex("/resourceGroups/(?P<name>[^/]+)/", var.spec.cosmosdb_account_id)["name"]
}

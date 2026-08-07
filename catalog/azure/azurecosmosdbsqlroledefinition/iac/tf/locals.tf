locals {
  # azurerm addresses Cosmos RBAC resources by the (resource group,
  # account, GUID) trio rather than an ARM ID, so both names are parsed
  # from the resolved account ID -- the spec models a single parent
  # reference and the module derives the rest (no redundant,
  # contradictable state). The named-group regexes fail the plan loudly
  # if the ID is not a Cosmos DB account ARM ID.
  cosmosdb_account_name = regex("/databaseAccounts/(?P<name>[^/]+)$", var.spec.cosmosdb_account_id)["name"]
  resource_group_name   = regex("/resourceGroups/(?P<name>[^/]+)/", var.spec.cosmosdb_account_id)["name"]

  # The definition's type. Unspecified deploys azurerm's own default
  # (CustomRole -- the only shape organizations author); the enum's
  # proto value names map to ARM's exact wire vocabulary.
  type_map = {
    "CUSTOM_ROLE"   = "CustomRole"
    "BUILT_IN_ROLE" = "BuiltInRole"
  }
  role_type = try(local.type_map[var.spec.type], null)
}

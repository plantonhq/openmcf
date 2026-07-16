locals {
  # azurerm addresses Cosmos children by the (resource group, account,
  # database, name) tuple rather than an ARM ID, so all three parents
  # are parsed from the resolved database ID -- the spec models a single
  # parent reference and the module derives the rest. The named-group
  # regexes fail the plan loudly if the ID is not a Cosmos DB Mongo
  # database ARM ID.
  mongo_database_name   = regex("/mongodbDatabases/(?P<name>[^/]+)$", var.spec.mongo_database_id)["name"]
  cosmosdb_account_name = regex("/databaseAccounts/(?P<name>[^/]+)/", var.spec.mongo_database_id)["name"]
  resource_group_name   = regex("/resourceGroups/(?P<name>[^/]+)/", var.spec.mongo_database_id)["name"]
}

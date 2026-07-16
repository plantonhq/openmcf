locals {
  # azurerm addresses Cosmos children by the (resource group, account,
  # database, name) tuple rather than an ARM ID, so all three parents
  # are parsed from the resolved database ID -- the spec models a single
  # parent reference and the module derives the rest. The named-group
  # regexes fail the plan loudly if the ID is not a Cosmos DB SQL
  # database ARM ID.
  sql_database_name     = regex("/sqlDatabases/(?P<name>[^/]+)$", var.spec.sql_database_id)["name"]
  cosmosdb_account_name = regex("/databaseAccounts/(?P<name>[^/]+)/", var.spec.sql_database_id)["name"]
  resource_group_name   = regex("/resourceGroups/(?P<name>[^/]+)/", var.spec.sql_database_id)["name"]

  # The spec's enums arrive as FULL proto value names; the maps carry
  # the complete vocabularies translated to ARM's wire values.
  partition_key_kind_map = {
    "HASH"       = "Hash"
    "MULTI_HASH" = "MultiHash"
  }

  # Unspecified kind means Hash -- azurerm's own default.
  partition_key_kind = (
    var.spec.partition_key_kind == null || var.spec.partition_key_kind == ""
    ? "Hash"
    : local.partition_key_kind_map[var.spec.partition_key_kind]
  )

  indexing_mode_map = {
    "CONSISTENT" = "consistent"
    "NONE"       = "none"
  }

  composite_index_order_map = {
    "ASCENDING"  = "Ascending"
    "DESCENDING" = "Descending"
  }

  conflict_resolution_mode_map = {
    "LAST_WRITER_WINS" = "LastWriterWins"
    "CUSTOM"           = "Custom"
  }
}

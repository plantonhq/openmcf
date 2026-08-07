output "sql_database_id" {
  description = "The Azure Resource Manager ID of the database (what containers reference)"
  value       = azurerm_cosmosdb_sql_database.main.id
}

output "sql_database_name" {
  description = "The database's name inside the account"
  value       = azurerm_cosmosdb_sql_database.main.name
}

output "cosmosdb_account_name" {
  description = "The name of the Cosmos DB account, parsed from the resolved account ID"
  value       = local.cosmosdb_account_name
}

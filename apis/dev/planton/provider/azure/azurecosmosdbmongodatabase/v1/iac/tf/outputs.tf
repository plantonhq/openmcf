output "mongo_database_id" {
  description = "The Azure Resource Manager ID of the database (what collections reference)"
  value       = azurerm_cosmosdb_mongo_database.main.id
}

output "mongo_database_name" {
  description = "The database's name inside the account"
  value       = azurerm_cosmosdb_mongo_database.main.name
}

output "cosmosdb_account_name" {
  description = "The name of the Cosmos DB account, parsed from the resolved account ID"
  value       = local.cosmosdb_account_name
}

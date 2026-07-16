output "mongo_collection_id" {
  description = "The Azure Resource Manager ID of the collection"
  value       = azurerm_cosmosdb_mongo_collection.main.id
}

output "mongo_collection_name" {
  description = "The collection's name inside the database"
  value       = azurerm_cosmosdb_mongo_collection.main.name
}

output "mongo_database_name" {
  description = "The name of the database the collection lives in, parsed from the resolved database ID"
  value       = local.mongo_database_name
}

output "cosmosdb_account_name" {
  description = "The name of the Cosmos DB account, parsed from the resolved database ID"
  value       = local.cosmosdb_account_name
}

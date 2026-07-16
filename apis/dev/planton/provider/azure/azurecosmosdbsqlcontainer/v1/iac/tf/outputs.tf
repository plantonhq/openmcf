output "sql_container_id" {
  description = "The Azure Resource Manager ID of the container"
  value       = azurerm_cosmosdb_sql_container.main.id
}

output "sql_container_name" {
  description = "The container's name inside the database"
  value       = azurerm_cosmosdb_sql_container.main.name
}

output "sql_database_name" {
  description = "The name of the database the container lives in, parsed from the resolved database ID"
  value       = local.sql_database_name
}

output "cosmosdb_account_name" {
  description = "The name of the Cosmos DB account, parsed from the resolved database ID"
  value       = local.cosmosdb_account_name
}

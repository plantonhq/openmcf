output "elastic_pool_id" {
  description = "The Azure Resource Manager ID of the elastic pool -- the seam AzureMssqlDatabase.elastic_pool_id references"
  value       = azurerm_mssql_elasticpool.main.id
}

output "elastic_pool_name" {
  description = "The name of the elastic pool"
  value       = azurerm_mssql_elasticpool.main.name
}

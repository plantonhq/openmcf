output "database_id" {
  description = "The Azure Resource Manager ID of the database -- the join key for copy/secondary/restore sources and failover-group membership"
  value       = azurerm_mssql_database.main.id
}

output "database_name" {
  description = "The name of the database -- the Database= segment of connection strings against the server's fqdn"
  value       = azurerm_mssql_database.main.name
}

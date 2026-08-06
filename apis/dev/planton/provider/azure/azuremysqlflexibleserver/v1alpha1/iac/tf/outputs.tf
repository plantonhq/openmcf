output "server_id" {
  description = "The Azure Resource Manager ID of the MySQL Flexible Server -- the join key for private endpoints and for replica/restore servers' source_server_id"
  value       = azurerm_mysql_flexible_server.main.id
}

output "server_name" {
  description = "The name of the MySQL Flexible Server"
  value       = azurerm_mysql_flexible_server.main.name
}

output "fqdn" {
  description = "The server's fully qualified domain name ({name}.mysql.database.azure.com); resolves to the private address on a VNet-injected server"
  value       = azurerm_mysql_flexible_server.main.fqdn
}

output "administrator_login" {
  description = "The administrator login, echoed for connection strings"
  value       = azurerm_mysql_flexible_server.main.administrator_login
}

output "database_ids" {
  description = "The ARM ID of each database declared in the spec, keyed by database name"
  value       = { for name, db in azurerm_mysql_flexible_database.main : name => db.id }
}

output "replica_capacity" {
  description = "How many read replicas the server can still accept (Azure computes this from the SKU; burstable SKUs report 0)"
  value       = azurerm_mysql_flexible_server.main.replica_capacity
}

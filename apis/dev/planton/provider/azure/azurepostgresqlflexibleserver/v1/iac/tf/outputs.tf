output "server_id" {
  description = "The Azure Resource Manager ID of the PostgreSQL Flexible Server -- the join key for private endpoints and for replica/restore servers' source_server_id"
  value       = azurerm_postgresql_flexible_server.main.id
}

output "server_name" {
  description = "The name of the PostgreSQL Flexible Server"
  value       = azurerm_postgresql_flexible_server.main.name
}

output "fqdn" {
  description = "The server's fully qualified domain name ({name}.postgres.database.azure.com); resolves to the private address on a VNet-injected server"
  value       = azurerm_postgresql_flexible_server.main.fqdn
}

output "administrator_login" {
  description = "The administrator login, echoed for connection strings; empty on an Entra-only server"
  value       = azurerm_postgresql_flexible_server.main.administrator_login
}

output "database_ids" {
  description = "The ARM ID of each database declared in the spec, keyed by database name"
  value       = { for name, db in azurerm_postgresql_flexible_server_database.main : name => db.id }
}

output "identity_principal_id" {
  description = "The principal ID of the server's system-assigned managed identity -- the subject for role assignments; empty unless the identity type includes SYSTEM_ASSIGNED"
  value       = local.has_system_identity ? azurerm_postgresql_flexible_server.main.identity[0].principal_id : ""
}

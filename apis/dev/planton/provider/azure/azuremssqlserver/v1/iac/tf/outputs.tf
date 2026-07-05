output "server_id" {
  description = "The Azure Resource Manager ID of the logical server -- the join key for AzureMssqlDatabase, AzureMssqlElasticPool, and AzurePrivateEndpoint references"
  value       = azurerm_mssql_server.main.id
}

output "server_name" {
  description = "The name of the logical server"
  value       = azurerm_mssql_server.main.name
}

output "fqdn" {
  description = "The server's fully qualified domain name ({name}.database.windows.net) -- the connection-string host"
  value       = azurerm_mssql_server.main.fully_qualified_domain_name
}

output "administrator_login" {
  description = "The administrator login, echoed for connection strings (empty on an Entra-only server)"
  value       = azurerm_mssql_server.main.administrator_login
}

output "identity_principal_id" {
  description = "The system-assigned identity's principal ID -- the AzureRoleAssignment seam (empty unless the identity type includes SYSTEM_ASSIGNED)"
  value       = try(azurerm_mssql_server.main.identity[0].principal_id, "")
}

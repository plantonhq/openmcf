# Cosmos DB authenticates with account keys rather than login/password
# pairs, so the keys and ready-made connection strings ARE the credential
# surface -- every key and connection string below is marked sensitive.
# When the spec disables local authentication they stop authenticating
# and data-plane access rides Entra ID instead.

output "cosmosdb_account_id" {
  description = "The Azure Resource Manager ID of the account"
  value       = azurerm_cosmosdb_account.main.id
}

output "cosmosdb_account_name" {
  description = "The account's name -- the globally unique DNS label"
  value       = azurerm_cosmosdb_account.main.name
}

output "endpoint" {
  description = "The document endpoint SDKs connect to"
  value       = azurerm_cosmosdb_account.main.endpoint
}

output "read_endpoints" {
  description = "Per-region read endpoints, in failover-priority order"
  value       = azurerm_cosmosdb_account.main.read_endpoints
}

output "write_endpoints" {
  description = "Per-region write endpoints (one entry unless multi-region writes are enabled)"
  value       = azurerm_cosmosdb_account.main.write_endpoints
}

output "primary_key" {
  description = "The primary read-write account key"
  value       = azurerm_cosmosdb_account.main.primary_key
  sensitive   = true
}

output "secondary_key" {
  description = "The secondary read-write account key (the rotation partner)"
  value       = azurerm_cosmosdb_account.main.secondary_key
  sensitive   = true
}

output "primary_readonly_key" {
  description = "The primary read-only account key"
  value       = azurerm_cosmosdb_account.main.primary_readonly_key
  sensitive   = true
}

output "secondary_readonly_key" {
  description = "The secondary read-only account key"
  value       = azurerm_cosmosdb_account.main.secondary_readonly_key
  sensitive   = true
}

output "primary_sql_connection_string" {
  description = "The primary SQL-API connection string"
  value       = azurerm_cosmosdb_account.main.primary_sql_connection_string
  sensitive   = true
}

output "secondary_sql_connection_string" {
  description = "The secondary SQL-API connection string"
  value       = azurerm_cosmosdb_account.main.secondary_sql_connection_string
  sensitive   = true
}

output "primary_readonly_sql_connection_string" {
  description = "The read-only primary SQL-API connection string"
  value       = azurerm_cosmosdb_account.main.primary_readonly_sql_connection_string
  sensitive   = true
}

output "secondary_readonly_sql_connection_string" {
  description = "The read-only secondary SQL-API connection string"
  value       = azurerm_cosmosdb_account.main.secondary_readonly_sql_connection_string
  sensitive   = true
}

output "primary_mongodb_connection_string" {
  description = "The primary MongoDB connection string (meaningful on MONGO_DB accounts)"
  value       = azurerm_cosmosdb_account.main.primary_mongodb_connection_string
  sensitive   = true
}

output "secondary_mongodb_connection_string" {
  description = "The secondary MongoDB connection string"
  value       = azurerm_cosmosdb_account.main.secondary_mongodb_connection_string
  sensitive   = true
}

output "primary_readonly_mongodb_connection_string" {
  description = "The read-only primary MongoDB connection string"
  value       = azurerm_cosmosdb_account.main.primary_readonly_mongodb_connection_string
  sensitive   = true
}

output "secondary_readonly_mongodb_connection_string" {
  description = "The read-only secondary MongoDB connection string"
  value       = azurerm_cosmosdb_account.main.secondary_readonly_mongodb_connection_string
  sensitive   = true
}

output "identity_principal_id" {
  description = "The principal ID of the system-assigned identity, when one is requested"
  value       = try(azurerm_cosmosdb_account.main.identity[0].principal_id, "")
}

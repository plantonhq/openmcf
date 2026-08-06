output "role_definition_id" {
  description = "The fully-scoped ARM ID of the role definition (what a Cosmos DB SQL role assignment binds)"
  value       = azurerm_cosmosdb_sql_role_definition.main.id
}

output "role_definition_guid" {
  description = "The definition's GUID resource name (pinned or generated at deploy time)"
  value       = azurerm_cosmosdb_sql_role_definition.main.role_definition_id
}

output "role_name" {
  description = "The role's display name as deployed"
  value       = azurerm_cosmosdb_sql_role_definition.main.name
}

output "cosmosdb_account_name" {
  description = "The name of the Cosmos DB account, parsed from the resolved account ID"
  value       = local.cosmosdb_account_name
}

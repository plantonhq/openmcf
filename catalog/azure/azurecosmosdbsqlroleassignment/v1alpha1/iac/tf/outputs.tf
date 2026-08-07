output "role_assignment_id" {
  description = "The fully-scoped ARM ID of the role assignment (the grant record's identity)"
  value       = azurerm_cosmosdb_sql_role_assignment.main.id
}

output "role_assignment_guid" {
  description = "The assignment's GUID resource name (pinned or generated at deploy time)"
  value       = azurerm_cosmosdb_sql_role_assignment.main.name
}

output "cosmosdb_account_name" {
  description = "The name of the Cosmos DB account, parsed from the resolved account ID"
  value       = local.cosmosdb_account_name
}

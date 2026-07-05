# The resource_manager_id attribute (rather than the resource id) keeps
# this output byte-identical across engines regardless of which
# addressing path the provider used -- role assignments (Storage Table
# Data Reader/Contributor) scope to it for table-level access.
output "table_id" {
  description = "The Azure Resource Manager ID of the table"
  value       = azurerm_storage_table.main.resource_manager_id
}

output "table_name" {
  description = "The name of the table"
  value       = azurerm_storage_table.main.name
}

# Parsed from the account ARM ID -- consumers frequently need the
# account/table name pair (SDK clients, Functions table bindings).
output "storage_account_name" {
  description = "The name of the storage account the table lives in"
  value       = local.storage_account_name
}

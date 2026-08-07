# The ARM id data-plane role assignments (Storage Blob Data
# Reader/Contributor) scope to for container-level access.
output "container_id" {
  description = "The Azure Resource Manager ID of the container"
  value       = azurerm_storage_container.main.id
}

output "container_name" {
  description = "The name of the container"
  value       = azurerm_storage_container.main.name
}

# Parsed from the account ARM ID -- consumers frequently need the
# account/container name pair (SDK clients, function bindings).
output "storage_account_name" {
  description = "The name of the storage account the container lives in"
  value       = local.storage_account_name
}

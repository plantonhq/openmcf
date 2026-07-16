output "encryption_scope_id" {
  description = "The Azure Resource Manager ID of the encryption scope"
  value       = azurerm_storage_encryption_scope.main.id
}

# What containers (default_encryption_scope), ADLS filesystems, and
# per-blob upload options reference within the account.
output "encryption_scope_name" {
  description = "The name of the encryption scope"
  value       = azurerm_storage_encryption_scope.main.name
}

# Parsed from the account ARM ID -- consumers frequently need the
# account/scope name pair.
output "storage_account_name" {
  description = "The name of the storage account the scope lives in"
  value       = local.storage_account_name
}

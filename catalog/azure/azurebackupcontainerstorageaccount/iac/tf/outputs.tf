output "backup_container_id" {
  description = "The Azure Resource Manager ID of the backup container registration"
  value       = azurerm_backup_container_storage_account.main.id
}

output "storage_account_id" {
  description = "The registered storage account's ARM ID, echoed for protected file shares to wire through (the reference carries both the value and the deploy-order edge)"
  value       = azurerm_backup_container_storage_account.main.storage_account_id
}

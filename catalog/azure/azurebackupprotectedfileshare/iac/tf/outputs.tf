output "backup_protected_file_share_id" {
  description = "The Azure Resource Manager ID of the protected item"
  value       = azurerm_backup_protected_file_share.main.id
}

output "backup_policy_id" {
  description = "The Azure Resource Manager ID of the backup policy"
  value       = azurerm_backup_policy_file_share.main.id
}

output "backup_policy_name" {
  description = "The policy's name -- unique on its vault"
  value       = azurerm_backup_policy_file_share.main.name
}

output "backup_vault_id" {
  description = "The Azure Resource Manager ID of the vault -- what backup policies and backup instances reference their vault by"
  value       = azurerm_data_protection_backup_vault.main.id
}

output "backup_vault_name" {
  description = "The vault's name"
  value       = azurerm_data_protection_backup_vault.main.name
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the vault's system-assigned identity, when one is enabled"
  value       = try(azurerm_data_protection_backup_vault.main.identity[0].principal_id, "")
}

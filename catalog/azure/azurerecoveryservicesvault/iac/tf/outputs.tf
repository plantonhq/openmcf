output "recovery_services_vault_id" {
  description = "The Azure Resource Manager ID of the vault"
  value       = azurerm_recovery_services_vault.main.id
}

output "recovery_services_vault_name" {
  description = "The vault's name -- what backup policies and protected items address their vault by"
  value       = azurerm_recovery_services_vault.main.name
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the vault's system-assigned identity, when one is enabled"
  value       = try(azurerm_recovery_services_vault.main.identity[0].principal_id, "")
}

output "resource_guard_association_id" {
  description = "The ARM ID of the vault's Resource Guard association, when spec.resource_guard_id composes one"
  value       = try(azurerm_recovery_services_vault_resource_guard_association.main[0].id, "")
}

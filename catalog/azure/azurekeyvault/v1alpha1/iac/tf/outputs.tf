output "key_vault_id" {
  description = "The vault's ARM resource ID -- the reference vault keys, certificates, and vault-scoped role assignments target"
  value       = azurerm_key_vault.main.id
}

output "key_vault_name" {
  description = "The vault's name"
  value       = azurerm_key_vault.main.name
}

output "vault_uri" {
  description = "The vault's data-plane URI (https://{name}.vault.azure.net/)"
  value       = azurerm_key_vault.main.vault_uri
}

output "tenant_id" {
  description = "The Azure AD tenant the vault authenticates against"
  value       = azurerm_key_vault.main.tenant_id
}

output "resource_group_name" {
  description = "The resource group the vault was created in"
  value       = azurerm_key_vault.main.resource_group_name
}

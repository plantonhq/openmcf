# The secret's VALUE is deliberately never an output -- consumers read
# it from the vault at runtime via the identifiers below.
output "secret_id" {
  description = "The secret's versioned data-plane ID (pins consumers to this version; value updates do not follow)"
  value       = azurerm_key_vault_secret.main.id
}

output "versionless_id" {
  description = "The secret's versionless data-plane ID -- the reference consumers should use so value updates propagate automatically"
  value       = azurerm_key_vault_secret.main.versionless_id
}

output "secret_name" {
  description = "The secret's name within the vault"
  value       = azurerm_key_vault_secret.main.name
}

output "version" {
  description = "The current version identifier"
  value       = azurerm_key_vault_secret.main.version
}

output "resource_id" {
  description = "The secret's versioned ARM resource ID (control-plane identity)"
  value       = azurerm_key_vault_secret.main.resource_id
}

output "resource_versionless_id" {
  description = "The secret's versionless ARM resource ID"
  value       = azurerm_key_vault_secret.main.resource_versionless_id
}

output "key_id" {
  description = "The key's versioned data-plane ID (pins consumers to this version; rotation does not follow)"
  value       = azurerm_key_vault_key.main.id
}

output "versionless_id" {
  description = "The key's versionless data-plane ID -- the reference CMK consumers should use so rotation propagates automatically"
  value       = azurerm_key_vault_key.main.versionless_id
}

output "key_name" {
  description = "The key's name within the vault"
  value       = azurerm_key_vault_key.main.name
}

output "version" {
  description = "The current version identifier"
  value       = azurerm_key_vault_key.main.version
}

output "resource_id" {
  description = "The key's versioned ARM resource ID (control-plane identity)"
  value       = azurerm_key_vault_key.main.resource_id
}

output "resource_versionless_id" {
  description = "The key's versionless ARM resource ID"
  value       = azurerm_key_vault_key.main.resource_versionless_id
}

output "public_key_pem" {
  description = "The public half of the key in PEM form"
  value       = azurerm_key_vault_key.main.public_key_pem
}

output "public_key_openssh" {
  description = "The public half of the key in OpenSSH form (RSA and P-256/P-384/P-521 EC keys)"
  value       = azurerm_key_vault_key.main.public_key_openssh
}

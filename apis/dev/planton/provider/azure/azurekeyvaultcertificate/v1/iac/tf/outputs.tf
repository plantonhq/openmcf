output "certificate_id" {
  description = "The certificate's versioned data-plane ID (pins consumers to this version; renewals do not follow)"
  value       = azurerm_key_vault_certificate.main.id
}

output "versionless_id" {
  description = "The certificate's versionless data-plane ID -- follows renewals automatically"
  value       = azurerm_key_vault_certificate.main.versionless_id
}

output "secret_id" {
  description = "The versioned ID of the certificate's secret face -- what TLS terminators (Application Gateway, App Service) consume"
  value       = azurerm_key_vault_certificate.main.secret_id
}

output "versionless_secret_id" {
  description = "The versionless ID of the certificate's secret face -- keeps TLS terminators on the current certificate across renewals"
  value       = azurerm_key_vault_certificate.main.versionless_secret_id
}

output "certificate_name" {
  description = "The certificate's name within the vault"
  value       = azurerm_key_vault_certificate.main.name
}

output "version" {
  description = "The current version identifier"
  value       = azurerm_key_vault_certificate.main.version
}

output "thumbprint" {
  description = "The SHA-1 thumbprint of the current certificate, hex-encoded"
  value       = azurerm_key_vault_certificate.main.thumbprint
}

output "resource_manager_id" {
  description = "The certificate's versioned ARM resource ID (control-plane identity)"
  value       = azurerm_key_vault_certificate.main.resource_manager_id
}

output "resource_manager_versionless_id" {
  description = "The certificate's versionless ARM resource ID"
  value       = azurerm_key_vault_certificate.main.resource_manager_versionless_id
}

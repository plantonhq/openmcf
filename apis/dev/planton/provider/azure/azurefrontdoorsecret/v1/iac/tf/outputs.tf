output "secret_id" {
  description = "The Azure Resource Manager ID of the secret (what a custom domain's tls.secret_id references)"
  value       = azurerm_cdn_frontdoor_secret.main.id
}

output "secret_name" {
  description = "The secret's name -- unique within its profile"
  value       = azurerm_cdn_frontdoor_secret.main.name
}

output "subject_alternative_names" {
  description = "The DNS names the wrapped certificate covers (read back from the certificate)"
  value       = azurerm_cdn_frontdoor_secret.main.secret[0].customer_certificate[0].subject_alternative_names
}

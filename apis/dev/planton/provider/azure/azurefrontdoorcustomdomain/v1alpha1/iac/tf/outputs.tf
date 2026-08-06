output "custom_domain_id" {
  description = "The Azure Resource Manager ID of the custom domain (what routes reference in custom_domain_ids and security policies scope WAFs to)"
  value       = azurerm_cdn_frontdoor_custom_domain.main.id
}

output "host_name" {
  description = "The hostname the domain serves -- CNAME this to the endpoint's host_name once validation passes"
  value       = azurerm_cdn_frontdoor_custom_domain.main.host_name
}

output "validation_token" {
  description = "The DNS validation challenge -- publish as a TXT record at _dnsauth.<host_name> to flip the domain from pending to approved"
  value       = azurerm_cdn_frontdoor_custom_domain.main.validation_token
}

output "expiration_date" {
  description = "When the current validation token expires (RFC-3339)"
  value       = azurerm_cdn_frontdoor_custom_domain.main.expiration_date
}

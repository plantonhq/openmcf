output "dns_forwarding_ruleset_id" {
  description = "The Azure Resource Manager ID of the DNS forwarding ruleset -- what virtual network links reference"
  value       = azurerm_private_dns_resolver_dns_forwarding_ruleset.main.id
}

output "dns_forwarding_ruleset_name" {
  description = "The name of the DNS forwarding ruleset resource"
  value       = azurerm_private_dns_resolver_dns_forwarding_ruleset.main.name
}

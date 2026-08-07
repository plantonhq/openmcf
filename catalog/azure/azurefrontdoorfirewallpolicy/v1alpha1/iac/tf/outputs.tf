output "firewall_policy_id" {
  description = "The Azure Resource Manager ID of the WAF policy (what AzureFrontDoorSecurityPolicy references in firewall_policy_id to attach it to a profile's domains)"
  value       = azurerm_cdn_frontdoor_firewall_policy.main.id
}

output "firewall_policy_name" {
  description = "The policy's name -- unique within its resource group"
  value       = azurerm_cdn_frontdoor_firewall_policy.main.name
}

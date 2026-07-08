output "security_policy_id" {
  description = "The Azure Resource Manager ID of the security policy (operational addressing -- nothing composes on the association itself)"
  value       = azurerm_cdn_frontdoor_security_policy.main.id
}

output "security_policy_name" {
  description = "The security policy's name -- unique within its profile"
  value       = azurerm_cdn_frontdoor_security_policy.main.name
}

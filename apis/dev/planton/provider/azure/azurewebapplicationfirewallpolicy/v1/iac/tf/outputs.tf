output "policy_id" {
  description = "The ARM ID of the WAF policy -- the join key Application Gateways attach through (gateway-wide, per listener, or per URL path rule)"
  value       = azurerm_web_application_firewall_policy.main.id
}

output "policy_name" {
  description = "The name of the WAF policy"
  value       = azurerm_web_application_firewall_policy.main.name
}

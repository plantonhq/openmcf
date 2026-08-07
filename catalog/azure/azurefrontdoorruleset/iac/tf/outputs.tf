output "rule_set_id" {
  description = "The Azure Resource Manager ID of the rule set (what routes reference in rule_set_ids to attach this delivery policy)"
  value       = azurerm_cdn_frontdoor_rule_set.main.id
}

output "rule_set_name" {
  description = "The rule set's name -- unique within its profile"
  value       = azurerm_cdn_frontdoor_rule_set.main.name
}

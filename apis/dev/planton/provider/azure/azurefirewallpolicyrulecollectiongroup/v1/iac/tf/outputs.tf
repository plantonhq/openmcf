output "rule_collection_group_id" {
  description = "The Azure Resource Manager ID of the rule collection group"
  value       = azurerm_firewall_policy_rule_collection_group.main.id
}

output "rule_collection_group_name" {
  description = "The name of the rule collection group resource"
  value       = azurerm_firewall_policy_rule_collection_group.main.name
}

output "data_collection_rule_id" {
  description = "The Azure Resource Manager ID of the data collection rule"
  value       = azurerm_monitor_data_collection_rule.main.id
}

output "data_collection_rule_name" {
  description = "The name of the data collection rule resource"
  value       = azurerm_monitor_data_collection_rule.main.name
}

output "immutable_id" {
  description = "The rule's immutable ID -- the identifier agents and the ingestion API address the rule by"
  value       = azurerm_monitor_data_collection_rule.main.immutable_id
}

output "identity_principal_id" {
  description = "The principal ID of the rule's system-assigned identity (empty when no identity is configured)"
  value       = try(azurerm_monitor_data_collection_rule.main.identity[0].principal_id, "")
}

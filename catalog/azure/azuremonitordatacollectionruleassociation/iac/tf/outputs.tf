output "data_collection_rule_association_id" {
  description = "The Azure Resource Manager ID of the association (scoped under the target machine)"
  value       = azurerm_monitor_data_collection_rule_association.main.id
}

output "data_collection_rule_association_name" {
  description = "The association's name on the target machine"
  value       = azurerm_monitor_data_collection_rule_association.main.name
}

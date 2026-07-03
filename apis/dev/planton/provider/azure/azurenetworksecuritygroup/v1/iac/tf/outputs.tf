output "network_security_group_id" {
  description = "The Azure Resource Manager ID of the Network Security Group"
  value       = azurerm_network_security_group.main.id
}

output "network_security_group_name" {
  description = "The name of the Network Security Group"
  value       = azurerm_network_security_group.main.name
}

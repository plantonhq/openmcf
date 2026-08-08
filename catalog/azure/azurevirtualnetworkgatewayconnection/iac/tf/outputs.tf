output "connection_id" {
  description = "The Azure Resource Manager ID of the connection"
  value       = azurerm_virtual_network_gateway_connection.main.id
}

output "connection_name" {
  description = "The name of the connection resource"
  value       = azurerm_virtual_network_gateway_connection.main.name
}

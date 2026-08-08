output "local_network_gateway_id" {
  description = "The Azure Resource Manager ID of the local network gateway"
  value       = azurerm_local_network_gateway.main.id
}

output "local_network_gateway_name" {
  description = "The name of the local network gateway resource"
  value       = azurerm_local_network_gateway.main.name
}

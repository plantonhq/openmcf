output "connection_id" {
  description = "The Azure Resource Manager ID of the connection (a child of the gateway)"
  value       = azurerm_vpn_gateway_connection.main.id
}

output "connection_name" {
  description = "The name of the connection"
  value       = azurerm_vpn_gateway_connection.main.name
}

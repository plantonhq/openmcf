output "express_route_gateway_id" {
  description = "The Azure Resource Manager ID of the ExpressRoute Gateway"
  value       = azurerm_express_route_gateway.main.id
}

output "express_route_gateway_name" {
  description = "The name of the gateway"
  value       = azurerm_express_route_gateway.main.name
}

output "connection_ids" {
  description = "The ARM ID of each connection on the gateway, keyed by the connection's name from the spec"
  value       = { for name, connection in azurerm_express_route_connection.connections : name => connection.id }
}

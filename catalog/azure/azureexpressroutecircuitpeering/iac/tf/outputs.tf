output "express_route_circuit_peering_id" {
  description = "The Azure Resource Manager ID of the peering -- the far side of another circuit's Global Reach connection"
  value       = azurerm_express_route_circuit_peering.main.id
}

output "azure_asn" {
  description = "Microsoft's Autonomous System Number on this peering -- the BGP neighbor ASN for your routers"
  value       = azurerm_express_route_circuit_peering.main.azure_asn
}

output "primary_azure_port" {
  description = "The Microsoft-edge identifier of the primary physical port"
  value       = azurerm_express_route_circuit_peering.main.primary_azure_port
}

output "secondary_azure_port" {
  description = "The Microsoft-edge identifier of the secondary physical port"
  value       = azurerm_express_route_circuit_peering.main.secondary_azure_port
}

output "connection_ids" {
  description = "The ARM ID of each Global Reach connection, keyed by the connection's name"
  value       = { for name, connection in azurerm_express_route_circuit_connection.connections : name => connection.id }
}

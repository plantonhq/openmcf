output "virtual_hub_connection_id" {
  description = "The Azure Resource Manager ID of the connection -- what a hub BGP peering references as its virtual_network_connection_id"
  value       = azurerm_virtual_hub_connection.main.id
}

output "virtual_hub_connection_name" {
  description = "The name of the connection"
  value       = azurerm_virtual_hub_connection.main.name
}

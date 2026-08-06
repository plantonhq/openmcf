output "private_endpoint_id" {
  description = "The Azure Resource Manager ID of the private endpoint"
  value       = azurerm_private_endpoint.endpoint.id
}

output "private_endpoint_name" {
  description = "The name of the private endpoint resource"
  value       = azurerm_private_endpoint.endpoint.name
}

output "private_ip_address" {
  description = "The private IP address allocated to the endpoint from the subnet"
  value       = azurerm_private_endpoint.endpoint.private_service_connection[0].private_ip_address
}

output "network_interface_id" {
  description = "The Azure Resource Manager ID of the auto-created network interface"
  value       = azurerm_private_endpoint.endpoint.network_interface[0].id
}

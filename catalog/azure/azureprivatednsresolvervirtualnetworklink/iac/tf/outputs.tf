output "virtual_network_link_id" {
  description = "The Azure Resource Manager ID of the virtual network link"
  value       = azurerm_private_dns_resolver_virtual_network_link.main.id
}

output "virtual_network_link_name" {
  description = "The name of the virtual network link resource"
  value       = azurerm_private_dns_resolver_virtual_network_link.main.name
}

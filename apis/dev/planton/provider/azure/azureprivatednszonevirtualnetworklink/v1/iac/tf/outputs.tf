output "link_id" {
  description = "The Azure Resource Manager ID of the virtual network link"
  value       = azurerm_private_dns_zone_virtual_network_link.main.id
}

output "link_name" {
  description = "The name of the virtual network link"
  value       = azurerm_private_dns_zone_virtual_network_link.main.name
}

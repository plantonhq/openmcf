output "network_interface_id" {
  description = "The Azure Resource Manager ID of the NIC"
  value       = azurerm_network_interface.main.id
}

output "network_interface_name" {
  description = "The name of the NIC"
  value       = azurerm_network_interface.main.name
}

output "private_ip_address" {
  description = "The primary configuration's private IP address"
  value       = azurerm_network_interface.main.private_ip_address
}

output "private_ip_addresses" {
  description = "The private IP addresses of all configurations, in configuration order"
  value       = azurerm_network_interface.main.private_ip_addresses
}

output "mac_address" {
  description = "The NIC's MAC address (populated once attached to a running VM)"
  value       = azurerm_network_interface.main.mac_address
}

output "internal_domain_name_suffix" {
  description = "The DNS suffix completing internal_dns_name_label into a resolvable FQDN"
  value       = azurerm_network_interface.main.internal_domain_name_suffix
}

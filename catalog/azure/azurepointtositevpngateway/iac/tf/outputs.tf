output "point_to_site_vpn_gateway_id" {
  description = "The Azure Resource Manager ID of the point-to-site VPN gateway"
  value       = azurerm_point_to_site_vpn_gateway.main.id
}

output "point_to_site_vpn_gateway_name" {
  description = "The name of the point-to-site VPN gateway"
  value       = azurerm_point_to_site_vpn_gateway.main.name
}

output "vpn_site_id" {
  description = "The Azure Resource Manager ID of the VPN site -- what a connection references as its remote_vpn_site_id"
  value       = azurerm_vpn_site.main.id
}

output "vpn_site_name" {
  description = "The name of the VPN site"
  value       = azurerm_vpn_site.main.name
}

output "link_ids" {
  description = "The ARM ID of each link on the site, keyed by the link's name from the spec -- what a connection's vpn_links reference as their vpn_site_link_id"
  value       = { for link in azurerm_vpn_site.main.link : link.name => link.id }
}

output "traffic_manager_profile_id" {
  description = "The Azure Resource Manager ID of the Traffic Manager profile -- what endpoints reference, and what alias DNS records target"
  value       = azurerm_traffic_manager_profile.main.id
}

output "traffic_manager_profile_name" {
  description = "The name of the Traffic Manager profile resource"
  value       = azurerm_traffic_manager_profile.main.name
}

output "fqdn" {
  description = "The profile's public DNS name ({relative_name}.trafficmanager.net) -- what users resolve, and what your own domain CNAMEs to"
  value       = azurerm_traffic_manager_profile.main.fqdn
}

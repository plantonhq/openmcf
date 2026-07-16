output "zone_id" {
  description = "The Azure Resource Manager ID of the Private DNS Zone"
  value       = azurerm_private_dns_zone.main.id
}

output "zone_name" {
  description = "The name of the Private DNS Zone"
  value       = azurerm_private_dns_zone.main.name
}

output "resource_group_name" {
  description = "The resource group the zone lives in"
  value       = azurerm_private_dns_zone.main.resource_group_name
}

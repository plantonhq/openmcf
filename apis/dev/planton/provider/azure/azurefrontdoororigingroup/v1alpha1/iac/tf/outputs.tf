output "origin_group_id" {
  description = "The Azure Resource Manager ID of the origin group (what origins reference as parent and routes reference as destination)"
  value       = azurerm_cdn_frontdoor_origin_group.main.id
}

output "origin_group_name" {
  description = "The origin group's name -- unique within its profile"
  value       = azurerm_cdn_frontdoor_origin_group.main.name
}

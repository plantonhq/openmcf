output "origin_id" {
  description = "The Azure Resource Manager ID of the origin (what routes list in origin_ids to sequence deployment)"
  value       = azurerm_cdn_frontdoor_origin.main.id
}

output "origin_name" {
  description = "The origin's name -- unique within its origin group"
  value       = azurerm_cdn_frontdoor_origin.main.name
}

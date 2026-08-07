output "route_id" {
  description = "The Azure Resource Manager ID of the route"
  value       = azurerm_cdn_frontdoor_route.main.id
}

output "route_name" {
  description = "The route's name -- unique within its endpoint"
  value       = azurerm_cdn_frontdoor_route.main.name
}

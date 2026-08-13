output "private_link_service_id" {
  description = "The Azure Resource Manager ID of the Private Link Service"
  value       = azurerm_private_link_service.main.id
}

output "private_link_service_name" {
  description = "The name of the Private Link Service resource"
  value       = azurerm_private_link_service.main.name
}

output "alias" {
  description = "The globally unique alias consumers use to request a private-endpoint connection"
  value       = azurerm_private_link_service.main.alias
}

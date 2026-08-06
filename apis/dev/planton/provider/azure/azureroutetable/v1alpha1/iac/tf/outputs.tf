output "route_table_id" {
  description = "The Azure Resource Manager ID of the route table"
  value       = azurerm_route_table.main.id
}

output "route_table_name" {
  description = "The name of the route table"
  value       = azurerm_route_table.main.name
}

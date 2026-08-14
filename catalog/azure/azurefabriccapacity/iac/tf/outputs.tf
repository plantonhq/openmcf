output "fabric_capacity_id" {
  description = "The Azure Resource Manager ID of the Fabric capacity"
  value       = azurerm_fabric_capacity.main.id
}

output "fabric_capacity_name" {
  description = "The capacity's name -- what Fabric workspaces assign themselves to"
  value       = azurerm_fabric_capacity.main.name
}

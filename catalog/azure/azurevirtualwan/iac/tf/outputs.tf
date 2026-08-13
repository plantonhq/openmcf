output "virtual_wan_id" {
  description = "The Azure Resource Manager ID of the Virtual WAN -- what virtual hubs reference as their virtual_wan_id"
  value       = azurerm_virtual_wan.main.id
}

output "virtual_wan_name" {
  description = "The name of the Virtual WAN"
  value       = azurerm_virtual_wan.main.name
}

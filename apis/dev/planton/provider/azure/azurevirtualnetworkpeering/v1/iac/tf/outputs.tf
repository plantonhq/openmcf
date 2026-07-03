output "peering_id" {
  description = "The Azure Resource Manager ID of the peering"
  value       = azurerm_virtual_network_peering.main.id
}

output "peering_name" {
  description = "The name of the peering within its local virtual network"
  value       = azurerm_virtual_network_peering.main.name
}

output "virtual_network_name" {
  description = "The name of the LOCAL virtual network the peering is written on, derived from the referenced network's ARM ID"
  value       = local.virtual_network_name
}

output "resource_group_name" {
  description = "The name of the resource group the local network lives in, derived from the referenced network's ARM ID"
  value       = local.resource_group_name
}

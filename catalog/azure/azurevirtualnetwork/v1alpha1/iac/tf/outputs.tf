output "virtual_network_id" {
  description = "The Azure Resource Manager ID of the virtual network"
  value       = azurerm_virtual_network.main.id
}

output "virtual_network_name" {
  description = "The name of the virtual network"
  value       = azurerm_virtual_network.main.name
}

output "guid" {
  description = "The stable GUID ARM assigns the virtual network at creation"
  value       = azurerm_virtual_network.main.guid
}

# Reflects the ACTUAL ranges the network carries: echoes the spec when
# self-managed, and surfaces the IPAM-provisioned CIDRs when
# ip_address_pools delegate allocation -- the only place those ranges are
# visible for downstream planning.
output "address_spaces" {
  description = "The address space actually carried by the network"
  value       = azurerm_virtual_network.main.address_space
}

output "subnet_id" {
  description = "The Azure Resource Manager ID of the subnet"
  value       = azurerm_subnet.main.id
}

output "subnet_name" {
  description = "The name of the subnet within its virtual network"
  value       = azurerm_subnet.main.name
}

# Read from the resource (not the spec) so IPAM-allocated subnets surface
# the ranges the Network Manager pool actually provisioned.
output "address_prefixes" {
  description = "The CIDR blocks actually assigned to the subnet"
  value       = azurerm_subnet.main.address_prefixes
}

output "virtual_network_name" {
  description = "The name of the parent virtual network, derived from the referenced network's ARM ID"
  value       = local.virtual_network_name
}

output "resource_group_name" {
  description = "The name of the resource group the subnet lives in, derived from the referenced network's ARM ID"
  value       = local.resource_group_name
}

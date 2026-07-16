output "nat_gateway_id" {
  description = "The Azure Resource Manager ID of the NAT gateway"
  value       = azurerm_nat_gateway.main.id
}

output "nat_gateway_name" {
  description = "The name of the NAT gateway"
  value       = azurerm_nat_gateway.main.name
}

output "resource_guid" {
  description = "The immutable GUID ARM assigns the gateway"
  value       = azurerm_nat_gateway.main.resource_guid
}

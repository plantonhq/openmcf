output "virtual_network_gateway_id" {
  description = "The Azure Resource Manager ID of the gateway"
  value       = azurerm_virtual_network_gateway.main.id
}

output "virtual_network_gateway_name" {
  description = "The name of the gateway resource"
  value       = azurerm_virtual_network_gateway.main.name
}

# The composed NAT rules' ARM ids, keyed by rule name -- connections opt
# into rules via their egress/ingress NAT rule id lists. (The gateway's
# public address is not an output here: it belongs to the referenced
# AzurePublicIp resource and surfaces through that kind's outputs.)
output "nat_rule_ids" {
  description = "The ARM ids of the gateway's NAT rules, keyed by rule name"
  value       = { for name, rule in azurerm_virtual_network_gateway_nat_rule.rules : name => rule.id }
}

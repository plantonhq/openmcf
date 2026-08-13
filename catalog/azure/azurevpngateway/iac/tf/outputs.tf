output "vpn_gateway_id" {
  description = "The Azure Resource Manager ID of the gateway -- what a connection references as its vpn_gateway_id"
  value       = azurerm_vpn_gateway.main.id
}

output "vpn_gateway_name" {
  description = "The name of the gateway"
  value       = azurerm_vpn_gateway.main.name
}

output "bgp_asn" {
  description = "The gateway's BGP autonomous system number -- what branch devices configure as the remote ASN"
  value       = try(azurerm_vpn_gateway.main.bgp_settings[0].asn, 0)
}

output "public_ip_addresses" {
  description = "The PUBLIC IPv4 address of each gateway instance -- what branch devices dial as their tunnel peer"
  value       = [for ip_configuration in azurerm_vpn_gateway.main.ip_configuration : ip_configuration.public_ip_address]
}

output "private_ip_addresses" {
  description = "The private IPv4 address of each gateway instance -- the tunnel endpoints for connections using local_azure_ip_address_enabled"
  value       = [for ip_configuration in azurerm_vpn_gateway.main.ip_configuration : ip_configuration.private_ip_address]
}

output "nat_rule_ids" {
  description = "The ARM ID of each NAT rule on the gateway, keyed by the rule's name from the spec"
  value       = { for name, nat_rule in azurerm_vpn_gateway_nat_rule.nat_rules : name => nat_rule.id }
}

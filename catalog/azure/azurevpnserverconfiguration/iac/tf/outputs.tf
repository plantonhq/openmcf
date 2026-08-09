output "vpn_server_configuration_id" {
  description = "The Azure Resource Manager ID of the VPN server configuration -- what a point-to-site VPN gateway references as its vpn_server_configuration_id"
  value       = azurerm_vpn_server_configuration.main.id
}

output "vpn_server_configuration_name" {
  description = "The name of the VPN server configuration"
  value       = azurerm_vpn_server_configuration.main.name
}

output "policy_group_ids" {
  description = "The ARM ID of each policy group on the configuration, keyed by the group's name from the spec"
  value       = { for name, policy_group in azurerm_vpn_server_configuration_policy_group.policy_groups : name => policy_group.id }
}

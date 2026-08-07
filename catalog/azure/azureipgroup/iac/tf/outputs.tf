# The ARM ID is the composition seam: firewall policy rules
# (source_ip_groups / destination_ip_groups) and intrusion-detection
# traffic bypasses reference it to target the group's address set.
output "ip_group_id" {
  description = "The Azure Resource Manager ID of the IP Group"
  value       = azurerm_ip_group.main.id
}

output "ip_group_name" {
  description = "The name of the IP Group resource"
  value       = azurerm_ip_group.main.name
}

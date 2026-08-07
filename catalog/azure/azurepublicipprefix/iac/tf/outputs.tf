output "public_ip_prefix_id" {
  description = "The Azure Resource Manager ID of the prefix"
  value       = azurerm_public_ip_prefix.main.id
}

# The actual reserved CIDR -- known only after creation, and the value
# partners and firewalls allowlist.
output "ip_prefix" {
  description = "The reserved CIDR range, e.g. \"20.42.0.16/28\""
  value       = azurerm_public_ip_prefix.main.ip_prefix
}

output "public_ip_prefix_name" {
  description = "The name of the prefix resource"
  value       = azurerm_public_ip_prefix.main.name
}

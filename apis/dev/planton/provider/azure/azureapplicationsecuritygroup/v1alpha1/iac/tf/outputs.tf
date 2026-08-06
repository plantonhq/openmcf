# The ARM ID is the composition seam: network interfaces, scale-set IP
# configurations, and NSG security rules reference it to declare membership
# or target the group in a rule.
output "application_security_group_id" {
  description = "The Azure Resource Manager ID of the application security group"
  value       = azurerm_application_security_group.main.id
}

output "application_security_group_name" {
  description = "The name of the application security group resource"
  value       = azurerm_application_security_group.main.name
}

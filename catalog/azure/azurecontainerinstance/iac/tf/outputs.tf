output "container_group_id" {
  description = "The Azure Resource Manager ID of the container group"
  value       = azurerm_container_group.main.id
}

output "container_group_name" {
  description = "The container group's name"
  value       = azurerm_container_group.main.name
}

output "ip_address" {
  description = "The group's IP address -- public or private per ip_address_type; empty for the \"None\" posture"
  value       = azurerm_container_group.main.ip_address
}

output "fqdn" {
  description = "The group's FQDN ({dns_name_label}.{region}.azurecontainer.io); empty unless dns_name_label is set on a public group"
  value       = azurerm_container_group.main.fqdn
}

output "identity_principal_id" {
  description = "The principal ID of the group's system-assigned managed identity; empty unless SYSTEM_ASSIGNED is enabled"
  value       = try(azurerm_container_group.main.identity[0].principal_id, "")
}

output "identity_tenant_id" {
  description = "The tenant ID of the group's system-assigned managed identity; empty unless SYSTEM_ASSIGNED is enabled"
  value       = try(azurerm_container_group.main.identity[0].tenant_id, "")
}

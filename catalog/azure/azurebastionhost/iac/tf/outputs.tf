output "bastion_host_id" {
  description = "The Azure Resource Manager ID of the Bastion host"
  value       = azurerm_bastion_host.main.id
}

output "bastion_host_name" {
  description = "The name of the Bastion host resource"
  value       = azurerm_bastion_host.main.name
}

output "dns_name" {
  description = "The DNS name sessions connect through"
  value       = azurerm_bastion_host.main.dns_name
}

output "private_only_enabled" {
  description = "Whether the host deployed private-only (a Premium host without a public IP)"
  value       = azurerm_bastion_host.main.private_only_enabled
}

output "environment_id" {
  description = "The Azure Resource Manager ID of the Container App Environment"
  value       = azurerm_container_app_environment.main.id
}

output "environment_name" {
  description = "The name of the Container App Environment"
  value       = azurerm_container_app_environment.main.name
}

output "default_domain" {
  description = "The default publicly resolvable domain for apps in this environment"
  value       = azurerm_container_app_environment.main.default_domain
}

output "static_ip_address" {
  description = "The static IP address of the environment"
  value       = azurerm_container_app_environment.main.static_ip_address
}

output "platform_reserved_cidr" {
  description = "The IP range reserved for environment infrastructure"
  value       = azurerm_container_app_environment.main.platform_reserved_cidr
}

output "platform_reserved_dns_ip_address" {
  description = "The IP address reserved for the internal DNS server"
  value       = azurerm_container_app_environment.main.platform_reserved_dns_ip_address
}

output "docker_bridge_cidr" {
  description = "The Docker bridge network address used inside the environment's infrastructure"
  value       = azurerm_container_app_environment.main.docker_bridge_cidr
}

output "custom_domain_verification_id" {
  description = "The TXT-record value proving ownership of a custom DNS suffix or per-app custom domain"
  value       = azurerm_container_app_environment.main.custom_domain_verification_id
}

output "identity_principal_id" {
  description = "The principal ID of the environment's system-assigned managed identity (empty unless SYSTEM_ASSIGNED is enabled)"
  value       = try(azurerm_container_app_environment.main.identity[0].principal_id, "")
}

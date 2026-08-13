output "domain_id" {
  description = "The Azure Resource Manager ID of the Event Grid domain"
  value       = azurerm_eventgrid_domain.main.id
}

output "domain_name" {
  description = "The domain's name (the first label of its endpoint hostname)"
  value       = azurerm_eventgrid_domain.main.name
}

output "endpoint" {
  description = "The HTTPS endpoint publishers POST events to (one endpoint for every topic in the domain)"
  value       = azurerm_eventgrid_domain.main.endpoint
}

output "primary_access_key" {
  description = "The primary access key (the aeg-sas-key header value)"
  value       = azurerm_eventgrid_domain.main.primary_access_key
  sensitive   = true
}

output "secondary_access_key" {
  description = "The secondary access key (the rotation partner)"
  value       = azurerm_eventgrid_domain.main.secondary_access_key
  sensitive   = true
}

output "identity_principal_id" {
  description = "The principal ID of the domain's system-assigned identity (empty when no identity is configured)"
  value       = try(azurerm_eventgrid_domain.main.identity[0].principal_id, "")
}

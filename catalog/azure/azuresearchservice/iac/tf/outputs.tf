output "search_service_id" {
  description = "The Azure Resource Manager ID of the service -- what shared private links and diagnostic settings reference"
  value       = azurerm_search_service.main.id
}

output "search_service_name" {
  description = "The name of the service. ARM addresses shared private links as children of this name, and the data-plane endpoint embeds it"
  value       = azurerm_search_service.main.name
}

output "endpoint" {
  description = "The service's data-plane endpoint (https://{name}.search.windows.net) -- what applications and SDKs call"
  value       = azurerm_search_service.main.endpoint
}

output "primary_key" {
  description = "The primary admin API key -- full read-write control of the service's data plane. Empty when local authentication is disabled"
  value       = azurerm_search_service.main.primary_key
  sensitive   = true
}

output "secondary_key" {
  description = "The secondary admin API key -- the rotation partner of primary_key. Empty when local authentication is disabled"
  value       = azurerm_search_service.main.secondary_key
  sensitive   = true
}

output "default_query_key" {
  description = "The service's built-in query API key -- read-only data-plane access for client applications. Empty when local authentication is disabled"
  value       = try(azurerm_search_service.main.query_keys[0].key, "")
  sensitive   = true
}

output "customer_managed_key_encryption_compliance_status" {
  description = "Whether the service's objects comply with the customer-managed-key enforcement posture (Compliant / NonCompliant)"
  value       = azurerm_search_service.main.customer_managed_key_encryption_compliance_status
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the service's system-assigned identity, when one is enabled"
  value       = try(azurerm_search_service.main.identity[0].principal_id, "")
}

output "shared_private_link_service_ids" {
  description = "The ARM ID of each shared private link on the service, keyed by the link's name from the spec"
  value       = { for name, link in azurerm_search_shared_private_link_service.links : name => link.id }
}

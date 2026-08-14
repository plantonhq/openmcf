output "data_factory_id" {
  description = "The Azure Resource Manager ID of the Data Factory"
  value       = azurerm_data_factory.main.id
}

output "data_factory_name" {
  description = "The factory's name"
  value       = azurerm_data_factory.main.name
}

output "identity_principal_id" {
  description = "The principal ID of the factory's system-assigned identity (empty when no identity is configured)"
  value       = try(azurerm_data_factory.main.identity[0].principal_id, "")
}

output "credential_ids" {
  description = "The ARM IDs of the factory's named credentials, keyed by credential name (both flavors -- names share one namespace)"
  value = merge(
    { for name, credential in azurerm_data_factory_credential_user_managed_identity.main : name => credential.id },
    { for name, credential in azurerm_data_factory_credential_service_principal.main : name => credential.id },
  )
}

output "managed_private_endpoint_ids" {
  description = "The ARM IDs of the factory's managed private endpoints, keyed by endpoint name"
  value       = { for name, endpoint in azurerm_data_factory_managed_private_endpoint.main : name => endpoint.id }
}

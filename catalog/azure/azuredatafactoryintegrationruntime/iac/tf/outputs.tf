# Exactly one of the 3 resources exists (the spec's variant block
# decides); all share the {factory_id}/integrationRuntimes/{name} ID
# shape.
output "integration_runtime_id" {
  description = "The Azure Resource Manager ID of the integration runtime ({factory_id}/integrationRuntimes/{name})"
  value = try(
    azurerm_data_factory_integration_runtime_azure.main[0].id,
    azurerm_data_factory_integration_runtime_azure_ssis.main[0].id,
    azurerm_data_factory_integration_runtime_self_hosted.main[0].id,
    null
  )
}

output "integration_runtime_name" {
  description = "The integration runtime's name -- what linked services, data flow activities, and the SSIS proxy resolve against"
  value = try(
    azurerm_data_factory_integration_runtime_azure.main[0].name,
    azurerm_data_factory_integration_runtime_azure_ssis.main[0].name,
    azurerm_data_factory_integration_runtime_self_hosted.main[0].name,
    null
  )
}

# Azure issues authorization keys only for a PRIMARY self-hosted
# registration (never for the managed flavors, never for a linked
# registration). Azure returns them readable and the provider does
# not flag them -- marked sensitive here so the catalog's contract
# holds.
output "primary_authorization_key" {
  description = "The primary key a self-hosted agent joins with (empty for the managed flavors and linked registrations)"
  value       = try(azurerm_data_factory_integration_runtime_self_hosted.main[0].primary_authorization_key, "")
  sensitive   = true
}

output "secondary_authorization_key" {
  description = "The secondary key (the rotation partner; empty for the managed flavors and linked registrations)"
  value       = try(azurerm_data_factory_integration_runtime_self_hosted.main[0].secondary_authorization_key, "")
  sensitive   = true
}

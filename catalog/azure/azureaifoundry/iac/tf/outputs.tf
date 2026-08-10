output "ai_foundry_id" {
  description = "The Azure Resource Manager ID of the hub -- what AzureAiFoundryProject resources reference as their ai_services_hub_id"
  value       = azurerm_ai_foundry.main.id
}

output "ai_foundry_name" {
  description = "The name of the hub"
  value       = azurerm_ai_foundry.main.name
}

output "workspace_guid" {
  description = "The hub's immutable GUID (distinct from the ARM ID) -- what some data-plane SDKs and diagnostic settings identify the hub by"
  value       = azurerm_ai_foundry.main.workspace_id
}

output "discovery_url" {
  description = "The hub's regional discovery URL -- where SDKs resolve the hub's data-plane service endpoints from"
  value       = azurerm_ai_foundry.main.discovery_url
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the hub's system-assigned identity, when one is enabled"
  value       = try(azurerm_ai_foundry.main.identity[0].principal_id, "")
}

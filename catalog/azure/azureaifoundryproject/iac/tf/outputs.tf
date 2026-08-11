output "ai_foundry_project_id" {
  description = "The Azure Resource Manager ID of the project (projects are ARM workspaces, like their hub)"
  value       = azurerm_ai_foundry_project.main.id
}

output "ai_foundry_project_name" {
  description = "The name of the project"
  value       = azurerm_ai_foundry_project.main.name
}

output "project_guid" {
  description = "The project's immutable GUID (distinct from the ARM ID) -- what Foundry SDKs and data-plane calls identify the project by"
  value       = azurerm_ai_foundry_project.main.project_id
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the project's system-assigned identity, when one is enabled"
  value       = try(azurerm_ai_foundry_project.main.identity[0].principal_id, "")
}

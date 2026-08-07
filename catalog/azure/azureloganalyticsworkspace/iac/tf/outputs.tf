output "workspace_id" {
  description = "The Azure Resource Manager ID of the Log Analytics Workspace -- the FK seam downstream kinds reference"
  value       = azurerm_log_analytics_workspace.main.id
}

output "workspace_name" {
  description = "The name of the Log Analytics Workspace"
  value       = azurerm_log_analytics_workspace.main.name
}

# The provider's workspace_id attribute is the CUSTOMER GUID (what agents
# authenticate against), not the ARM resource ID -- exported under the
# unambiguous name so the two can never be confused.
output "workspace_customer_id" {
  description = "The workspace customer ID (GUID) agents and ingestion APIs identify the workspace by"
  value       = azurerm_log_analytics_workspace.main.workspace_id
}

output "resource_group_name" {
  description = "The name of the resource group containing the workspace"
  value       = azurerm_log_analytics_workspace.main.resource_group_name
}

output "primary_shared_key" {
  description = "The primary shared key for agent authentication (unusable when local authentication is disabled)"
  value       = azurerm_log_analytics_workspace.main.primary_shared_key
  sensitive   = true
}

output "secondary_shared_key" {
  description = "The secondary shared key for agent authentication -- rotate via the primary/secondary swap"
  value       = azurerm_log_analytics_workspace.main.secondary_shared_key
  sensitive   = true
}

output "identity_principal_id" {
  description = "The principal ID of the system-assigned managed identity (empty unless SYSTEM_ASSIGNED is enabled)"
  value       = try(azurerm_log_analytics_workspace.main.identity[0].principal_id, "")
}

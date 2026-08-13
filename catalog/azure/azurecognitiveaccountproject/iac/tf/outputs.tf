output "project_id" {
  description = "The Azure Resource Manager ID of the project"
  value       = azurerm_cognitive_account_project.main.id
}

output "project_name" {
  description = "The project's name"
  value       = azurerm_cognitive_account_project.main.name
}

output "endpoints" {
  description = "The project's data-plane endpoints, keyed by service label as ARM reports them"
  value       = azurerm_cognitive_account_project.main.endpoints
}

output "is_default" {
  description = "Whether ARM made this the account's default project (the first project created on an account becomes the default)"
  value       = azurerm_cognitive_account_project.main.default
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the project's system-assigned identity, when one is enabled"
  value       = try(azurerm_cognitive_account_project.main.identity[0].principal_id, "")
}

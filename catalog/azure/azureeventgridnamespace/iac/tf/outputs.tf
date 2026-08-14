output "namespace_id" {
  description = "The Azure Resource Manager ID of the Event Grid namespace"
  value       = azurerm_eventgrid_namespace.main.id
}

output "namespace_name" {
  description = "The namespace's name"
  value       = azurerm_eventgrid_namespace.main.name
}

output "identity_principal_id" {
  description = "The principal ID of the namespace's system-assigned identity (empty when no identity is configured)"
  value       = try(azurerm_eventgrid_namespace.main.identity[0].principal_id, "")
}

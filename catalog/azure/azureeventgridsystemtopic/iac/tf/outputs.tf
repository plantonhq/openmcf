output "system_topic_id" {
  description = "The Azure Resource Manager ID of the Event Grid system topic"
  value       = azurerm_eventgrid_system_topic.main.id
}

output "system_topic_name" {
  description = "The system topic's name"
  value       = azurerm_eventgrid_system_topic.main.name
}

output "metric_resource_id" {
  description = "The GUID-style identifier Azure Monitor uses for the topic's metrics"
  value       = azurerm_eventgrid_system_topic.main.metric_resource_id
}

output "identity_principal_id" {
  description = "The principal ID of the topic's system-assigned identity (empty when no system-assigned identity is configured)"
  value       = try(azurerm_eventgrid_system_topic.main.identity[0].principal_id, "")
}

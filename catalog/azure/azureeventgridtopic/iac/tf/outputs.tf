output "topic_id" {
  description = "The Azure Resource Manager ID of the Event Grid topic"
  value       = azurerm_eventgrid_topic.main.id
}

output "topic_name" {
  description = "The topic's name (the first label of its endpoint hostname)"
  value       = azurerm_eventgrid_topic.main.name
}

output "endpoint" {
  description = "The HTTPS endpoint publishers POST events to"
  value       = azurerm_eventgrid_topic.main.endpoint
}

output "primary_access_key" {
  description = "The primary access key (the aeg-sas-key header value)"
  value       = azurerm_eventgrid_topic.main.primary_access_key
  sensitive   = true
}

output "secondary_access_key" {
  description = "The secondary access key (the rotation partner)"
  value       = azurerm_eventgrid_topic.main.secondary_access_key
  sensitive   = true
}

output "identity_principal_id" {
  description = "The principal ID of the topic's system-assigned identity (empty when no identity is configured)"
  value       = try(azurerm_eventgrid_topic.main.identity[0].principal_id, "")
}

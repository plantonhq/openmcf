output "namespace_topic_id" {
  description = "The Azure Resource Manager ID of the namespace topic"
  value       = azurerm_eventgrid_namespace_topic.main.id
}

output "namespace_topic_name" {
  description = "The namespace topic's name"
  value       = azurerm_eventgrid_namespace_topic.main.name
}

output "domain_topic_id" {
  description = "The Azure Resource Manager ID of the domain topic ({domain_id}/topics/{name})"
  value       = azurerm_eventgrid_domain_topic.main.id
}

output "domain_topic_name" {
  description = "The topic's name (the value publishers put in the event's topic field)"
  value       = azurerm_eventgrid_domain_topic.main.name
}

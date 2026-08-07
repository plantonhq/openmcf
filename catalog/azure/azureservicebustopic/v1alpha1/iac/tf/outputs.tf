# The parent reference for AzureServiceBusSubscription and topic-scoped
# SAS rules, and the scope for topic-level data-plane role assignments
# (Azure Service Bus Data Sender).
output "topic_id" {
  description = "The Azure Resource Manager ID of the topic"
  value       = azurerm_servicebus_topic.main.id
}

output "topic_name" {
  description = "The name of the topic"
  value       = azurerm_servicebus_topic.main.name
}

# Parsed from the namespace ARM ID -- consumers frequently need the
# namespace/topic name pair (SDK clients, function bindings).
output "namespace_name" {
  description = "The name of the Service Bus namespace the topic lives in"
  value       = local.namespace_name
}

output "subscription_id" {
  description = "The Azure Resource Manager ID of the subscription"
  value       = azurerm_servicebus_subscription.main.id
}

output "subscription_name" {
  description = "The name of the subscription"
  value       = azurerm_servicebus_subscription.main.name
}

# Parsed from the topic ARM ID -- consumers receive by the
# namespace/topic/subscription triple.
output "topic_name" {
  description = "The name of the topic the subscription attaches to"
  value       = local.topic_name
}

output "namespace_name" {
  description = "The name of the Service Bus namespace"
  value       = local.namespace_name
}

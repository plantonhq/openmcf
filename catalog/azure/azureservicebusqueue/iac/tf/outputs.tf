# The scope for queue-level data-plane role assignments (Azure Service
# Bus Data Receiver/Sender) and the parent reference for queue-scoped
# SAS rules.
output "queue_id" {
  description = "The Azure Resource Manager ID of the queue"
  value       = azurerm_servicebus_queue.main.id
}

output "queue_name" {
  description = "The name of the queue"
  value       = azurerm_servicebus_queue.main.name
}

# Parsed from the namespace ARM ID -- consumers frequently need the
# namespace/queue name pair (SDK clients, function bindings).
output "namespace_name" {
  description = "The name of the Service Bus namespace the queue lives in"
  value       = local.namespace_name
}

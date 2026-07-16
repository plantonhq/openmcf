# The parent reference for consumer groups and hub-scoped authorization
# rules, and the scope for hub-level data-plane RBAC (grant "Azure Event
# Hubs Data Receiver/Sender" on exactly this hub).
output "event_hub_id" {
  description = "The Azure Resource Manager ID of the event hub"
  value       = azurerm_eventhub.main.id
}

output "event_hub_name" {
  description = "The name of the event hub (the Kafka topic name on the Kafka endpoint)"
  value       = azurerm_eventhub.main.name
}

# What partition-aware consumers enumerate (e.g. "0", "1", ...).
output "partition_ids" {
  description = "The identifiers of the hub's partitions"
  value       = azurerm_eventhub.main.partition_ids
}

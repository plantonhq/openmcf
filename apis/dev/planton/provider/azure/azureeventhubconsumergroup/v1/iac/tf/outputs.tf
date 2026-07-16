output "consumer_group_id" {
  description = "The Azure Resource Manager ID of the consumer group"
  value       = azurerm_eventhub_consumer_group.main.id
}

# What consumer applications pass to their SDK client alongside the hub
# name.
output "consumer_group_name" {
  description = "The name of the consumer group"
  value       = azurerm_eventhub_consumer_group.main.name
}

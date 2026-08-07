# What an AzureEventHubNamespace's dedicated_cluster_id references to
# place the namespace on this cluster.
output "cluster_id" {
  description = "The Azure Resource Manager ID of the Event Hubs cluster"
  value       = azurerm_eventhub_cluster.main.id
}

output "cluster_name" {
  description = "The name of the Event Hubs cluster"
  value       = azurerm_eventhub_cluster.main.name
}

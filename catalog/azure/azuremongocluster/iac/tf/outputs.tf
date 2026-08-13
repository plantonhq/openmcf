output "mongo_cluster_id" {
  description = "The Azure Resource Manager ID of the Mongo vCore cluster"
  value       = azurerm_mongo_cluster.main.id
}

output "mongo_cluster_name" {
  description = "The cluster's name (the first label of its hostname)"
  value       = azurerm_mongo_cluster.main.name
}

output "connection_string" {
  description = "The cluster's primary MongoDB connection string, administrator credentials substituted in (empty without a native administrator)"
  value       = try(azurerm_mongo_cluster.main.connection_strings[0].value, "")
  sensitive   = true
}

output "connection_strings" {
  description = "Every connection string Azure publishes for the cluster, keyed by Azure's name for it"
  value       = { for cs in azurerm_mongo_cluster.main.connection_strings : cs.name => cs.value }
  sensitive   = true
}

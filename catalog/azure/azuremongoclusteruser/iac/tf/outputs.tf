output "mongo_cluster_user_id" {
  description = "The Azure Resource Manager ID of the user grant"
  value       = azurerm_mongo_cluster_user.main.id
}

output "mongo_cluster_user_name" {
  description = "The grant's ARM name -- the granted principal's object id"
  value       = azurerm_mongo_cluster_user.main.object_id
}

output "access_policy_assignment_id" {
  description = "The Azure Resource Manager ID of the access policy assignment"
  value       = azurerm_redis_cache_access_policy_assignment.main.id
}

output "access_policy_assignment_name" {
  description = "The assignment's name within the cache"
  value       = azurerm_redis_cache_access_policy_assignment.main.name
}

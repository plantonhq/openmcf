output "access_policy_id" {
  description = "The Azure Resource Manager ID of the access policy"
  value       = azurerm_redis_cache_access_policy.main.id
}

# What AzureRedisCacheAccessPolicyAssignment.access_policy_name references
# to grant this policy to an identity.
output "access_policy_name" {
  description = "The policy's name within the cache"
  value       = azurerm_redis_cache_access_policy.main.name
}

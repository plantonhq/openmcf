output "access_policy_assignment_id" {
  description = "The Azure Resource Manager ID of the access policy assignment"
  value       = azurerm_managed_redis_access_policy_assignment.main.id
}

# Azure names the assignment after the granted object ID, so the name
# equals the principal's GUID.
output "access_policy_assignment_name" {
  description = "The assignment's name (the granted object ID)"
  value       = azurerm_managed_redis_access_policy_assignment.main.object_id
}

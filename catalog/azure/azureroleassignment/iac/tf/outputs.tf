output "role_assignment_id" {
  description = "The fully-scoped Azure Resource Manager ID of the role assignment"
  value       = azurerm_role_assignment.main.id
}

output "name" {
  description = "The assignment's GUID resource name (pinned or generated)"
  value       = azurerm_role_assignment.main.name
}

output "scope" {
  description = "The scope the role was granted at"
  value       = azurerm_role_assignment.main.scope
}

output "role_definition_id" {
  description = "The role definition ID Azure resolved (even when the spec referenced the role by name)"
  value       = azurerm_role_assignment.main.role_definition_id
}

output "principal_id" {
  description = "The Azure AD object ID of the principal the role was granted to"
  value       = azurerm_role_assignment.main.principal_id
}

output "principal_type" {
  description = "The principal type Azure recorded for the assignment"
  value       = azurerm_role_assignment.main.principal_type
}

output "role_definition_id" {
  description = "The fully-scoped ARM ID of the role definition (what an AzureRoleAssignment's role_definition_id consumes)"
  # azurerm names the fully-scoped ARM ID `role_definition_resource_id` and
  # the bare GUID `role_definition_id`. Planton's Azure surface consistently
  # uses `role_definition_id` for the fully-scoped form (that is what a role
  # assignment binds), so the output name intentionally follows the platform
  # contract rather than the provider attribute name.
  value = azurerm_role_definition.main.role_definition_resource_id
}

output "role_definition_guid" {
  description = "The definition's GUID resource name (pinned or generated)"
  value       = azurerm_role_definition.main.role_definition_id
}

output "role_name" {
  description = "The role's tenant-unique display name as deployed"
  value       = azurerm_role_definition.main.name
}

output "scope" {
  description = "The scope the definition was created at"
  value       = azurerm_role_definition.main.scope
}

output "assignable_scopes" {
  description = "The assignable scopes Azure recorded (the definition's own scope when the spec omitted them)"
  value       = azurerm_role_definition.main.assignable_scopes
}

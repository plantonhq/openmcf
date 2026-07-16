output "schema_group_id" {
  description = "The Azure Resource Manager ID of the schema group"
  value       = azurerm_eventhub_namespace_schema_group.main.id
}

# What schema-registry serializers address at runtime, alongside the
# namespace's fully-qualified hostname.
output "schema_group_name" {
  description = "The name of the schema group"
  value       = azurerm_eventhub_namespace_schema_group.main.name
}

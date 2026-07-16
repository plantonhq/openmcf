# The parent reference for every Service Bus child kind (queue, topic,
# authorization rule, geo-DR pairing), the scope for namespace-wide
# data-plane RBAC, and the private-endpoint target (subresource
# "namespace").
output "namespace_id" {
  description = "The Azure Resource Manager ID of the Service Bus namespace"
  value       = azurerm_servicebus_namespace.main.id
}

output "namespace_name" {
  description = "The name of the Service Bus namespace"
  value       = azurerm_servicebus_namespace.main.name
}

output "endpoint" {
  description = "The Service Bus endpoint URL (https://{name}.servicebus.windows.net:443/)"
  value       = azurerm_servicebus_namespace.main.endpoint
}

# Empty unless the identity block includes SYSTEM_ASSIGNED. Grant this
# principal access on other resources (e.g. Key Vault for CMK).
output "identity_principal_id" {
  description = "The system-assigned identity's principal ID"
  value       = try(azurerm_servicebus_namespace.main.identity[0].principal_id, "")
}

# The root SAS rule's (RootManageSharedAccessKey) credential faces --
# full manage rights over the whole namespace. Quick-start/break-glass
# credentials; production workloads mint least-privilege rules with
# AzureServiceBusAuthorizationRule or go keyless.
output "default_primary_connection_string" {
  description = "The root SAS rule's primary connection string"
  value       = azurerm_servicebus_namespace.main.default_primary_connection_string
  sensitive   = true
}

output "default_secondary_connection_string" {
  description = "The root SAS rule's secondary connection string (rotation partner)"
  value       = azurerm_servicebus_namespace.main.default_secondary_connection_string
  sensitive   = true
}

output "default_primary_key" {
  description = "The root SAS rule's primary key"
  value       = azurerm_servicebus_namespace.main.default_primary_key
  sensitive   = true
}

output "default_secondary_key" {
  description = "The root SAS rule's secondary key (rotation partner)"
  value       = azurerm_servicebus_namespace.main.default_secondary_key
  sensitive   = true
}

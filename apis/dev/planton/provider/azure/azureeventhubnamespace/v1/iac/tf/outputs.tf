# The parent reference for every Event Hubs child kind (hub, consumer
# group, authorization rule, schema group, geo-DR pairing, CMK), the
# scope for namespace-wide data-plane RBAC, and the private-endpoint
# target (subresource "namespace").
output "namespace_id" {
  description = "The Azure Resource Manager ID of the Event Hubs namespace"
  value       = azurerm_eventhub_namespace.main.id
}

output "namespace_name" {
  description = "The name of the Event Hubs namespace"
  value       = azurerm_eventhub_namespace.main.name
}

# Empty unless the identity block includes SYSTEM_ASSIGNED. Grant this
# principal access on other resources (e.g. Storage for capture, Key
# Vault for CMK).
output "identity_principal_id" {
  description = "The system-assigned identity's principal ID"
  value       = try(azurerm_eventhub_namespace.main.identity[0].principal_id, "")
}

# The root SAS rule's (RootManageSharedAccessKey) credential faces --
# full manage rights over the whole namespace. Quick-start/break-glass
# credentials; production workloads mint least-privilege rules with
# AzureEventHubAuthorizationRule or go keyless.
output "default_primary_connection_string" {
  description = "The root SAS rule's primary connection string"
  value       = azurerm_eventhub_namespace.main.default_primary_connection_string
  sensitive   = true
}

output "default_secondary_connection_string" {
  description = "The root SAS rule's secondary connection string (rotation partner)"
  value       = azurerm_eventhub_namespace.main.default_secondary_connection_string
  sensitive   = true
}

output "default_primary_key" {
  description = "The root SAS rule's primary key"
  value       = azurerm_eventhub_namespace.main.default_primary_key
  sensitive   = true
}

output "default_secondary_key" {
  description = "The root SAS rule's secondary key (rotation partner)"
  value       = azurerm_eventhub_namespace.main.default_secondary_key
  sensitive   = true
}

# Alias faces address the geo-DR alias hostname -- only populated when
# the namespace carries an AzureEventHubDisasterRecoveryConfig pairing;
# empty otherwise. DR-aware clients hold these so a failover needs no
# reconfiguration.
output "default_primary_connection_string_alias" {
  description = "The root SAS rule's primary connection string addressing the geo-DR alias"
  value       = azurerm_eventhub_namespace.main.default_primary_connection_string_alias
  sensitive   = true
}

output "default_secondary_connection_string_alias" {
  description = "The secondary alias connection string (rotation partner)"
  value       = azurerm_eventhub_namespace.main.default_secondary_connection_string_alias
  sensitive   = true
}

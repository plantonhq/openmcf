# Azure materializes the one logical policy on BOTH accounts, so there
# are two ARM IDs sharing one policy GUID.
output "source_object_replication_id" {
  description = "The policy's ARM ID on the source account"
  value       = azurerm_storage_object_replication.main.source_object_replication_id
}

output "destination_object_replication_id" {
  description = "The policy's ARM ID on the destination account (the authoritative copy)"
  value       = azurerm_storage_object_replication.main.destination_object_replication_id
}

# What `az storage account or-policy show --policy-id` and the
# monitoring surfaces key on.
output "policy_id" {
  description = "The server-assigned policy GUID shared by both sides"
  value       = local.policy_id
}

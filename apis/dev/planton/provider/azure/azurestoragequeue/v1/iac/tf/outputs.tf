# The ARM id data-plane role assignments (Storage Queue Data
# Contributor / Message Processor / Message Sender) scope to for
# queue-level access.
output "queue_id" {
  description = "The Azure Resource Manager ID of the queue"
  value       = azurerm_storage_queue.main.id
}

output "queue_name" {
  description = "The name of the queue"
  value       = azurerm_storage_queue.main.name
}

# Parsed from the account ARM ID -- consumers frequently need the
# account/queue name pair (SDK clients, Functions queue triggers).
output "storage_account_name" {
  description = "The name of the storage account the queue lives in"
  value       = local.storage_account_name
}

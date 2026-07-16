# The Storage queue, addressed by the parent account's ARM ID (the
# control-plane path -- the account-name form is the provider's legacy
# data-plane path, removed in azurerm v5). Queues carry no Azure tags:
# ARM does not support tags on queueServices/queues, so the platform's
# identity tags live on the account.
resource "azurerm_storage_queue" "main" {
  name               = var.spec.queue_name
  storage_account_id = var.spec.storage_account_id

  # Queue metadata is NOT Azure tags -- free-form key/value pairs on the
  # queue itself, visible to anyone who can read queue properties.
  metadata = var.spec.metadata
}

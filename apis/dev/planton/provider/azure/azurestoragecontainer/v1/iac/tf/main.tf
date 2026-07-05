# The blob container, addressed by the parent account's ARM ID (the
# control-plane path -- the account-name form is the provider's legacy
# data-plane path, removed in azurerm v5). Containers carry no Azure
# tags: ARM does not support tags on blobServices/containers, so the
# platform's identity tags live on the account.
resource "azurerm_storage_container" "main" {
  name               = var.spec.container_name
  storage_account_id = var.spec.storage_account_id

  # Anonymous access also requires the ACCOUNT's
  # allow_nested_items_to_be_public to be true; when the account forbids
  # it, Azure forces private regardless of this value.
  container_access_type = local.container_access_type

  # Sub-account key isolation: blobs without their own scope encrypt
  # under this one. Both fields are fixed at creation.
  default_encryption_scope          = local.default_encryption_scope
  encryption_scope_override_enabled = var.spec.encryption_scope_override_enabled

  metadata = var.spec.metadata
}

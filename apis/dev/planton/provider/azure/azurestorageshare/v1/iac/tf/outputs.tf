# The management-plane ARM id -- what ARM reads and policy target.
output "share_id" {
  description = "The Azure Resource Manager ID of the share"
  value       = azurerm_storage_share.main.id
}

# Azure Files RBAC scopes to a DIFFERENT segment than the management ID
# (.../fileServices/default/fileshares/{name}); this output spares role
# assignments (Storage File Data SMB Share Reader/Contributor) from
# rewriting the management ID by hand.
output "rbac_scope_id" {
  description = "The scope data-plane role assignments target for share-level file access"
  value       = azurerm_storage_share.main.rbac_scope_id
}

output "share_name" {
  description = "The name of the share"
  value       = azurerm_storage_share.main.name
}

# Parsed from the account ARM ID -- consumers frequently need the
# account/share name pair (mount commands, CSI volume definitions).
output "storage_account_name" {
  description = "The name of the storage account the share lives in"
  value       = local.storage_account_name
}

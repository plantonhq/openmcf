# Register the file share under the policy's protection. Creation only
# REGISTERS protection -- the first backup runs on the policy's
# schedule, not immediately. The share's storage account must already
# be REGISTERED with the vault (AzureBackupContainerStorageAccount):
# the provider runs an Inquire pass to discover protectable shares
# inside the registered container and fails loudly when the account is
# not registered. ARM names the protected item by the share's SYSTEM
# name (AzureFileShare;{system-name}), not its friendly name.
#
# Destroy semantics kept deliberately at the engines' defaults:
# destroying stops protection AND deletes the backup data (vault soft
# delete may hold it 14 days) -- recorded on the spec.
resource "azurerm_backup_protected_file_share" "main" {
  resource_group_name       = var.spec.resource_group
  recovery_vault_name       = var.spec.recovery_vault_name
  source_storage_account_id = var.spec.source_storage_account_id
  source_file_share_name    = var.spec.source_file_share_name

  # The spec's ONLY updatable field -- re-pointing the policy updates
  # in place; everything else is ForceNew on the provider.
  backup_policy_id = var.spec.backup_policy_id
}

# Register the storage account with the vault as a backup container
# (.../vaults/{vault}/backupFabrics/Azure/protectionContainers/
# StorageContainer;storage;{sa-rg};{sa-name}) -- the prerequisite for
# protecting any of the account's file shares. Registration is free
# and moves no data. ARM derives the container's own name from the
# storage account's group and name.
#
# Every argument is ForceNew (ARM has no update on protection
# containers). While registered, Azure Backup places a resource lock
# on the storage account; destroying this resource unregisters and
# removes the lock -- and REFUSES while any of the account's shares
# are still protected (destroy the protections first; the GUIDE
# carries the teardown recipe).
resource "azurerm_backup_container_storage_account" "main" {
  resource_group_name = var.spec.resource_group
  recovery_vault_name = var.spec.recovery_vault_name
  storage_account_id  = var.spec.storage_account_id
}

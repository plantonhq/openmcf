# Azure Backup Protected File Share -- Operational Guide

Judgment that saves real time when protecting file shares. The field reference lives in the API Explorer; this is the operational layer above it.

## The chain is vault → registration → policy → THIS

Protection fails without its prerequisites live: the storage account must be REGISTERED with the vault (AzureBackupContainerStorageAccount) and the policy must exist in the SAME vault. The failure when the account is not registered is explicit -- `fileshare not found in protectable or protected fileshares, make sure Storage Account ... is registered` -- and it costs minutes of discovery time before it surfaces. Keep the default reference (through the registration's echoed output) and the ordering is automatic.

## The first backup is NOT at create time

Creating the binding registers protection; the first recovery point lands at the policy's next scheduled run. A share protected at 09:00 against a 23:00 daily policy has NO restore point until 23:00 -- run an on-demand backup (portal/CLI) if you need one sooner.

## Only the policy re-points in place

`backupPolicyId` is the spec's single updatable field. Everything else -- vault, account, share name -- replaces the protection with a NEW protected item, and the old item's backup data follows the vault's soft-delete rules. Treat those fields as identity.

## Destroy deletes the data, and soft delete holds it

Destroying this resource stops protection AND deletes the backup data. Vault soft delete -- always on since Azure's secure-by-default policy -- holds the deleted item for 14 days: it still counts against the vault, blocks unregistering the storage account, and blocks vault deletion. Teardown order is protections → registration → vault, and a just-deleted protection can delay the registration's unregister until the soft-deleted item is purged or expires.

## Expect the 80-minute timeout class

Protection configuration is a long-running ARM operation (the provider allows 80 minutes for create/update/delete). Normal runs finish in minutes; budget pipelines for the class, not the average.

## Azure renames the share internally

The protected item's ARM ID carries the share's SYSTEM name (`AzureFileShare;{system-name}`), which differs from its friendly name. When correlating with `az backup item list`, match by the friendly name field, not the ID segment.

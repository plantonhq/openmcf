# Overview

The **AzureBackupContainerStorageAccount** component registers a storage account with a Recovery Services vault as a backup container -- the one-time prerequisite that lets the account's file shares be protected. ONE registration per storage-account-and-vault pair; each share then gets its own AzureBackupProtectedFileShare binding. Registration is free and moves no data.

## Purpose

- **The missing link made explicit**: Azure discovers protectable file shares only inside REGISTERED storage accounts -- protection without registration fails. Modeling the registration as its own resource makes the dependency visible, orderable, and destroyable in the right sequence.
- **One registration, many protections**: the account registers once; every share protection in the account rides the same registration -- exactly the provider's own resource shape.
- **Wire-through dependency ordering**: the registration echoes the storage account's ARM ID as an output; protected shares reference THAT output, so the registration always deploys first -- the reference carries both the value and the deploy-order edge.

## Key Features

- Full azurerm v5 surface (three arguments -- the resource is deliberately small).
- All-ForceNew semantics documented where they bite: ARM has no update on protection containers.
- The operational edges on the record: Azure Backup's resource lock on the registered account, and unregister's refusal while shares remain protected.

## Use Cases

- **The file-share backup chain**: vault → registration → policy → protected shares. The registration is step two.
- **Charts**: landing-zone or app charts that provision a storage account and protect its shares in one deploy -- the registration node makes the ordering automatic.

## Future Enhancements

- Blob and disk protection under the modern Data Protection family use their own binding kinds -- this registration serves the classic (Recovery Services) file-share path.

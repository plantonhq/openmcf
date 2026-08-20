# Azure Backup Container (Storage Account) -- Operational Guide

Judgment that saves real time when registering storage accounts for backup. The field reference lives in the API Explorer; this is the operational layer above it.

## One registration per account-and-vault pair -- never per share

Registering an already-registered account fails ("already exists") -- the registration is a singleton on the (vault, storage account) pair. Deploy ONE AzureBackupContainerStorageAccount per account, then any number of AzureBackupProtectedFileShare bindings ride it. If a chart provisions several shares in one account, they all reference the same registration.

## Wire protections THROUGH the registration

AzureBackupProtectedFileShare's `sourceStorageAccountId` should reference THIS resource's `storage_account_id` output (its default reference does exactly that), not the storage account directly. The value is identical -- the registration echoes the account's ARM ID -- but the reference edge is what guarantees the registration deploys before the protection. Reference the account directly and a fresh chart deploy can race: the protection's discovery pass finds an unregistered account and fails.

## Expect the resource lock -- and it blocks share deletes too

While registered, Azure Backup places a DoNotDelete ARM lock on the storage account (protecting the backups' source from accidental deletion). Teams auditing locks should expect it; it is removed at unregister, not by hand. The lock's scope covers the account's children: deleting ANY file share in a registered account fails with `ScopeLocked` naming the account -- even a share that was never protected. If infrastructure automation tears down shares alongside this registration, the unregister must run FIRST; a share delete that hits ScopeLocked is ordering, not a stuck lock.

## Teardown runs strictly backwards

Unregister REFUSES while any of the account's shares are still protected. The order is always: destroy the AzureBackupProtectedFileShare bindings, then this registration, then (if ever) the vault. One more wrinkle: vault soft delete -- always on since Azure's secure-by-default policy -- can hold a deleted protection's data for 14 days, and a held item can delay the unregister. If unregister fails right after destroying protections, the soft-deleted items are why; undelete-and-purge them or wait out the window.

## Same region, or it fails at apply

Azure Files backup is regional: the storage account must live in the vault's region. Cross-region registration fails at apply time with a service error -- there is nothing the manifest can pre-check across two live resources, so plan regions deliberately.

# Azure Data Protection Backup Vault -- Operational Guide

Judgment calls that matter when you run Data Protection vaults in production.

## Two vault generations -- pick the right one first

Azure Backup has two families. The CLASSIC Recovery Services vault protects IaaS VMs and Azure Files shares. THIS vault -- Data Protection -- protects managed disks, blob storage, AKS clusters, MySQL/PostgreSQL flexible servers, and Data Lake storage. They are separate services with separate policies and separate protections; a workload belongs to exactly one family. Check which family covers your datasource before creating either vault.

## Datastore tier and redundancy are create-time decisions

Both `datastoreType` and `redundancy` replace the vault when changed. `VaultStore` + `GeoRedundant` is the production default posture; drop to `LocallyRedundant` only for dev/test where losing backups with the region is acceptable. There is no in-place migration -- a wrong choice means a new vault and re-protecting everything.

## Cross-region restore is a one-way door in one direction

Enabling `crossRegionRestoreEnabled` on a geo-redundant vault is an in-place update. DISABLING it replaces the vault. Enable it when the paired-region restore capability justifies the storage premium; treat the decision as permanent.

## Immutability: Unlocked is the trial run, Locked is forever

`Unlocked` immutability blocks backup deletion but can itself be turned off -- run production vaults there first, until retention settings have survived a few review cycles. `Locked` is permanent: leaving it replaces the vault, and Azure will not shorten retention or delete backups inside it. Lock only when the org is certain.

## Soft delete has its own ratchet

`On` (the default) retains deleted backup data for `retentionDurationInDays` (14-180 days) -- deletion becomes recoverable. `AlwaysOn` makes that permanent: the setting can never leave AlwaysOn without replacing the vault. `Off` exists for dev/test churn -- and is exactly the class of privileged operation a Resource Guard exists to gate in production.

## CMK encryption can never be removed

Once the `encryption` block is applied, customer-managed-key encryption is on for the vault's lifetime -- Azure has no path back to Microsoft-managed keys (the provider's delete for the composed CMK arm is a documented no-op). The KEY can be rotated freely: the reference's versionless default means new key versions propagate automatically. The vault's SYSTEM-assigned identity unwraps the key (Azure hardcodes this), so grant it wrap/unwrap on the Key Vault key before enabling.

## Deletion takes longer than the API admits

Azure's vault delete returns before the vault is fully gone; the module polls until the name is actually free (the provider's own workaround). Budget a few extra minutes in teardown automation, and remember a vault refuses deletion while backup instances remain inside it.

## Multi-user authorization is an org-structure decision

For vaults whose backups must survive a compromised administrator, pair the vault family with a Data Protection Resource Guard living in a DIFFERENT administrator's scope. The guard's value comes entirely from that separation -- a guard the same admin controls is a speed bump, not a control.

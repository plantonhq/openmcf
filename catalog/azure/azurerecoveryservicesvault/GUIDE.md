# Azure Recovery Services Vault -- Operational Guide

Judgment that saves real time when running Recovery Services vaults. The field reference lives in the API Explorer; this is the operational layer above it.

## Decide redundancy before the first protection, not after

`storageModeType` looks like an in-place update -- and it is, right up until the vault protects its first item, after which Azure locks it. GeoRedundant (the default) is the production posture; LocallyRedundant roughly halves backup storage cost for workloads whose loss tolerance genuinely allows a single-region copy. If you catch yourself wanting to change redundancy on a populated vault, the real operation is a new vault and a migration of protections.

## Cross-region restore is a one-way door in one direction

Enabling `crossRegionRestoreEnabled` on a geo-redundant vault is a free in-place update. DISABLING it replaces the vault -- the provider's one-way ForceNew. Enable it when the paired-region restore capability is worth the modest storage premium; treat disabling as "rebuild".

## Immutability: Unlocked is the trial run, Locked is forever

The transitions run `Disabled <-> Unlocked -> Locked`. Run production vaults at `Unlocked` first -- it enforces the same protections but stays reversible. Move to `Locked` only when the organization has lived with immutability long enough to trust its own retention settings: leaving Locked means REPLACING the vault, and a locked vault's backup data cannot be deleted early by anyone. Setting `Locked` directly in a fresh manifest works (the module stages it through Unlocked automatically), but consider whether day one is really the moment.

## Deleting a vault is deliberately hard

Azure refuses to delete a vault holding protected items, and these modules keep that guard (the engine-level purge switch stays off). The teardown order is: flip protections off (or destroy the protection resources -- which deletes their backup data), wait for soft-deleted items to purge if soft delete holds them, then delete the vault. Budget real time for this in environment teardowns -- soft delete holds backup data for 14 days by default, and a "why won't my resource group delete" investigation usually ends at a vault.

## CMK encryption is a ratchet with three teeth

Once `encryption` is set: it can never be removed, `infrastructureEncryptionEnabled` can never change, and the `sku` freezes. The identity needs wrap/unwrap on the key BEFORE the vault deploys -- with a user-assigned identity you can compose the grant first (that is the reason to prefer it over system-assigned for CMK); with system-assigned there is a bootstrap hop (deploy without encryption, grant, then update). The key reference targets the VERSIONLESS URI by default, so rotation propagates without touching the vault.

## The identity never downgrades

Once the vault has an identity, removing it -- or switching from both flavors to just one -- is rejected by the provider (Azure's CMK machinery breaks silently otherwise). Adding flavors is fine. Decide the identity story once, at creation.

## Monitoring defaults are all-on; turn off deliberately

All five alert switches default ON service-side -- an unset `monitoring` block IS the sensible posture. The three Site Recovery-related switches are new in provider v5 and only the Terraform engine can turn them OFF (the Pulumi engine's SDK predates them and its module fails loudly on an explicit false); job-failure and critical-operation alerts switch freely on both engines.

## Multi-user authorization is an org-structure decision

`resourceGuardId` binds the vault to a Resource Guard so privileged operations (disabling soft delete, reducing retention) need an approval through the guard. The security value comes from the guard living in a DIFFERENT administrator's scope -- a guard the same admin controls is ceremony, not protection. One guard per vault (ARM pins the association name).

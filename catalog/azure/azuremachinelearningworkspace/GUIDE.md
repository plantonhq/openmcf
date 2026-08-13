# Azure Machine Learning Workspace -- Operational Guide

Judgment that saves real time when running ML workspaces. The field reference lives in the API Explorer; this is the operational layer above it.

## Deletion is a soft delete -- the name stays taken

A deleted workspace becomes a purgeable ghost that keeps holding the workspace NAME (the Key Vault recycle-bin pattern). Recreating under the same name fails until the ghost is purged -- `az ml workspace list --archived` shows them. The provider purges on destroy only when the `machine_learning.purge_soft_deleted_workspace_on_destroy` features flag is set; without it, delete-recreate cycles under one name will surprise you.

## The companion services are a one-way door

Storage account, key vault, application insights, and container registry are all ForceNew: re-pointing any of them replaces the workspace -- and with it every datastore, compute, and endpoint registered on it. Treat the companion set as part of the workspace's identity: provision them together, size the storage account for the workspace's lifetime, and never share a "temporary" vault you plan to swap later.

## The default storage account has hard shape requirements

ARM rejects a workspace whose default storage account has hierarchical namespace (Data Lake Gen2) enabled, and premium storage is not supported as default workspace storage. General-purpose v2, HNS off, same region -- boring on purpose. Data Lake accounts belong in DATASTORES pointing at them, not under the workspace itself.

## Isolation mode changes are one-directional in practice

Moving from `DISABLED` toward `ALLOW_ONLY_APPROVED_OUTBOUND` tightens the workspace and works in place; loosening back out of approved-outbound is rejected by ARM once the managed network has provisioned. Decide the target posture up front. Under approved-outbound, remember the built-in rules Azure creates for its own services -- your rule lists only ADD to them, and all three rule types share ONE name namespace.

## The managed network provisions lazily -- and that bites the first job

By default the managed VNet materializes on first compute creation, which adds minutes to the first job and can fail late on quota. Set `managedNetwork.provisionOnCreationEnabled: true` on isolated workspaces so the network exists (and fails, if it must) at deploy time instead.

## Grant the identity before you need it

The workspace identity needs Blob Data Contributor on the storage account for identity-based datastore access and wrap/unwrap on the CMK key BEFORE encryption is configured. With `USER_ASSIGNED` identities you can compose those grants ahead of workspace creation -- that is the reason to prefer them in locked-down estates.

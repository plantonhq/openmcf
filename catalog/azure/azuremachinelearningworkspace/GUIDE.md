# Azure Machine Learning Workspace -- Operational Guide

Judgment that saves real time when running ML workspaces. The field reference lives in the API Explorer; this is the operational layer above it.

## Deletion is a soft delete -- the name stays taken

A deleted workspace becomes a purgeable ghost that keeps holding the workspace NAME (the Key Vault recycle-bin pattern). Both Planton modules enable the provider's `machine_learning.purge_soft_deleted_workspace_on_destroy` features flag, so a Planton destroy purges the ghost and the name frees immediately -- destroy means destroy, and there is no soft-delete recovery window to lean on. A workspace deleted OUTSIDE Planton (portal, `az ml workspace delete` without `--permanently-delete`) still ghosts and blocks a redeploy under the same name until purged. Know the boundary: NO CLI or REST API lists ghosts -- `az ml workspace list` has no soft-delete flag and Resource Graph indexes active resources only. The Azure portal's "Recently deleted" view (Azure Machine Learning service page, per region) is the one listing surface, and purging an existing ghost happens there too. The tell that a ghost is in your way: a create fails on a name conflict while no active workspace of that name exists anywhere.

## The companion services are a one-way door

Storage account, key vault, application insights, and container registry are all ForceNew: re-pointing any of them replaces the workspace -- and with it every datastore, compute, and endpoint registered on it. Treat the companion set as part of the workspace's identity: provision them together, size the storage account for the workspace's lifetime, and never share a "temporary" vault you plan to swap later.

## The default storage account has hard shape requirements

ARM rejects a workspace whose default storage account has hierarchical namespace (Data Lake Gen2) enabled, and premium storage is not supported as default workspace storage. General-purpose v2, HNS off, same region -- boring on purpose. Data Lake accounts belong in DATASTORES pointing at them, not under the workspace itself.

## Isolation mode changes are one-directional in practice

Moving from `DISABLED` toward `ALLOW_ONLY_APPROVED_OUTBOUND` tightens the workspace and works in place; loosening back out of approved-outbound is rejected by ARM once the managed network has provisioned. Decide the target posture up front. Under approved-outbound, remember the built-in rules Azure creates for its own services -- your rule lists only ADD to them, and all three rule types share ONE name namespace.

## The managed network provisions lazily -- and that bites the first job

By default the managed VNet materializes on first compute creation, which adds minutes to the first job and can fail late on quota. Set `managedNetwork.provisionOnCreationEnabled: true` on isolated workspaces so the network exists (and fails, if it must) at deploy time instead.

## Approved-outbound + provision-on-creation is a fragile create, and a PE outbound rule is a fragile delete

Two live attempts at a Basic workspace with `ALLOW_ONLY_APPROVED_OUTBOUND` and `provision_on_creation_enabled` failed in different ways. Without a private-endpoint outbound rule, ARM rolled the workspace back mid-create after ~6 minutes (`Bad request to get identity secret: The workspace identity has been deleted` / workspace NotFound). With a PE rule to a Key Vault, create eventually succeeded and then destroy 409-looped for 40+ minutes (`privateEndpointConnectionProxies/validate` failing; workspace delete returning `InternalServerError` / 409) while the workspace stayed `Succeeded`. Manual recovery is `az ml workspace delete --yes --permanently-delete` (itself ~15 minutes; first `az ml` hangs forever unless the `ml` extension is already installed). Prefer the default network for day-one workspaces; if you need approved-outbound, add PE rules after the workspace is up, and delete those rules (and wait for the target-side connection to drop) before deleting the workspace.

## Grant the identity before you need it

The workspace identity needs Blob Data Contributor on the storage account for identity-based datastore access and wrap/unwrap on the CMK key BEFORE encryption is configured. With `USER_ASSIGNED` identities you can compose those grants ahead of workspace creation -- that is the reason to prefer them in locked-down estates.

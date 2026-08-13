# Azure Disk Snapshot

Deploys a managed disk snapshot -- a point-in-time copy of a disk used for backup, cloning, and as the source of gallery image versions. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Disk snapshot** -- the point-in-time copy: creation mode and source, incremental mode, network posture, optional legacy ADE encryption settings, tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **A source** -- "Copy" mode references a managed disk (or another snapshot); "Import" mode references a VHD blob and its storage account.

### Azure Subscription

- **Prefer incremental snapshots** -- they store only the delta on standard storage; full snapshots store the whole disk. The choice is fixed at creation.
- **Snapshots are regional** -- create the snapshot where its source disk lives.
- **This is a manual/pipeline snapshot, not a backup policy** -- for scheduled, retention-managed VM backups use the Recovery Services kinds (AzureBackupPolicyVm and friends); this kind is the single point-in-time artifact.

## Deploy

### Console

Open the deployment store, find **Azure Disk Snapshot**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Incremental Disk Backup** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f disk-snapshot.yaml
```

## After Deploy

Use the `snapshot_id` output where the copy is consumed: a gallery image version's `osDiskSnapshotId` (the golden-image chain), a new managed disk's creation source, or a cross-region copy job. A snapshot bills only for the storage it holds -- delete pre-change backups once the change settles.

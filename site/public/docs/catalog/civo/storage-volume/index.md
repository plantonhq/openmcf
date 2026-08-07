---
title: "Storage Volume"
description: "Storage Volume deployment documentation"
icon: "package"
order: 100
componentName: "civovolume"
---

# Storage Volume on Civo

Deploys a block storage volume on Civo Cloud with configurable size, region, and optional filesystem pre-formatting. Volumes provide persistent, expandable storage that can be attached to Civo compute instances independently of instance lifecycle. Integrates with Planton's Provider Connections for Civo credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Block Storage Volume** -- a `civo_volume` resource created in the target region with the requested capacity in GiB and the name specified in `volumeName`
- **Civo Labels** -- metadata labels derived from the resource identity, applied automatically to the volume for organizational tracking

The volume is created in a detached state. Attach it to a compute instance separately via a volume attachment resource or a Kubernetes CSI driver.

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A target region** -- the volume must be created in the same region as any instance that will attach to it. Check available regions via the Civo CLI (`civo region ls`) or Civo dashboard.
- **Sufficient block storage quota** in the target region for the requested volume size.
- **A volume snapshot ID** if creating the volume from a snapshot (optional, currently limited to CivoStack private cloud deployments).

## Deploy

### Console

Open the deployment store, find **Storage Volume on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **General-Purpose ext4** preset in the [Presets](#presets) tab for a ready-to-use storage volume.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoVolume
metadata:
  name: app-data
  org: acme-corp
  env: prod
spec:
  volumeName: app-data
  region: lon1
  sizeGib: 50
```

```shell
planton apply -f civo-volume.yaml
```

This creates a 50 GiB unformatted block storage volume in Civo's London region. No filesystem formatting, snapshot source, or tags are configured. A Stack Job tracks the provisioning in real time.

Attach the volume to a compute instance after creation -- volumes are provisioned in a detached state by default.

## Key Configuration

These are the most important decisions when configuring a Civo storage volume. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Volume name** -- The `volumeName` field sets the volume's name on Civo. Must be lowercase letters, numbers, and hyphens only, starting with a letter and ending with a letter or number (maximum 64 characters).

**Volume size** -- The `sizeGib` field sets the volume capacity in GiB, from 1 to 16,000. Start with a size that matches your current needs -- volumes can be resized later if more capacity is required.

**Region** -- The `region` field accepts a Civo region code (e.g., `lon1`, `nyc1`, `fra1`). The volume must be in the same region as any instance you plan to attach it to. Volumes cannot be moved across regions after creation.

**Filesystem type** -- The `filesystemType` field requests a pre-formatted filesystem (`ext4` or `xfs`). Note that the upstream Civo provider does not currently expose filesystem formatting -- the volume is created unformatted regardless of this setting. Format the volume manually or via cloud-init after attachment.

**Snapshot source** -- The `snapshotId` field creates the volume from an existing snapshot, inheriting the snapshot's data and minimum size. Snapshot-based creation is currently available only on CivoStack (private cloud) deployments, not on public Civo Cloud.

**Tags** -- The `tags` field applies organizational tags to the volume. Note that the upstream Civo Volume provider does not currently apply tags to the cloud resource -- tags are recorded in Planton metadata only. Civo labels derived from the resource identity are applied automatically.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `volume_id` | Unique identifier of the created Civo volume | Volume attachment, monitoring dashboards |
| `attached_instance_id` | ID of the instance the volume is attached to | Populated after attachment via a separate resource |
| `device_path` | Device path of the volume on the attached instance | Populated after attachment via a separate resource |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**General-purpose ext4 volume** -- 50 GiB volume intended for application data, logs, and file storage. ext4 is the most widely compatible Linux filesystem and a safe default for most workloads. The volume requires manual formatting after attachment since the Civo provider does not expose pre-formatting. Start from the **General-Purpose ext4** preset.

**XFS database volume** -- 100 GiB volume optimized for database workloads such as PostgreSQL and MySQL data directories. XFS excels at large sequential writes and parallel I/O. Like the ext4 preset, manual formatting is required after attachment. Start from the **XFS Database** preset.

## Works With

This component operates independently and does not reference other deployment components via foreign keys.
# Storage Volume on DigitalOcean

Deploys a DigitalOcean block storage volume with configurable size, region, and optional filesystem pre-formatting. Volumes provide persistent, network-attached storage that can be attached to Droplets independently of Droplet lifecycle. Integrates with Planton's Provider Connections for DigitalOcean credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Block Storage Volume** -- a `digitalocean_volume` resource in the specified region with the given size, optional filesystem formatting, optional description, and tags
- **DigitalOcean Tags** -- tags from the spec applied directly to the volume resource for organizational tracking and cost allocation

The volume is created in a detached state. Attach it to a Droplet separately via a volume attachment or the DigitalOcean control panel.

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A target region** -- the volume must be created in the same region as any Droplet that will attach to it. Volumes cannot cross regions. Check available regions via the DigitalOcean API or control panel.
- **Sufficient block storage quota** in the target region for the requested volume size.
- **A volume snapshot ID** if creating the volume from an existing snapshot (optional). The snapshot must exist in the same region, and `sizeGib` must be at least the snapshot's size.

## Deploy

### Console

Open the deployment store, find **Storage Volume on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **General-Purpose ext4** preset in the [Presets](#presets) tab for a ready-to-use storage volume.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanVolume
metadata:
  name: app-data
  org: acme-corp
  env: prod
spec:
  volumeName: app-data
  region: nyc1
  sizeGib: 50
```

```shell
planton apply -f do-volume.yaml
```

This creates a 50 GiB unformatted block storage volume in DigitalOcean's NYC1 region. No filesystem formatting, snapshot source, or tags are configured. A Stack Job tracks the provisioning in real time.

Attach the volume to a Droplet after creation -- volumes are provisioned in a detached state by default.

## Key Configuration

These are the most important decisions when configuring a DigitalOcean storage volume. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Volume name** -- The `volumeName` field sets the volume's name on DigitalOcean. Must be lowercase letters, numbers, and hyphens only, starting with a letter and ending with a letter or number (maximum 64 characters).

**Volume size** -- The `sizeGib` field sets the volume capacity in GiB, from 1 to 16,000. DigitalOcean volumes can be resized up to 16 TiB without detaching the volume from its Droplet, so start with a size that matches current needs and scale as data grows.

**Region** -- The `region` field accepts a DigitalOcean region slug (e.g., `nyc1`, `sfo3`, `fra1`, `lon1`). The volume must be in the same region as any Droplet you plan to attach it to. Volumes cannot be moved across regions after creation.

**Filesystem type** -- The `filesystemType` field optionally pre-formats the volume with `ext4` or `xfs`. Pre-formatting eliminates the need for manual `mkfs` after attaching to a Droplet -- the volume is immediately mountable. Use `ext4` for general-purpose storage and `xfs` for database workloads with heavy sequential writes. Leave unset for unformatted volumes.

**Description** -- The `description` field adds a human-readable description (up to 100 characters) stored on the DigitalOcean resource. Useful for documenting the volume's purpose, especially when managing multiple volumes in the same project.

**Snapshot source** -- The `snapshotId` field creates the volume from an existing volume snapshot. The volume inherits the snapshot's data, and `sizeGib` must be at least the snapshot's size. Useful for disaster recovery or cloning environments from production data.

**Tags** -- The `tags` field applies organizational tags to the volume for cost allocation and filtering. Tags are applied directly to the DigitalOcean resource and visible in the control panel and API responses.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `volume_id` | UUID of the created DigitalOcean volume | Volume attachment, monitoring dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**General-purpose ext4 volume** -- 50 GiB volume pre-formatted with ext4 for application data, logs, and file storage. ext4 is the most widely compatible Linux filesystem and a safe default for most workloads. The volume is immediately mountable after attaching to a Droplet. Start from the **General-Purpose ext4** preset.

**Database XFS volume** -- 100 GiB volume pre-formatted with XFS, optimized for database workloads such as PostgreSQL and MySQL data directories. XFS provides superior write performance for sequential and random I/O patterns. The larger starting size accommodates typical database growth. Start from the **Database XFS** preset.

## Works With

This component operates independently and does not reference other deployment components via foreign keys.
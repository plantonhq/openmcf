# Scaleway Block Volume

Deploys a network-attached NVMe block storage volume on Scaleway with configurable performance tier (5K or 15K IOPS) and size. Block volumes persist independently of Instance lifecycle and can be moved between Instances in the same Availability Zone. Supports snapshot-based provisioning for cloning data from existing volumes.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Block Volume** -- a `scaleway_block_volume` in the specified Availability Zone with the configured performance tier (SBS 5K or SBS 15K) and size
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair (Access Key + Secret Key). The IaC module authenticates through the Scaleway provider configuration.
- **Choose an Availability Zone** -- block volumes are zonal resources (`fr-par-1`, `nl-ams-1`, `pl-waw-2`). The volume can only attach to Instances in the same zone. Cannot be changed after creation.
- **Snapshot ID** (optional) -- if creating a volume from a snapshot, the snapshot must exist in the same zone and the volume `sizeGb` must be >= the snapshot's source volume size.

## Deploy

### Console

Open the deployment store, find **Scaleway Block Volume**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a 20 GB volume with the 5K IOPS tier.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayBlockVolume
metadata:
  name: app-data
  org: acme-corp
  env: prod
spec:
  zone: fr-par-1
  sizeGb: 50
  performanceTier: sbs_5k
```

```shell
planton apply -f scaleway-block-volume.yaml
```

This creates a 50 GB block volume with standard 5K IOPS in the Paris-1 zone. No snapshot source is configured. The volume is a raw block device -- format it (`mkfs.ext4`, `mkfs.xfs`) and mount it after attaching to an Instance. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a block volume. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Performance tier** -- The `performanceTier` field selects between `sbs_5k` (5,000 IOPS) for general-purpose workloads and `sbs_15k` (15,000 IOPS) for databases and latency-sensitive applications. The tier can be changed in-place after creation without recreating the volume. Instances using `sbs_15k` volumes need at least 3 GiB/s of block bandwidth to utilize the full IOPS capacity.

**Size** -- The `sizeGb` field accepts values from 5 to 10,240 GB. Volumes can be increased in-place (hot resize) but cannot be shrunk -- the Scaleway provider rejects any plan that decreases size. After increasing via IaC, grow the filesystem inside the OS using `growpart` and `resize2fs` or `xfs_growfs`.

**Zone** -- The `zone` field sets the Availability Zone (`fr-par-1`, `nl-ams-1`, `pl-waw-2`). Cannot be changed after creation. Must match the zone of the Instance the volume will attach to.

**Snapshot source** -- Set `snapshotId` to clone an existing Block Storage snapshot into a new volume. The snapshot must be in the same zone. When omitted, a blank volume is created.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `volume_id` | Zoned ID of the created volume (`{zone}/{uuid}`) | Instance additional volume attachment, monitoring dashboards |
| `volume_name` | Name of the volume in Scaleway Block Storage | Observability labels, alert routing |
| `zone` | Availability Zone where the volume is deployed | Zone co-location verification for Instance attachment |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard storage** -- A 20 GB volume with 5,000 IOPS for general-purpose persistent storage. Suitable for application data, logs, uploads, and development environments. Start from the **Standard** preset.

**High-performance storage** -- A 100 GB volume with 15,000 IOPS for databases, search engines, and I/O-intensive workloads that require low-latency access to persistent data. Start from the **High-Performance** preset.

## Works With

This component operates independently and does not reference other components.
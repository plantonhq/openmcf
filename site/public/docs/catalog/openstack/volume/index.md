---
title: "Volume"
description: "Volume deployment documentation"
icon: "package"
order: 100
componentName: "openstackvolume"
---

# OpenStack Volume

Deploys an OpenStack Cinder block storage volume with configurable size, volume type, and availability zone. Volumes can be created blank, initialized from a Glance image for bootable root disks, restored from a snapshot, or cloned from an existing volume. Supports ValueFromRef for wiring image dependencies in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cinder Block Storage Volume** -- an `openstack_blockstorage_volume_v3` resource with the specified size, optional volume type, and availability zone. Source options (image, snapshot, clone) are mutually exclusive.
- **Volume Metadata** -- key-value pairs stored on the volume, visible in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **A Cinder volume type** available in the target project if specifying `volumeType` (e.g., `SSD`, `HDD`, `ceph-ssd`). If omitted, Cinder uses the project's default volume type.
- **A Glance image UUID** if creating a bootable volume via `imageId`. The image can be referenced via ValueFromRef from an OpenStackImage Cloud Resource.
- **A Cinder snapshot UUID** if restoring from a snapshot via `snapshotId`.
- **An existing volume UUID** if cloning via `sourceVolId`.

## Deploy

### Console

Open the deployment store, find **OpenStack Volume**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Blank Data** preset in the [Presets](#presets) tab for an empty data volume.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackVolume
metadata:
  name: app-data
  org: acme-corp
  env: prod
spec:
  size: 100
  description: Application data volume
```

```shell
planton apply -f openstack-volume.yaml
```

This creates a blank 100 GB Cinder volume using the project's default volume type and availability zone. No source image, snapshot, or clone is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the volume to a Glance image deployed in the same InfraPipeline:

```yaml
spec:
  imageId:
    valueFrom:
      kind: OpenStackImage
      name: ubuntu-base
      fieldPath: status.outputs.image_id
```

The InfraPipeline resolves the dependency graph, deploys the image first, then provisions the bootable volume from it.

## Key Configuration

These are the most important decisions when configuring an OpenStack volume. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Volume size** -- The `size` field sets the capacity in gigabytes. Required for blank and image-based volumes. For snapshot or clone sources, if omitted or set to 0, the source size is used. If provided, it must be greater than or equal to the source size.

**Volume type** -- The `volumeType` field selects the Cinder backend storage class (e.g., `SSD`, `HDD`, `ceph-ssd`, `lvm`). If omitted, the project default is used. Changing the type on an existing volume triggers a retype operation.

**Volume source** -- The `snapshotId`, `sourceVolId`, and `imageId` fields are mutually exclusive. Use `imageId` for bootable root disks, `snapshotId` for disaster recovery or environment cloning, and `sourceVolId` for direct volume clones. All source fields are immutable after creation.

**Availability zone** -- The `availabilityZone` field selects where the volume is created. Must match the instance's availability zone for attachment to work in most OpenStack deployments. Changing the AZ requires recreating the volume.

**Metadata** -- The `metadata` field accepts key-value pairs stored on the volume. Use for tagging, billing, or organizational purposes. Metadata is visible in the OpenStack API and Horizon dashboard.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackImage** (optional) | `imageId` | `status.outputs.image_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `volume_id` | UUID of the created Cinder volume | OpenStackVolumeAttach, instance block device mapping |
| `name` | Name of the volume, derived from metadata.name | Audit logs, operational dashboards |
| `size` | Provisioned size of the volume in gigabytes | Capacity monitoring |
| `volume_type` | Cinder volume type (backend storage class) | Storage tier verification |
| `availability_zone` | Availability zone where the volume was created | Instance co-location planning |
| `region` | OpenStack region where the volume was created | Multi-region deployment coordination |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Blank data volume** -- Empty Cinder volume for application data, database storage, or log aggregation. Attach to an instance via OpenStackVolumeAttach, then partition and format from within the instance. Start from the **Blank Data** preset.

**Bootable volume from image** -- Volume pre-populated with a Glance image for use as a persistent root disk. Create the volume first, then reference it in an instance's block device mapping or attach via OpenStackVolumeAttach. Start from the **Bootable from Image** preset.

## Works With

- [**OpenStack Image**](/cloud-catalog/openstack-image) -- provides Glance images that can be referenced via `imageId` for bootable volumes
---
title: "Volume"
description: "Volume deployment documentation"
icon: "package"
order: 100
componentName: "hetznercloudvolume"
---

# Hetzner Cloud Volume

Deploys a network-attached block storage volume on Hetzner Cloud with optional server attachment and filesystem formatting. Volumes persist independently of servers -- detaching or deleting a server does not affect the volume's data. Size ranges from 10 GB to 10 TB and can be increased after creation but never decreased. Ideal for databases, application state, and any data that must survive server replacement.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Volume** -- an `hcloud_volume` resource with the specified size, location, and optional filesystem format
- **Volume Attachment** (optional) -- an `hcloud_volume_attachment` resource connecting the volume to a server when `serverId` is specified

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **Location selection** -- choose a location (fsn1, nbg1, hel1, ash, hil, sin). The volume can only be attached to servers in the same location.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Volume**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudVolume
metadata:
  name: db-data
  org: acme-corp
  env: prod
spec:
  size: 50
  location: fsn1
  format: ext4
  serverId:
    value: "12345678"
  automount: true
```

```shell
planton apply -f hetznercloud-volume.yaml
```

This creates a 50 GB ext4-formatted volume in Falkenstein, attaches it to the specified server, and auto-mounts it. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a server environment, use ValueFromRef to attach the volume to a server:

```yaml
spec:
  size: 100
  location: fsn1
  format: ext4
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: db-server
      fieldPath: status.outputs.server_id
  automount: true
```

The InfraPipeline resolves the dependency graph, creates the server first, then provisions and attaches the volume.

## Key Configuration

These are the most important decisions when configuring a volume. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Size** -- The `size` field sets the volume capacity in GB (10-10240). Size can be increased after creation but never decreased -- the Hetzner Cloud API rejects size reductions.

**Location** -- The `location` field determines the physical datacenter. Must match the server's location for attachment. Changing forces replacement (data loss).

**Format** -- The `format` field selects the initial filesystem (ext4 or xfs). If unset, the volume is created raw and must be formatted manually. This is a create-time-only setting.

**Server attachment** -- The `serverId` field optionally attaches the volume to a server. Accepts a literal server ID or a ValueFromRef reference. If omitted, the volume is created unattached.

**Automount** -- The `automount` field automatically mounts the volume on the server after attachment. Only meaningful when `serverId` is set.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **HetznerCloudServer** (optional) | `serverId` | `status.outputs.server_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `volume_id` | Hetzner Cloud numeric ID of the volume | API operations, resource tracking |
| `linux_device` | Linux device path on the attached server | OS-level mount commands |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Database storage** -- A 50-100 GB ext4 volume attached to a database server with automount enabled. The volume survives server replacement, preserving database files.

**Unattached reservation** -- Create a volume without server attachment for later use. Useful for pre-provisioning storage before the target server exists.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- volumes are attached to servers for persistent block storage

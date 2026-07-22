# GCP Filestore Instance

Deploys a Google Cloud Filestore instance (`google_filestore_instance`) — fully managed, high-performance NFS file storage for workloads that need a shared POSIX filesystem: media rendering, EDA, genomics, web serving, and shared volumes for GKE and Compute Engine.

## Overview

The spec covers the full instance lifecycle surface:

- **Tiers** — all eight service tiers, from cost-effective HDD (`STANDARD`/`BASIC_HDD`) through mid-tier SSD (`PREMIUM`/`BASIC_SSD`), legacy `HIGH_SCALE_SSD`, modern `ZONAL` with IOPS tuning, and the multi-zone `REGIONAL`/`ENTERPRISE` tiers.
- **File share** — exactly one per instance (a GCP constraint the spec models honestly), with capacity, NFS export access controls (IP allowlists, read-only exports, root squash), and create-time restore from an existing backup (`sourceBackup`).
- **Networking** — one VPC attachment per instance with three connect modes: `DIRECT_PEERING` (default), `PRIVATE_SERVICE_ACCESS` (required for Shared VPC; rides an existing service-networking connection), and `PRIVATE_SERVICE_CONNECT`. Optional reserved `/29` range and IP version modes.
- **Performance** — fixed IOPS or IOPS-per-TiB provisioning on `ZONAL`/`REGIONAL`/`ENTERPRISE` tiers.
- **Protection and encryption** — `deletionProtectionEnabled` as the destroy guard, CMEK via a `GcpKmsKey` reference, NFSv4.1 on the four modern tiers.
- **Replication** — create-time cross-instance replication (`initialReplication`): this instance joins as the STANDBY replica of an existing ACTIVE peer, referenced as another `GcpFilestoreInstance`.
- **Organization** — user `labels` merged beneath Planton's attribution labels, and Resource Manager `tags` for org-policy and IAM conditions.

`instanceName` (falls back to `metadata.name` when omitted), `location`, `tier`, `protocol`, the network configuration, `kmsKeyName`, and `initialReplication` are immutable — changing them replaces the instance and its data. File share capacity grows in place but never shrinks.

## When to Use

- **Shared volumes for GKE** — ReadWriteMany PersistentVolumes via the Filestore CSI driver
- **Lift-and-shift NFS workloads** — applications that expect a POSIX filesystem, not an object API
- **High-throughput pipelines** — rendering, EDA, genomics, ML training data staging
- **Home directories and content management** — shared access across many VMs

## Prerequisites

- GCP credentials with Filestore admin permissions on the target project (the Filestore API is enabled automatically)
- A VPC network for the instance to attach to (reference a `GcpVpcNetwork` or name one literally)
- For CMEK: a `GcpKmsKey` the Filestore service agent can use
- For `PRIVATE_SERVICE_ACCESS`: an existing service-networking connection on the VPC

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpFilestoreInstance
metadata:
  name: my-nfs
spec:
  projectId:
    value: my-gcp-project
  location: us-central1-a
  tier: BASIC_SSD
  fileShare:
    name: vol1
    capacityGb: 2560
  networkConfig:
    network:
      value: default
```

This creates a 2.5 TiB SSD-backed NFS server named `my-nfs` (from `metadata.name`) in `us-central1-a`, mountable from any client on the default VPC. Note the `projectId` shape: a `StringValueOrRef` takes `value:` for a literal or `valueFrom:` for a reference to a `GcpProject` — never a bare scalar.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Target project; reference a `GcpProject` or set a literal value. Empty falls back to the provider's default project. Immutable |
| `instanceName` | string | No | GCP resource name (2-63 lowercase chars); empty falls back to `metadata.name`. Immutable |
| `location` | string | Yes | Zone for zonal tiers (`us-central1-a`), region for `REGIONAL`/`ENTERPRISE` (`us-central1`). Immutable |
| `tier` | string | Yes | `STANDARD`, `PREMIUM`, `BASIC_HDD`, `BASIC_SSD`, `HIGH_SCALE_SSD`, `ZONAL`, `REGIONAL`, `ENTERPRISE`. Immutable |
| `description` | string | No | Human-readable description |
| `protocol` | string | No | `NFS_V3` (default) or `NFS_V4_1` (HIGH_SCALE_SSD/ZONAL/REGIONAL/ENTERPRISE only). Immutable |
| `kmsKeyName` | StringValueOrRef | No | CMEK key (reference a `GcpKmsKey`); empty uses Google-managed keys. Immutable |
| `deletionProtectionEnabled` | bool | No | The destroy guard: a protected instance cannot be deleted until this is flipped false |
| `deletionProtectionReason` | string | No | Informational reason for the protection |
| `fileShare.name` | string | Yes | Export path name (letter + up to 15 letters/numbers/underscores). Immutable |
| `fileShare.capacityGb` | int | Yes | Capacity in GiB (min 1024; 2560 for BASIC_SSD/PREMIUM, 10240 for HIGH_SCALE_SSD). Grows, never shrinks |
| `fileShare.nfsExportOptions[]` | list | No | Up to 10: `ipRanges` (max 64 total), `accessMode` (`READ_WRITE` default / `READ_ONLY`), `squashMode` (`NO_ROOT_SQUASH` default / `ROOT_SQUASH`), `anonUid`/`anonGid` (default 65534). Empty allows all clients read-write |
| `fileShare.sourceBackup` | string | No | Restore from `projects/{p}/locations/{l}/backups/{b}`; capacity must cover the backup's source. Create-time only |
| `networkConfig.network` | StringValueOrRef | Yes | VPC network (reference a `GcpVpcNetwork`). Immutable |
| `networkConfig.connectMode` | string | No | `DIRECT_PEERING` (default), `PRIVATE_SERVICE_ACCESS` (Shared VPC), `PRIVATE_SERVICE_CONNECT`. Immutable |
| `networkConfig.reservedIpRange` | string | No | A `/29` block for the instance; empty lets GCP pick. Immutable |
| `networkConfig.modes[]` | list | No | `MODE_IPV4` and/or `MODE_IPV6`; empty means `["MODE_IPV4"]`. Immutable |
| `performanceConfig.fixedIops.maxIops` | int | No | Constant IOPS (multiple of 1000). Mutually exclusive with `iopsPerTb` |
| `performanceConfig.iopsPerTb.maxIopsPerTb` | int | No | IOPS scaling with capacity. Mutually exclusive with `fixedIops` |
| `initialReplication.role` | string | No | This instance's role: `STANDBY` (default — the receiving replica) or `ACTIVE` |
| `initialReplication.peerInstances[]` | list | Yes* | Peer instances (reference `GcpFilestoreInstance` or a literal full resource path). *Required when the block is set. Create-time only |
| `labels` | map | No | User labels, merged beneath platform attribution labels (platform wins on conflicts) |
| `tags` | map | No | Resource Manager tags (`tagKeys/{id}` → `tagValues/{id}`). Create-time only |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `instance_id` | Fully qualified ID (`projects/{p}/locations/{l}/instances/{i}`) — what peers reference |
| `instance_name` | Short instance name |
| `ip_addresses` | Addresses on the VPC; mount with `mount <ip>:/<share> /mnt/nfs` |
| `file_share_name` | Share name for the mount path |
| `create_time` | Creation timestamp (RFC3339) |
| `reserved_ip_range` | The `/29` block as resolved by GCP (populated even when auto-picked) |
| `etag` | Server-specified ETag guarding concurrent updates |

See the [presets](presets/) for remixable starting points and [docs/README.md](docs/README.md) for the deep dive.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

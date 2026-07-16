# GCP Filestore Instance

Creates a Google Cloud Filestore instance — fully managed, high-performance NFS file storage for shared POSIX filesystems — with a single file share, VPC connectivity across three connect modes, NFS export access controls, IOPS tuning, CMEK encryption, create-time replication and backup restore, and deletion protection. Supports all eight Filestore tiers from cost-effective HDD to enterprise-grade regional HA.

## What Gets Created

- The Filestore API is enabled on the project (never disabled on destroy)
- A `google_filestore_instance` carrying your labels merged beneath Planton's attribution labels (`planton-ai_resource`, `planton-ai_name`, `planton-ai_kind`, plus org/env/id when set) and any Resource Manager tags
- One NFS file share on the instance (a GCP constraint: exactly one per instance), mountable at `<ip>:/<share_name>`
- One VPC network attachment via `DIRECT_PEERING` (default), `PRIVATE_SERVICE_ACCESS` (Shared VPC; rides an existing service-networking connection), or `PRIVATE_SERVICE_CONNECT`

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **IAM permissions** — Filestore admin on the target project
- **A VPC network** for the instance to attach to (reference a `GcpVpcNetwork`)
- For CMEK: a `GcpKmsKey` the Filestore service agent can use
- For `PRIVATE_SERVICE_ACCESS`: an existing service-networking connection on the VPC

## Quick Start

Create a file `filestore.yaml`:

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

Deploy:

```shell
planton apply -f filestore.yaml
```

This creates a 2.5 TiB SSD-backed NFS server named `my-nfs` (the instance name falls back to `metadata.name`) in `us-central1-a`, mountable from any client on the default VPC: `mount <ip>:/vol1 /mnt/nfs`.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Target project; reference a `GcpProject` or set a literal value. Empty falls back to the provider's default project. Immutable |
| `instanceName` | string | No | GCP resource name; empty falls back to `metadata.name`. Immutable |
| `location` | string | Yes | Zone for zonal tiers, region for `REGIONAL`/`ENTERPRISE`. Immutable |
| `tier` | string | Yes | `STANDARD`, `PREMIUM`, `BASIC_HDD`, `BASIC_SSD`, `HIGH_SCALE_SSD`, `ZONAL`, `REGIONAL`, `ENTERPRISE`. Immutable |
| `description` | string | No | Human-readable description |
| `protocol` | string | No | `NFS_V3` (default) or `NFS_V4_1` (modern tiers only). Immutable |
| `kmsKeyName` | StringValueOrRef | No | CMEK key (reference a `GcpKmsKey`). Immutable |
| `deletionProtectionEnabled` | bool | No | Destroy guard: must be flipped false before a protected instance can be deleted |
| `deletionProtectionReason` | string | No | Informational reason for the protection |
| `fileShare.name` | string | Yes | Export path name (max 16 chars). Immutable |
| `fileShare.capacityGb` | int | Yes | Capacity in GiB (tier-specific minimums). Grows, never shrinks |
| `fileShare.nfsExportOptions[]` | list | No | Up to 10 access controls: `ipRanges`, `accessMode`, `squashMode`, `anonUid`/`anonGid` |
| `fileShare.sourceBackup` | string | No | Restore from an existing Filestore backup. Create-time only |
| `networkConfig.network` | StringValueOrRef | Yes | VPC network (reference a `GcpVpcNetwork`). Immutable |
| `networkConfig.connectMode` | string | No | `DIRECT_PEERING` (default), `PRIVATE_SERVICE_ACCESS`, `PRIVATE_SERVICE_CONNECT`. Immutable |
| `networkConfig.reservedIpRange` | string | No | A `/29` block; empty lets GCP pick. Immutable |
| `networkConfig.modes[]` | list | No | `MODE_IPV4` / `MODE_IPV6`; empty means IPv4. Immutable |
| `performanceConfig` | object | No | `fixedIops.maxIops` or `iopsPerTb.maxIopsPerTb` (mutually exclusive; ZONAL/REGIONAL/ENTERPRISE) |
| `initialReplication` | object | No | Create-time replication: `role` (`STANDBY` default / `ACTIVE`) + `peerInstances[]` referencing other `GcpFilestoreInstance` resources |
| `labels` | map | No | User labels (merged beneath platform labels) |
| `tags` | map | No | Resource Manager tags (`tagKeys/{id}` → `tagValues/{id}`). Create-time only |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `instance_id` | Fully qualified resource ID — what replication peers reference |
| `instance_name` | Short instance name |
| `ip_addresses` | Addresses on the VPC network; use the first for NFS mounts |
| `file_share_name` | Share name for the mount path |
| `create_time` | Creation timestamp (RFC3339) |
| `reserved_ip_range` | The `/29` block as resolved by GCP |
| `etag` | Server-specified ETag guarding concurrent updates |

## Related Resources

- **GcpVpcNetwork** — the VPC the instance attaches to via `networkConfig.network`
- **GcpKmsKey / GcpKmsKeyRing** — CMEK protection via `kmsKeyName`
- **GcpProject** — target project via `projectId`
- **GcpGkeCluster** — mount shares as ReadWriteMany PersistentVolumes via the Filestore CSI driver
- **GcpComputeInstance** — VMs mount shares directly over NFS
- **GcpFilestoreInstance** — another instance as the ACTIVE replication peer

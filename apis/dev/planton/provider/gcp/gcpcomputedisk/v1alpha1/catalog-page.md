# GCP Compute Disk

Creates a zonal Google Compute Engine persistent disk — the durable block device behind stateful VMs — with source control (empty, image, snapshot, or clone of another disk), the full pd and hyperdisk type families with in-place IOPS/throughput tuning, customer-managed encryption, and a snapshot-before-destroy recovery net.

## What Gets Created

- The Compute Engine API is enabled on the project (never disabled on destroy)
- A `google_compute_disk` carrying your labels merged beneath Planton's attribution labels (`planton-ai_resource`, `planton-ai_name`, `planton-ai_kind`, plus org/env/id when set)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **IAM permissions** — `roles/compute.storageAdmin` on the target project
- For CMEK: a `GcpKmsKey` the Compute Engine service agent can use (`roles/cloudkms.cryptoKeyEncrypterDecrypter`)

## Quick Start

Create a file `disk.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpComputeDisk
metadata:
  name: app-data
spec:
  zone: us-central1-a
  type: pd-balanced
  sizeGb: 100
  labels:
    team: platform
```

Deploy:

```shell
planton apply -f disk.yaml
```

This creates an empty 100 GB pd-balanced data disk whose data survives any VM it is later attached to — attach it via a `GcpComputeInstance`'s `attachedDisks[].source`.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | Owning project; empty uses the provider default. Immutable |
| `diskName` | string | No | Disk name; falls back to `metadata.name`. Immutable |
| `zone` | string | Yes | Zone, e.g. `us-central1-a`. Immutable |
| `description` | string | No | Human-readable description |
| `type` | string | No | `pd-standard` / `pd-balanced` / `pd-ssd` / `pd-extreme` / `hyperdisk-*`. Immutable |
| `sizeGb` | int | Empty disks | Size in GB; grows in place, never shrinks |
| `image` | string | No | Source image (family or self link) — bootable disk. Create-time only |
| `sourceSnapshot` | string | No | Snapshot to restore from. Create-time only |
| `sourceDisk` | StringValueOrRef | No | Existing `GcpComputeDisk` to clone. Create-time only |
| `kmsKey` | StringValueOrRef | No | CMEK key (reference a `GcpKmsKey`). Immutable |
| `provisionedIops` | int | pd-extreme, hyperdisk-extreme | IOPS; in-place tunable on hyperdisk (every 4h at most) |
| `provisionedThroughput` | int | No | MB/s on hyperdisk types; in-place tunable (every 4h at most) |
| `accessMode` | string | No | `READ_WRITE_SINGLE` (default), `READ_WRITE_MANY`, `READ_ONLY_MANY` |
| `architecture` | string | No | `X86_64` or `ARM64`. Immutable |
| `enableConfidentialCompute` | bool | No | Confidential-compute disk (hyperdisk SKUs; requires `kmsKey`) |
| `physicalBlockSizeBytes` | int | No | 4096 (default) or 16384 |
| `createSnapshotBeforeDestroy` | bool | No | Final snapshot at destroy — last-resort recovery net |
| `snapshotBeforeDestroyPrefix` | string | No | Name prefix for that final snapshot |
| `storagePool` | string | No | Hyperdisk storage pool |
| `labels` | map | No | User labels (merged beneath platform labels) |
| `resourceManagerTags` | map | No | `tagKeys/{id}` → `tagValues/{id}` (create-time only) |

At most one source (`image` / `sourceSnapshot` / `sourceDisk`); none creates an empty disk, which then requires `sizeGb`.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Name of the disk in GCP |
| `disk_id` | Server-assigned unique numeric identifier |
| `self_link` | Self-link URL — what instance attachments consume |
| `zone` | Zone (plain zone name) |
| `size_gb` | Provisioned size in GB |
| `type` | Disk type (plain type name) |

## Related Resources

- **GcpComputeInstance** — attach this disk via `attachedDisks[].source` or boot from it via `bootDisk.sourceDisk` (both reference the `self_link` output)
- **GcpComputeDisk** — clone this disk into another via `sourceDisk` (self-reference by kind)
- **GcpKmsKey** — CMEK protection via `kmsKey`
- **GcpProject** — the owning project via `projectId`

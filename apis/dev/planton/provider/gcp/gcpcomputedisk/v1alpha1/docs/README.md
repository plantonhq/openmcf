# GcpComputeDisk — Deep Dive

## The problem this resource solves

A VM's most valuable property is usually not the VM — it is the data on its disks. When the disk lives inside the instance definition, replacing the instance (a machine-type change, an image upgrade, a bad rollout) puts the data's fate in the hands of whatever the instance's deletion flags happen to say. This kind models the zonal persistent disk as a first-class node with an independent lifecycle: create the disk once, attach it to a `GcpComputeInstance` by reference, and the data survives instance replacement, resizing, and rescheduling. The dangerous decisions — encryption, sources, the destroy-time snapshot — are deliberate spec choices reviewed like any other resource.

## Where it sits in the composition

The disk's `self_link` output is the composition key:

- **GcpComputeInstance `attachedDisks[].source`** — the standard data-volume pattern: the instance mounts the disk, the disk outlives the instance. A CMEK-encrypted disk's attachment must also present the same `kmsKey` the disk was created with.
- **GcpComputeInstance `bootDisk.sourceDisk`** — boot from a pre-created bootable disk (one built from an `image`). Size/type/encryption on the instance side are ignored; the disk already owns them.
- **GcpComputeDisk `sourceDisk`** — self-clone: stamp out copies of a golden disk by referencing another `GcpComputeDisk` of this same kind.

Inbound references the disk itself makes: `projectId` → `GcpProject`, `kmsKey` → `GcpKmsKey`, `sourceDisk` → another `GcpComputeDisk`.

Because a disk is zonal, it attaches only to instances in the same zone — pick the zone once, deliberately.

## Disk types

The pd family:

- **`pd-standard`** — HDD-backed; cheapest per GB, lowest performance. Cold data, batch scratch.
- **`pd-balanced`** — SSD-backed, GCP's default and the sensible general choice; performance scales with size.
- **`pd-ssd`** — higher IOPS per GB for latency-sensitive databases.
- **`pd-extreme`** — provisioned-IOPS pd; `provisionedIops` is required and fixed at create time.

The hyperdisk family (`hyperdisk-balanced`, `hyperdisk-extreme`, `hyperdisk-throughput`, `hyperdisk-ml`, on supported machine families) decouples capacity from performance: `provisionedIops` and `provisionedThroughput` are purchased independently of size and — the key operational property — **update in place, at most every 4 hours**. A pd disk's performance is a function of its size; a hyperdisk's performance is a dial you turn without touching the data. `hyperdisk-ml` additionally supports `accessMode: READ_WRITE_MANY`/`READ_ONLY_MANY` for shared datasets across many VMs.

## Sources: at most one, or none

Exactly four ways to birth a disk, and the spec enforces choosing at most one:

- **No source (empty)** — the common case for data volumes. `sizeGb` is then required (there is nothing to infer a size from).
- **`image`** — initializes the disk from an image family (`debian-cloud/debian-12`) or a specific image self link, making it bootable. The starting point for golden images.
- **`sourceSnapshot`** — restores from a snapshot (name or self link). The recovery and cloning-across-time path.
- **`sourceDisk`** — clones another `GcpComputeDisk` directly, referenced by kind. The fan-out path for golden disks.

All sources are create-time only. With a source, `sizeGb` is optional — the source's size is used, and a larger value grows the disk at birth.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `diskName`, `zone`, `type`, `image`, `sourceSnapshot`, `sourceDisk`, `kmsKey`, `architecture`, `projectId` | Immutable (ForceNew) — replacement destroys the disk and its data |
| `sizeGb` | Grows in place; shrinking is impossible (a smaller value forces replacement) |
| `provisionedIops` / `provisionedThroughput` | In-place on hyperdisk types (at most every 4 hours); create-time on pd-extreme; unsettable on other pd types |
| `labels`, `description` | Mutable in place |
| Deletion | Fails while the disk is attached to a running instance — detach first (or delete the instance). With `createSnapshotBeforeDestroy: true`, a final snapshot is taken before the disk goes |

`diskName` falls back to `metadata.name` when omitted — both IaC engines derive the identical cloud-side name.

## Encryption (CMEK)

Setting `kmsKey` (a `GcpKmsKey` reference, resolving to `status.outputs.key_id`) switches the disk from Google-managed to customer-managed encryption. One prerequisite trips everyone once: the **Compute Engine service agent** (`service-<project-number>@compute-system.iam.gserviceaccount.com`) — not your deploy credential — must hold `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key before create, or the disk fails to provision. The key choice is immutable; snapshots taken by `createSnapshotBeforeDestroy` reuse the disk's key. Confidential-compute mode (`enableConfidentialCompute`, hyperdisk SKUs only) requires CMEK, and the spec enforces that pairing at validation time.

## The safety posture

- **Deleting an attached disk fails.** GCP refuses; this is a feature, not a bug — a disk in active use cannot be yanked out from under its instance by a stack teardown.
- **`createSnapshotBeforeDestroy`** is the last-resort recovery net for precious volumes: even if a destroy is a mistake, the data survives as a snapshot named `<disk-name>-YYYYMMDD-HHmmss` (prefix overridable via `snapshotBeforeDestroyPrefix`). It is not a backup strategy — schedule real snapshots for that — it is the net under the trapeze.
- **Grow-only sizing** means a resize can never be a data-loss event.

## Deliberately not modeled

- **Regional (dual-zone replicated) disks** — a separate GCP resource (`google_compute_region_disk`) with a materially different surface: no image/architecture/confidential arms, a different CMEK attribute, replica-zone semantics. Deferred as its own future kind rather than bolted on as a mode flag.
- **Raw CSEK keys** (customer-supplied encryption) — the key is secret material that would live in manifests and state. CMEK via `GcpKmsKey` is the modeled path.
- **Async primary-disk replication pairs** — depends on the regional-disk kind; deferred with it.
- **`guest_os_features` / `licenses`** — image-derived and computed in practice; declaring them by hand invites drift against what the image actually carries.
- **`source_instant_snapshot` / `source_storage_object` import arms** — specialist migration paths, deferred until composition demand pulls them in.

## Best practices

- **CMEK for regulated data.** Anything with compliance obligations gets a `GcpKmsKey` reference — and grant the Compute Engine service agent on the key first.
- **Snapshot-before-destroy for precious volumes.** Any disk whose loss would hurt sets `createSnapshotBeforeDestroy: true`. The cost of an unnecessary snapshot rounds to zero; the cost of the alternative does not.
- **Hyperdisk when performance needs will change.** If you cannot predict IOPS needs today, `hyperdisk-balanced` lets you retune every 4 hours without replacing the disk; with pd types the escape hatch is growing the disk or migrating.
- **Size for growth, not peak.** Disks grow in place at any time and never shrink — start honest and grow when monitoring says so.

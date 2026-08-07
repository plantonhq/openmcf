# High Performance Zonal

The throughput posture: the modern ZONAL tier with IOPS provisioned per
terabyte, so performance scales automatically as the share grows, plus
customer-managed encryption.

## What this preset creates

A ZONAL instance whose `performanceConfig.iopsPerTb` ties IOPS to
capacity — 2 TiB × 4000 IOPS/TiB = 8000 IOPS at creation, rising
automatically with every capacity increase (capacity only ever grows).
Data at rest is encrypted with a customer-managed key referenced from a
`GcpKmsKey` resource.

## Prerequisites

- A `GcpKmsKey` named `storage-cmek` that the Filestore service agent
  can use (replace with yours, or drop `kmsKeyName` for Google-managed
  keys).
- A VPC network named `ml-vpc` (replace, or reference a `GcpVpcNetwork`
  via `valueFrom`).

## iopsPerTb vs. fixedIops

`iopsPerTb` is the set-and-forget model: grow the share, get more IOPS.
`fixedIops` pins a constant number (a multiple of 1000) regardless of
capacity — right when the working set grows but the IOPS demand does
not, so you stop paying for performance you don't use. They are mutually
exclusive; the spec rejects both at once.

## Remix ideas

- Add `protocol: NFS_V4_1` — supported on ZONAL — for NFSv4.1 semantics.
- Add `nfsExportOptions` to restrict mounts to the training subnets.
- Need zone-failure tolerance at this performance level? Move to
  `REGIONAL` with a region `location` — but note tier and location are
  immutable, so that is a new instance, not an edit.

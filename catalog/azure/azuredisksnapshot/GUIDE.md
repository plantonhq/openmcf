# Azure Disk Snapshot -- Operational Guide

Judgment calls that matter when you run disk snapshots in production.

## Incremental unless you have a reason, decided at creation

Incremental snapshots store deltas on standard storage and are what every backup chain and cross-region copy workflow wants; full snapshots store the entire disk and exist mainly for compatibility with older tooling. The mode is ForceNew -- a chain started full stays full -- so set `incrementalEnabled: true` as the reflex and document the exception when you break it. One caveat worth knowing: the first incremental snapshot of a disk stores the full disk; the savings start at the second.

## This kind is an artifact, not a backup strategy

A snapshot is one point-in-time copy with no schedule, no retention, and no restore orchestration. Scheduled VM protection belongs to the Recovery Services kinds (vault + policy + protected VM); Data Protection covers the modern per-datasource backups. Reach for THIS kind when a pipeline or a human wants one deliberate copy: before a risky change, as an image-pipeline handoff, as a clone source. If you find yourself cron-ing snapshot manifests, you wanted a backup policy.

## The pairing rules Azure enforces and the schema does not

The provider deliberately does not tie `createOption` to its source fields -- the spec preserves that honesty, so a manifest with `createOption: Copy` and no `sourceResourceId` validates and then fails at Azure. The working pairs: **Copy** reads `sourceResourceId` (a disk or another snapshot); **Import** reads `sourceUri` + `storageAccountId`. Get the pair right in the manifest; the offline plan cannot catch a missing source, only the live create can.

## Network posture: snapshots are exfiltration surfaces

A snapshot's data plane supports export -- anyone with the right role can generate a download SAS. The public defaults (`AllowAll` + public access on) are fine for build artifacts; for snapshots of sensitive disks set `networkAccessPolicy: AllowPrivate` with a disk-access resource (private endpoint) or `DenyAll` when nothing should ever export it, and turn `publicNetworkAccessEnabled` off alongside. Treat the posture like you treat a storage account's.

## ADE settings ride along, one way

The `encryptionSettings` block exists for sources encrypted with legacy in-guest Azure Disk Encryption: it records where the BitLocker/dm-crypt secrets live so a restored disk can boot. Platform-managed and customer-managed (disk encryption set) encryption need NOTHING here. Once set, the block cannot be removed in place (Azure cannot disable encryption on a snapshot) -- removal replaces the resource.

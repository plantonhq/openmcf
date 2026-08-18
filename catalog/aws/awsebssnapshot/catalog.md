# AWS EBS Snapshot

Point-in-time disk backups as declarative resources: capture a volume, replicate it across regions, re-encrypt it under a new key, or bring a VMDK/VHD/RAW image into AWS — then restore volumes from it anywhere in the region.

## What Gets Managed

- The snapshot, from one of three sources: a live volume, another snapshot (copy), or a disk image in S3 (import).
- Its storage tier (standard or archive) and archive-restore dials.
- Fast snapshot restore per availability zone, and createVolumePermission grants to other accounts.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EC2 permissions.

### AWS Prerequisites

- Volume arm: the source volume.
- Import arm: the disk image staged in S3 and the VM Import/Export service role (`vmimport` by convention) with its documented trust and S3 policy.

## After You Deploy

- Restore volumes from `snapshot_id` (any zone in the region — snapshots are regional).
- Shared accounts create their own volumes from the snapshot without copying it; encrypted snapshots additionally need the KMS key shared.

## Common Changes

- Archive an old snapshot: set `storage_tier: archive` (in place; restore before next use, minimum 24-90 day archive billing).
- Cheap DR: a copy arm instance per region, re-encrypting under the DR region's key.
- Automate the cadence: recurring snapshots belong to AwsDlmLifecyclePolicy, not hand-managed snapshot resources.

# AWS EBS Volume

Block storage as its own resource: provision a disk once — sized, typed, encrypted — and attach it to the instance that needs it today, not the one it happened to be born with. Data outlives compute.

## What Gets Managed

- The volume: zone, type (gp3/io1/io2/...), size, provisioned IOPS and throughput, encryption, multi-attach — or a copy of an existing volume.
- Its attachments: which instances see it, at which device names, and how detach behaves on teardown.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with EC2 permissions.

### AWS Prerequisites

- For attachments: a running instance in the SAME availability zone (EBS never attaches across zones).
- For snapshot restores: the snapshot id (same region).

## After You Deploy

- Reference `volume_id` from snapshots, copies, and attachments.
- Mount the device on the instance (the attachment presents the disk; the filesystem is yours).

## Common Changes

- Grow it: raise `size_gb` — in place, no detach (grow the filesystem after).
- Tune it: `iops`/`throughput_mibps` update in place on gp3/io-family volumes.
- Move it across zones: there is no move — snapshot, then restore in the target zone.

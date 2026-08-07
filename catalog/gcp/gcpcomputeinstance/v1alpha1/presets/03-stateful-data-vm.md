# Stateful Data VM

A database-on-VM pattern where everything durable outlives the instance:
data on a referenced `GcpComputeDisk`, a stable internal address from a
referenced `GcpAddress`, deletion protection on the VM object, and
live-migration for zero-downtime host maintenance.

## What this preset creates

An `n2-standard-4` VM with a `pd-balanced` boot disk and one attached
data disk. The disk is a first-class `GcpComputeDisk` attached by its
`self_link` output — the instance holds only the attachment, so
replacing or deleting the VM never touches the data. The reserved
internal IP keeps client configuration and DNS stable across rebuilds.

## Prerequisites

- A `GcpComputeDisk` named `postgres-data` in the same zone as the
  instance (replace the reference with your own).
- An INTERNAL `GcpAddress` named `data-vm-ip` in the VPC/subnetwork the
  interface attaches to (replace the reference with your own).

## The data-safety split

`deletionProtection: true` guards the *instance object* — deleting it
fails until the flag is flipped. The *data* is guarded by the disk's own
lifecycle: it is not created by this VM and is not deleted with it.
Inside the guest, the disk appears as
`/dev/disk/by-id/google-data` (from `deviceName`).

## Remix ideas

- Attach the same disk `READ_ONLY` to multiple VMs for shared read
  workloads (`mode: READ_ONLY`).
- Add a `resourcePolicies` entry on the boot disk for scheduled
  snapshots, or manage snapshot schedules on the `GcpComputeDisk`.
- Set `desiredStatus: TERMINATED` to stop the VM (compute billing stops;
  the disks and their data remain) and `RUNNING` to start it again.
- Move the interface to a `GcpSubnetwork` reference and add a dedicated
  `GcpServiceAccount`, as in the hardened preset, for production.

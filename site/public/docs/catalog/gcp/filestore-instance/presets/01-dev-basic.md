---
title: "Dev Basic"
description: "The minimal working Filestore instance for development, testing, and CI: SSD-backed, on the default VPC, easy to tear down."
type: "preset"
rank: "01"
presetSlug: "01-dev-basic"
componentSlug: "filestore-instance"
componentTitle: "Filestore Instance"
provider: "gcp"
icon: "package"
order: 1
---

# Dev Basic

The minimal working Filestore instance for development, testing, and CI:
SSD-backed, on the default VPC, easy to tear down.

## What this preset creates

A BASIC_SSD instance at the tier's 2.5 TiB minimum, exporting a single
share `vol1` that any client on the default VPC can mount read-write
(`mount <ip>:/vol1 /mnt/nfs` — the address surfaces in the
`ip_addresses` output). `instanceName` is omitted, so the instance takes
its name from `metadata.name`. No deletion protection, so cleanup is
frictionless.

## Remix ideas

- Add `nfsExportOptions` with `ipRanges` to restrict which subnets can
  mount the share — empty export options allow all clients.
- Switch `tier` to `STANDARD` for HDD-backed storage at lower cost
  (1 TiB minimum) when throughput doesn't matter.
- Grow `capacityGb` any time; it can never shrink, so start at the
  minimum.

## When to graduate

Production workloads want the enterprise preset: a regional tier for
zone-failure tolerance, `deletionProtectionEnabled` as the destroy
guard, and `PRIVATE_SERVICE_ACCESS` networking.

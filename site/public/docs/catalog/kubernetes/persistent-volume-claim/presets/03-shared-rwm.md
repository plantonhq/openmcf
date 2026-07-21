---
title: "Shared RWM"
description: "This preset creates a volume that many pods — across many nodes — mount read-write simultaneously: `ReadWriteMany`. It is the shape for shared asset directories, upload buffers, and legacy..."
type: "preset"
rank: "03"
presetSlug: "03-shared-rwm"
componentSlug: "persistent-volume-claim"
componentTitle: "Persistent Volume Claim"
provider: "kubernetes"
icon: "package"
order: 3
---

# Shared RWM

This preset creates a volume that many pods — across many nodes — mount read-write simultaneously: `ReadWriteMany`. It is the shape for shared asset directories, upload buffers, and legacy applications that expect a common filesystem.

**ReadWriteMany requires a shared-filesystem driver.** Block-storage drivers (EBS, GCE PD, Azure Disk) do not support it — a ReadWriteMany claim against a block-storage class does not error at creation; it simply **never provisions** and pends forever. The class this preset names must be backed by a shared filesystem: EFS (`efs.csi.aws.com`) on AWS, Filestore (`filestore.csi.storage.gke.io`) on GCP, Azure Files (`file.csi.azure.com`) on Azure, or CephFS/NFS on self-managed clusters.

## When to Use

- Several pods (or several workloads) reading and writing the same files
- Shared asset/media directories behind a fleet of web servers
- Legacy applications that coordinate through a common filesystem instead of an API or queue

## Key Configuration Choices

- **`access_modes: [ReadWriteMany]`** — the request; whether it can be honored is entirely a property of the class's driver (see above)
- **`storage_class_name.value`** — must name a shared-filesystem class; the cluster default is almost never one
- **`storage_request: 100Gi`** — note that some shared-filesystem drivers (e.g. EFS) are elastic and treat the size as nominal, while others (Filestore, Azure Files) provision it
- **Consider the consistency model** — a shared filesystem is not a coordination mechanism; concurrent writers to the same file still need application-level discipline

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace the claim lives in — pods mounting it must be in the same namespace | Your namespace management |
| `<your-shared-fs-class>` | A StorageClass backed by a shared-filesystem driver (EFS, Filestore, Azure Files, CephFS, NFS) | `kubectl get storageclass` — check the PROVISIONER column |

The size `100Gi` is a working example — replace it with the actual requirement.

## Related Presets

- **02-pinned-class** — the single-node ReadWriteOnce shape
- **01-dynamic-default-class** — the simplest claim, default class

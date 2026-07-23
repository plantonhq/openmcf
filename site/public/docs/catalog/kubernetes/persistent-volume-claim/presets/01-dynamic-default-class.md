---
title: "Dynamic Default Class"
description: "This preset creates the simplest useful claim: a name and a size, provisioned through the cluster's DEFAULT StorageClass. Omitting `storage_class_name` is itself the request — the claim receives..."
type: "preset"
rank: "01"
presetSlug: "01-dynamic-default-class"
componentSlug: "persistent-volume-claim"
componentTitle: "Persistent Volume Claim"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dynamic Default Class

This preset creates the simplest useful claim: a name and a size, provisioned through the cluster's DEFAULT StorageClass. Omitting `storage_class_name` is itself the request — the claim receives whatever class carries the `storageclass.kubernetes.io/is-default-class` annotation. Access modes default to `["ReadWriteOnce"]` (one node at a time, the mode every block-storage driver supports) and volume mode to `filesystem`.

## When to Use

- Ad-hoc or development storage where the cluster default's characteristics are acceptable
- The first claim on a new cluster, before a platform storage tier exists
- Any claim where "10Gi of whatever the cluster gives me" is the actual requirement

## Key Configuration Choices

- **No `storage_class_name`** — absent means "cluster default". This is deliberately different from an EMPTY class name, which opts out of dynamic provisioning entirely (that shape is `disable_dynamic_provisioning: true`)
- **`storage_request: 10Gi`** — a Kubernetes quantity; growing later requires the default class to allow volume expansion (Kubernetes never shrinks)
- **Defaults do the rest** — `access_modes` falls back to `["ReadWriteOnce"]` and `volume_mode` to `filesystem`, applied identically by both IaC engines
- **Pending until consumed is normal** — if the default class binds on first consumer (the norm for zonal cloud disks and kind's local-path), the claim stays **Pending until a pod uses it — correct behavior, not an error**

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace the claim lives in — workloads must be in the same namespace to mount it | Your namespace management |

The size `10Gi` is a working example — replace it with the actual requirement.

## Related Presets

- **02-pinned-class** — the production shape: name the class instead of inheriting the default
- **03-shared-rwm** — a volume shared by several pods
- **04-clone-from-claim** — start from a copy of an existing claim's data

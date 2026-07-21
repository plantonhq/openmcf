---
title: "Single-Instance Database"
description: "This preset deploys a single-replica database: one pod with a stable name (`my-database-0`), one PersistentVolumeClaim stamped from the `data` template that survives pod restarts and rescheduling,..."
type: "preset"
rank: "01"
presetSlug: "01-database"
componentSlug: "statefulset"
componentTitle: "StatefulSet"
provider: "kubernetes"
icon: "package"
order: 1
---

# Single-Instance Database

This preset deploys a single-replica database: one pod with a stable name (`my-database-0`), one PersistentVolumeClaim stamped from the `data` template that survives pod restarts and rescheduling, and TCP probes so the pod only receives connections once the database is accepting them. The mount is wired by name — the `pvc.claimName` in the container's volume mounts matches the volume claim template's `name`.

Clients connect through the StatefulSet's governing Service (exported as the `service` and `kubeEndpoint` outputs). External exposure, if ever needed, is composed with a first-class ingress kind referencing those outputs — never embedded in the workload.

## When to Use

- Development and staging databases, or small production databases where one instance is acceptable
- Any single-writer stateful system: a relational database, a single-node queue, a cache with persistence

## Key Configuration Choices

- **Single replica (the default)** — no `availability` block needed; the one pod is `<name>-0` and re-attaches to its PVC wherever it reschedules
- **`data` volume claim template, 10Gi, ReadWriteOnce** — each replica gets its own claim; growing later requires a StorageClass with volume expansion enabled
- **Default StorageClass** — `storageClass` is left unset so the cluster's default provisioner serves the claim; pin one explicitly to select SSD-backed or expandable storage
- **TCP probes on the database port** — readiness gates client traffic, liveness restarts a wedged process; swap in your engine's native health command via an `exec` probe if it has one
- **PVC retention defaults to Retain** — deleting the StatefulSet keeps the data; re-creating it with the same name re-adopts the volume

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-database-image>` | Database container image (e.g., `postgres`) | Your container registry |
| `<your-image-tag>` | Image tag or version (e.g., `16.3`) | Your registry or CI/CD pipeline output |

## Related Presets

- **02-ha-quorum-cluster** — Three replicas with anti-affinity and a quorum-guarding disruption budget
- **03-hardened-database** — Restricted-profile security hardening with a composed ServiceAccount identity

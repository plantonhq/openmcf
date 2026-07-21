---
title: "Hardened Database"
description: "This preset passes the Kubernetes restricted Pod Security Standard while running persistent storage: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all..."
type: "preset"
rank: "03"
presetSlug: "03-hardened-database"
componentSlug: "statefulset"
componentTitle: "StatefulSet"
provider: "kubernetes"
icon: "package"
order: 3
---

# Hardened Database

This preset passes the Kubernetes restricted Pod Security Standard while running persistent storage: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all Linux capabilities dropped, the runtime-default seccomp filter, and no API token mount. Identity is composed — pods run as a `KubernetesServiceAccount` resource referenced by output, so permissions and cloud federation live on the identity, not the workload.

The storage-specific piece of the hardening is `fsGroup`: persistent volumes are provisioned root-owned, and a non-root database cannot write to them without it.

## When to Use

- Production databases in clusters that enforce Pod Security Standards
- Multi-tenant or security-sensitive environments running stateful workloads
- As the template for any stateful service that handles customer data

## Key Configuration Choices

- **`fsGroup: 10001` + `fsGroupChangePolicy: OnRootMismatch`** — the volume is group-owned by the database's GID so the non-root process can write it; `OnRootMismatch` skips the recursive chown when ownership already matches, which matters on large data volumes where a full re-chown can add minutes to pod start
- **`runAsNonRoot` + pinned UID/GID 10001** — refuses to start images that silently default to root; pinned IDs make on-disk file ownership deterministic across image upgrades
- **`readOnlyRootFilesystem` + EmptyDir /tmp** — the container cannot modify its own image; the data volume and the size-limited EmptyDir cover every legitimate write
- **`drop: ["ALL"]` capabilities + `RuntimeDefault` seccomp** — the restricted-profile baseline; databases rarely need any capability added back
- **`automountServiceAccountToken: false`** — a database never calls the Kubernetes API and should not carry credentials for it
- **ServiceAccount by reference** — the identity (with any workload-identity binding for cloud-native backups) deploys with the chart and flows in as a reference

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your `KubernetesNamespace` resource |
| `<your-container-registry>/<your-database-image>` | Database container image | Your container registry |
| `<your-image-tag>` | Image tag or version | Your registry or CI/CD pipeline output |
| `<your-service-account-resource>` | The `KubernetesServiceAccount` resource pods run as | Your identity resources |

## Related Presets

- **01-database** — The same shape without the hardening
- **02-ha-quorum-cluster** — Multi-member availability posture; combine with this preset's security block for hardened clusters

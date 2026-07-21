---
title: "Hardened Production Service"
description: "This preset passes the Kubernetes restricted Pod Security Standard: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all Linux capabilities dropped, the..."
type: "preset"
rank: "04"
presetSlug: "04-hardened-production"
componentSlug: "deployment"
componentTitle: "Deployment"
provider: "kubernetes"
icon: "package"
order: 4
---

# Hardened Production Service

This preset passes the Kubernetes restricted Pod Security Standard: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all Linux capabilities dropped, the runtime-default seccomp filter, no API token mount, and replicas spread across availability zones. Identity is composed — pods run as a `KubernetesServiceAccount` resource referenced by output, so permissions and cloud federation live on the identity, not the workload.

## When to Use

- Production services in clusters that enforce Pod Security Standards
- Multi-tenant or security-sensitive environments
- As the template for any service that handles customer data

## Key Configuration Choices

- **`runAsNonRoot` + pinned UID/GID 10001** — refuses to start images that silently default to root; the pinned IDs make file ownership deterministic
- **`readOnlyRootFilesystem` + EmptyDir /tmp** — the container cannot modify its own image; the size-limited EmptyDir covers legitimate scratch writes
- **`drop: ["ALL"]` capabilities + `RuntimeDefault` seccomp** — the restricted-profile baseline; add back only what the app demonstrably needs (e.g. `NET_BIND_SERVICE` to bind below 1024)
- **`automountServiceAccountToken: false`** — app pods that never call the Kubernetes API should not carry credentials for it
- **Zone topology spread + 3 replicas + PDB `minAvailable: "2"`** — a zone outage or a drain can never take the service below two serving pods
- **ServiceAccount by reference** — the identity (with any workload-identity binding) deploys with the chart and flows in as a reference

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your `KubernetesNamespace` resource |
| `<your-container-registry>/<your-image>` | Container image repository | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-service-account-resource>` | The `KubernetesServiceAccount` resource pods run as | Your identity resources |

## Related Presets

- **02-web-service-with-hpa** — Same production rollout posture without the security hardening

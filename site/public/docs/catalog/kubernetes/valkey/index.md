---
title: "Valkey"
description: "Valkey deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesvalkey"
---

# Kubernetes Valkey

Deploys Valkey — the Linux Foundation's Redis-compatible in-memory data
store — from the official Valkey Helm chart, in one of two topologies:
standalone (the default, for caches and development) or primary/replica
replication with a dedicated read Service. Durability is declared
through typed valkey.conf fields (append-only file, RDB snapshots,
memory ceiling and eviction policy), authentication through ACL users
materialized as a Kubernetes Secret, and TLS through an existing
certificate Secret.

> **Redis compatibility**: Valkey is the BSD-licensed fork of Redis 7.2
> governed by the Linux Foundation — a drop-in replacement. Every Redis
> client library, the `RESP` protocol, and the command surface work
> unchanged; point your existing Redis clients at the exported endpoint.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Helm release** (official `valkey` chart, pinned 0.11.0) with the
  chart fullname pinned to `metadata.name`: a Deployment (standalone)
  or a StatefulSet (replication), the write Service `<name>`, and — in
  replication mode — the read Service `<name>-read` and the headless
  Service `<name>-headless`
- **Auth Secret** (`<name>-auth`, when ACL users are declared) — one
  key per username; the chart consumes it, passwords never appear in
  rendered values
- **PersistentVolumeClaims** — one for standalone persistence, or one
  per pod in replication mode (where persistence is required)
- **Metrics sidecar and Service** (when metrics are enabled) — the
  redis_exporter the chart ships, plus an optional ServiceMonitor
- **PodDisruptionBudget** — replication mode only, when enabled

## Prerequisites

- A Kubernetes namespace that already exists, or set `create_namespace`
- A StorageClass for persistence (most managed clusters provide a
  default; or reference a KubernetesStorageClass)
- For TLS: an existing kubernetes.io/tls Secret — issue a
  KubernetesCertificate and reference it via `tls.certificate_secret`
- For `metrics.service_monitor_enabled`: the Prometheus operator CRDs
  on the cluster (the release fails to install without them)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesValkey
metadata:
  name: sessions
spec:
  namespace:
    value: sessions
  create_namespace: true
  replication:
    replicas: 2
    persistence:
      size: 10Gi
    min_replicas_to_write: 1
  config:
    append_only: true
    max_memory: 1gb
  auth:
    users:
      - name: default
        password: <set-a-strong-password>
```

The chart brings up one primary and two replicas; applications write
through the `<name>` Service (the exported `kube_endpoint`) and can
read through `<name>-read`, authenticating with the `default` user's
password from the `<name>-auth` Secret.

## Configuration

### Topologies

Omit `replication` for a standalone instance (Deployment-backed;
persistence optional — without it the cache starts empty after every
restart). Declare `replication` for one primary plus N replicas
(StatefulSet-backed; persistence required — replicas full-sync from the
primary's dataset). Replication is NOT automated failover: the chart
does not ship Sentinel or Cluster mode yet, so a dead primary is
restarted by Kubernetes rather than replaced by a promoted replica —
durability across that restart comes from the append-only file on a
persistent volume. `min_replicas_to_write` adds write safety: a primary
that cannot see enough in-sync replicas refuses writes.

### Durability and memory

The typed `config` block renders valkey.conf: `append_only` for
lossless restarts, `rdb_save_points` or `snapshots_disabled` for
snapshot control, `max_memory` and `max_memory_policy` for the memory
ceiling and eviction behavior (`allkeys-lru` for caches; the
`noeviction` default for durable stores), and `extra_directives` for
anything beyond the typed fields. Size `resources` memory above
`max_memory` — Valkey needs headroom for replication buffers and
fork-based persistence.

### Authentication

Auth is OFF by default (the chart default — anyone with network reach
has full access; pair with a NetworkPolicy or keep it to development).
Declaring `auth.users` turns on ACL authentication; the list must
include the `default` user, and passwords land in the `<name>-auth`
Secret referenced by the exported outputs.

### TLS

`tls.enabled` with `certificate_secret` serves TLS from an existing
kubernetes.io/tls Secret — reference a KubernetesCertificate for the
cert-manager path; `require_client_certificate` adds mutual TLS.

### Escape hatch

`helm_values` merges additional chart values LAST (Helm `-f`
semantics) for the chart surface beyond the typed fields — never the
primary interface, and never for secrets.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the instance runs in |
| `service` | Write Service (`<name>`) — the primary in replication mode |
| `read_service` | Read Service (`<name>-read`) — empty outside replication mode |
| `headless_service` | Pod-discovery Service (`<name>-headless`) — replication mode only |
| `kube_endpoint` | In-cluster connection host (`<name>.<namespace>.svc.cluster.local:<port>`) |
| `port_forward_command` | Workstation access when no exposure is composed |
| `username` | ACL username (`default` when auth is declared; empty when auth is off) |
| `password_secret` | `{name, key}` of the password in the `<name>-auth` Secret — unset when auth is off |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the target namespace via reference
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — issues the kubernetes.io/tls Secret `tls.certificate_secret` wires
  in
- [KubernetesNetworkPolicy](/docs/catalog/kubernetes/network-policy)
  — restricts who can reach the store; essential when auth is off

## Next Steps

Compose exposure when the store must be reachable from outside — a
KubernetesService of type LoadBalancer or a Gateway TCP route against
the `service` output (this component never embeds one). Declare ACL
users before anything production-shaped touches the store, set
`max_memory` with an explicit eviction policy, and turn on
`append_only` with persistence wherever a restart must not lose data.

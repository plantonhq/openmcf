---
title: "Qdrant"
description: "Qdrant deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesqdrant"
---

# Kubernetes Qdrant

Deploys a Qdrant vector database — the search engine behind RAG,
semantic search, agent-memory and recommendation architectures — from
the official Qdrant Helm chart. One resource is one cluster:
distributed mode is always on, pod 0 bootstraps the Raft consensus,
and scaling is the `replicas` field. API keys are off by default
(the upstream posture); declared keys materialize in a chart-owned
Kubernetes Secret (never plaintext in rendered values). Storage is a
persistent volume per pod with an optional separate snapshots volume,
TLS on the client listeners mounts a cert-manager-shaped Secret, and
the Service stays ClusterIP — external reachability composes from
ingress and gateway kinds.

> **Why a vector database**: RAG retrieval fetches the passages whose
> embeddings sit closest to the question's; agent memory recalls
> experiences by similarity. Qdrant serves those nearest-neighbor
> queries over HNSW indexes with payload filtering during the search
> — the data layer those architectures stand on.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Helm release** (official `qdrant` chart, pinned 1.18.2, named
  `metadata.name`): the StatefulSet (`replicas` pods, distributed
  mode on), the ClusterIP Service (REST 6333, gRPC 6334, p2p 6335)
  plus its headless twin, one storage PersistentVolumeClaim per pod
  (default 10Gi; a separate snapshots PVC when declared), and the
  engine's `production.yaml` ConfigMap
- **API-key Secret** (`<name>-apikey`, when a key is declared) —
  chart-owned, keys `api-key` / `read-only-api-key`; the generate arm
  creates a random key once (stable across upgrades), the
  existing-secret arm copies from a Secret you own — key material
  never lands in rendered Helm values

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- A StorageClass for the volumes (most managed clusters provide a
  default; or reference a KubernetesStorageClass)
- For `api_key.existing_secret` / `read_only_api_key.existing_secret`:
  the referenced Secret must exist before the install — the chart
  reads it at template time
- For `tls`: an existing TLS Secret carrying `tls.crt`/`tls.key` (the
  shape cert-manager produces; a KubernetesCertificate reference
  works directly)
- For `service_monitor_enabled`: the Prometheus Operator CRDs

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesQdrant
metadata:
  name: vectors
spec:
  namespace:
    value: vector
  create_namespace: true
  api_key:
    generate: true
  storage:
    size: 20Gi
```

SDKs connect at the exported `grpc_endpoint`
(`vectors.vector.svc.cluster.local:6334`) with the key from the
`vectors-apikey` Secret (key `api-key`); REST rides the
`http_endpoint` (port 6333).

## Configuration

### Cluster size

`replicas: 1` (the default) is a single node — still a one-member
cluster, since distributed mode is always on. Any higher count forms
a Raft cluster over the p2p port; 3 is the production posture (quorum
survives one loss). Collections replicate per their own
`replication_factor` through the Qdrant API — shards and replication
factor are data operations, not spec fields.

### Authentication

Off by default (upstream posture — private dev clusters only).
Declare `api_key` with `generate: true` or `existing_secret`; add
`read_only_api_key` for query-only consumers. A read-only key
requires a read-write key — the spec enforces it. Key material lives
in the chart-owned `<name>-apikey` Secret, exported in the outputs.

### Storage and snapshots

`storage` sizes the per-pod data volume (vectors, payloads, WAL —
default 10Gi). When using snapshots or snapshot shard transfers,
declare the separate `snapshots` volume and size it like `storage` —
upstream's recommendation so a large snapshot cannot crash the node.
Size container memory (`resources`) for the vectors you plan to hold;
Qdrant keeps hot segments and indexes in RAM.

### Exposure and TLS

The Service is ClusterIP. Compose a KubernetesIngress or Gateway API
route over `service_name` for external access (gRPC needs an
HTTP/2-aware route). The `tls` block enables HTTPS/gRPC-TLS on the
client listeners from an existing certificate Secret; inter-node p2p
TLS rides `helm_values` (`config.cluster.p2p`).

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `release_name` | Helm release name (= `metadata.name`) |
| `service_name` | Main Qdrant Service (http 6333, grpc 6334) |
| `http_endpoint` | In-cluster REST endpoint (`http(s)://...:6333`) |
| `grpc_endpoint` | In-cluster gRPC endpoint SDKs default to (`...:6334`) |
| `api_key_secret_name` | Chart-owned `<name>-apikey` Secret (key `api-key`); empty when unauthenticated |
| `read_only_api_key_secret_name` | The same Secret (key `read-only-api-key`); empty when not declared |
| `port_forward_command` | Workstation REST access when no exposure is composed |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/storage-class)
  — backs the storage and snapshots volumes via reference
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — issues the TLS Secret the `tls` block mounts
- [KubernetesIngress](/docs/catalog/kubernetes/ingress) —
  composes REST exposure over the service handle

## Next Steps

Declare an API key before anything beyond a private dev cluster
depends on the instance — the unauthenticated default is the upstream
posture, not a production one. Size memory for the corpus you expect
and the storage volume for its vectors and WAL; add the snapshots
volume the moment backups enter the picture. Scale to 3 replicas with
`pod_anti_affinity` for availability, and remember that collection
shards and replication factor are declared through the Qdrant API,
not the manifest.

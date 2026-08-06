# Kubernetes Neo4j

Deploys a Neo4j graph database — the standard engine for knowledge
graphs, GraphRAG and agent-memory architectures — from the official
Neo4j Helm chart. One resource is one server: the community edition
(the default) is single-instance by license, and Enterprise clustering
is built from multiple resources sharing a `cluster_name`. Admin
credentials materialize as a Kubernetes Secret the chart consumes
(never plaintext in rendered values), storage is a persistent volume
sized in the spec, memory tuning is typed, and the chart's default
LoadBalancer exposure is deliberately overridden to ClusterIP —
external reachability composes from ingress and gateway kinds.

> **Why a graph database**: GraphRAG retrieval walks a knowledge graph
> at query time; agent memory accumulates entities and relationships
> across sessions. That relationship-shaped data is what a native
> graph engine serves without the JOIN-depth penalty.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Auth Secret** (`<name>-auth`, when `auth.password` is declared) —
  the chart's contract: one key, `NEO4J_AUTH`, value
  `neo4j/<password>`; created BEFORE the release (the chart looks it
  up at template time) and the only place the credential lands
- **Helm release** (official `neo4j` chart, pinned 2026.6.0, named
  `metadata.name`): the single-pod StatefulSet, the data
  PersistentVolumeClaim (default 10Gi), the default ClusterIP Service
  (bolt 7687, http 7474, https 7473), and the exposure Service —
  pinned to ClusterIP instead of the chart's LoadBalancer default

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- A StorageClass for the data volume (most managed clusters provide a
  default; or reference a KubernetesStorageClass)
- For `auth.existing_secret`: a Secret already carrying
  `NEO4J_AUTH: neo4j/<password>` — it must exist before the install
- For `ssl`: Secrets carrying `private.key`/`public.crt` (the chart's
  expected keys — cert-manager Secrets carry `tls.key`/`tls.crt`;
  bridge the key names explicitly)
- For `edition: enterprise`: a valid Neo4j license and
  `accept_license_agreement: true`
- For `service_monitor_enabled`: the Prometheus Operator CRDs

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesNeo4j
metadata:
  name: knowledge-graph
spec:
  namespace:
    value: graph
  create_namespace: true
  auth:
    password: <set-a-strong-password>
  data_volume:
    size: 50Gi
```

Drivers connect at the exported `bolt_endpoint`
(`neo4j://knowledge-graph.graph.svc.cluster.local:7687`) with the
`neo4j` user's password from the `knowledge-graph-auth` Secret; the
Browser rides the `http_endpoint` (port 7474).

## Configuration

### Editions and clustering

`community` (default) runs exactly one instance — the spec rejects
`cluster_name` on it. `enterprise` enables clustering: install one
KubernetesNeo4j resource per member, all sharing `cluster_name`, each
accepting the license agreement.

### Sizing

The chart rejects installs below its floor (500m CPU / 2Gi memory);
empty `resources` uses the chart's own defaults. The typed `memory`
block renders heap and page-cache into neo4j.conf; empty lets Neo4j
auto-compute from container memory.

### Exposure

The default is in-cluster only (ClusterIP — a deliberate override of
the chart's LoadBalancer default). Compose a KubernetesIngress or
Gateway API route for HTTP/Browser; for external bolt drivers, use a
TCP route or set `service.type: LoadBalancer` with cloud annotations.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the server runs in |
| `release_name` | Helm release name (= `metadata.name`) |
| `service_name` | Main Neo4j Service (bolt/http/https ports) |
| `bolt_endpoint` | In-cluster bolt endpoint (`neo4j://...:7687`) |
| `http_endpoint` | In-cluster HTTP API / Browser endpoint (port 7474) |
| `auth_secret_name` | Secret holding the admin credentials (`<name>-auth`, or the referenced existing Secret); empty when the chart generated a random password |
| `port_forward_command` | Workstation bolt access when no exposure is composed |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) —
  provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/kubernetesstorageclass)
  — backs the data volume via reference
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — issues TLS material for the `ssl` scopes (with the
  `private.key`/`public.crt` key bridge)
- [KubernetesIngress](/docs/catalog/kubernetes/kubernetesingress) —
  composes HTTP/Browser exposure over the service handle

## Next Steps

Declare a real password (or an existing Secret) before anything
depends on the instance — the random-password path logs it once and
exports no Secret handle. Size the data volume and memory for the
graph you expect, and compose exposure only where drivers live outside
the cluster. For Enterprise clusters, add members as separate
resources sharing `cluster_name`.

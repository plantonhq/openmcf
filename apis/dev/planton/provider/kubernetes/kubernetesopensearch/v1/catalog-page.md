# Kubernetes OpenSearch

Declares an OpenSearch cluster — the Apache-2.0 search and analytics
engine, drop-in compatible with the Elasticsearch APIs at the 7.10
fork line, with its own 2.x/3.x feature line since — reconciled by the
OpenSearch Kubernetes Operator. One resource carries node pools with
roles and per-pool storage, the TLS and security-plugin posture, an
optional OpenSearch Dashboards console, Prometheus monitoring,
keystore entries, and snapshot repositories. Workloads connect at the
exported https endpoint; exposure composes from ingress and gateway
kinds over the exported service handles.

> **Migrating from Elasticsearch**: existing Elasticsearch clients and
> tooling speak to OpenSearch unchanged at the API-compatibility line —
> the fork is Apache-2.0 end to end, which is the reason it exists.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **OpenSearchCluster** (`opensearch.opster.io/v1`, named
  `metadata.name`) — the ONLY custom resource: topology, security,
  Dashboards, monitoring, keystore, snapshot repositories

The operator reconciles it into node StatefulSets per pool
(`<cluster>-<pool>`), the main Service (`<name>`), a discovery
Service, TLS Secrets, the `<name>-admin-password` credentials Secret,
and — when enabled — the `<name>-dashboards` Deployment and Service.

## Prerequisites

- The OpenSearch Kubernetes Operator on the cluster
  (KubernetesOpenSearchOperator) — it must watch this cluster's
  namespace (the default posture watches all namespaces)
- A StorageClass for PVC-backed pools (most managed clusters provide a
  default; or reference a KubernetesStorageClass)
- For provided certificates: a kubernetes.io/tls Secret — issue a
  KubernetesCertificate and reference it in the TLS blocks
- For `monitoring.enabled`: the Prometheus Operator CRDs (enabling it
  without them fails reconciliation)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenSearch
metadata:
  name: main
spec:
  namespace:
    value: search
  create_namespace: true
  version: "2.19.4"
  node_pools:
    - name: core
      replicas: 3
      roles: [cluster_manager, data]
      jvm: "-Xmx1G -Xms1G"
  security:
    transport_tls: {}
    http_tls: {}
  dashboards:
    enabled: true
```

The operator forms the cluster and bootstraps security; workloads
connect at the exported `http_endpoint`
(`https://main.search.svc.cluster.local:9200`) with the credentials
from the `main-admin-password` Secret.

## Security: Read This Before Production

The HTTP API serves **https in every posture** — even without a
`security` block, the image's demo security configuration is active.
And WITHOUT a custom `security.config`, the bootstrapped credentials
in `<name>-admin-password` are the image's **well-known demo admin
credentials** — the operator does not generate a random password at
this release. Inside a private cluster that is fine for development;
for anything real, bring a custom `security.config` (your own
internal_users.yml and admin credentials) or rotate the admin password
through the security API immediately after install.

## Configuration

### Topology

`node_pools` is the cluster's shape: roles (underscore forms —
`cluster_manager`, `data`, `ingest`, ...), replicas, JVM, resources,
disk size, and storage backing (PVC by default — data survives pod
loss; emptyDir for throwaway experiments). One 3-replica all-roles
pool is the smallest real cluster; one replica works for development.

### Backups

`snapshot_repositories` registers repositories (`s3`, `gcs`, `azure`,
`fs`, ...); the matching repository plugin rides `plugins_list`, and
credentials go in the `keystore` (Secrets loaded into every node, with
key renaming) — or nowhere, when the nodes carry instance/workload
identity (IRSA, GKE Workload Identity): the keyless path.

### Dashboards

A section of the same resource: `dashboards.enabled` with optional
TLS, replicas, a base path for path-rewriting proxies, and a
LoadBalancer service type for quick exposure.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | OpenSearchCluster resource name (= `metadata.name`) |
| `service_name` | Main Service fronting all nodes (`<name>`) |
| `http_endpoint` | In-cluster HTTP API endpoint — always https |
| `admin_credentials_secret_name` | `<name>-admin-password` (fields `username`/`password`); empty with a custom security config |
| `dashboards_service_name` | `<name>-dashboards`; empty when not enabled |
| `dashboards_endpoint` | In-cluster Dashboards endpoint (port 5601); empty when not enabled |
| `port_forward_command` | Workstation access when no exposure is composed |

## Related Components

- [KubernetesOpenSearchOperator](/docs/catalog/kubernetes/kubernetesopensearchoperator)
  — the engine; must be installed and watching this namespace
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) —
  provides the target namespace via reference
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — issues the TLS Secrets the provided-certificate arms reference
- [KubernetesIngress](/docs/catalog/kubernetes/kubernetesingress) —
  composes external exposure over the exported service handles

## Next Steps

Replace the demo bootstrap before anything real touches the cluster —
a custom `security.config` or an immediate admin-password rotation.
Register a snapshot repository with keystore (or keyless) credentials
and take a first snapshot. Compose exposure from KubernetesIngress or
Gateway API kinds over `service_name` and `dashboards_service_name` —
this component never embeds it.

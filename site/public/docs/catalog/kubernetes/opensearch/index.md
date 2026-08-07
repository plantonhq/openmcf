---
title: "OpenSearch"
description: "OpenSearch deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesopensearch"
---

# OpenSearch

Deploy an [OpenSearch](https://opensearch.org) cluster — the Apache-2.0 search and analytics engine (a drop-in replacement for the Elasticsearch APIs at the 7.10 fork line, with its own 2.x/3.x feature line since). The cluster is declared as an `OpenSearchCluster` custom resource reconciled by the OpenSearch Kubernetes Operator, which manages the full lifecycle: node StatefulSets per pool, cluster bootstrap, TLS, the security plugin's admin bootstrap, safe rolling upgrades, and an optional OpenSearch Dashboards console.

The topology is yours to declare: node pools carry roles (`cluster_manager`, `data`, `ingest`, …), counts, storage, and sizing. The smallest working dev shape is one all-roles pool with 2 replicas — a single manager-eligible replica cannot survive the operator's bootstrap handoff (verified live).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **OpenSearchCluster custom resource** — one StatefulSet per node pool, the cluster Services, and (when enabled) the Dashboards Deployment, all reconciled by the operator
- **Generated TLS** — a CA plus per-layer certificates for node-to-node and client traffic (the default posture), and the security plugin's admin bootstrap
- **Snapshot repositories and keystore entries** — registered on the cluster at startup when declared

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.
- **OpenSearch Operator** — a **Kubernetes OpenSearch Operator** resource must be running and watching the target namespace (its default watch is cluster-wide). Deploy it first.

### Cluster Side

- **Kernel tuning** — OpenSearch needs `vm.max_map_count` raised; the default-on privileged init container handles it. On clusters that forbid privileged init containers, tune nodes out of band and disable the dial.
- **Prometheus Operator CRDs** — only if you enable monitoring; a missing ServiceMonitor CRD fails reconciliation.

## Deploy

### Console

Open the deployment store, find **OpenSearch**, and click **Deploy**. The creation wizard walks you through namespace placement, the engine version, the node-pool topology (with the manager-quorum floor held live), engine settings, node runtime, the TLS/security posture, the optional Dashboards console, monitoring, backups, and air-gap sourcing. Start from the **Dev Minimal** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenSearch
metadata:
  name: dev-opensearch
  org: acme-corp
  env: dev
spec:
  namespace:
    value: dev-opensearch
  create_namespace: true
  version: 2.19.4
  node_pools:
    - name: all
      replicas: 2
      roles:
        - cluster_manager
        - data
        - ingest
      jvm: -Xmx1G -Xms1G
      disk_size: 10Gi
```

```shell
planton apply -f opensearch-cluster.yaml
```

### InfraChart

Compose the cluster behind its operator and give Dashboards real exposure:

```yaml
spec:
  namespace:
    value: search
  version: 2.19.4
  dashboards:
    enabled: true
```

Reference the exported `dashboards_service_name` from a **Kubernetes Ingress** or Gateway API route for team access.

## Key Configuration

**The manager floor** — a lone 1-replica `cluster_manager` pool is stranded by the operator's temporary bootstrap manager (every write returns `cluster_manager_not_discovered`). Declare at least 2 manager-eligible replicas; production runs 3 dedicated managers.

**The admin password is the image's demo credential** — the operator does not generate a random admin password at this release. The bootstrapped `<name>-admin-password` Secret holds the well-known demo pair: rotate it through the security API immediately after install, or bring a custom security config, before anything real uses the cluster.

**Storage is per-pool and durable by default** — a PVC per node on the default StorageClass. emptyDir is for throwaway experiments only; heap should be about half the container memory (`jvm` sets it; the rest is OS page cache).

**Version bumps are day-2 friendly** — the operator rolls nodes one at a time with drain ordering. Check its compatibility table before a major-line jump, and keep an enabled Dashboards aligned with the engine version.

**Snapshot credentials live in the keystore** — repository settings pass verbatim into the CR (never put credentials there), and a cloud repository type needs its plugin (`repository-s3`, …) in the node plugins list.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the cluster runs |
| `spec.node_pools[].persistence.pvc.storage_class` | KubernetesStorageClass (`status.outputs.storage_class_name`) | Per-pool volume class |
| `spec.security.transport_tls.secret` / `spec.security.http_tls.secret` | KubernetesCertificate (`status.outputs.secret_name`) | Bring-your-own TLS (the cert-manager seam) |
| `spec.security.transport_tls.ca_secret` | KubernetesSecret (`metadata.name`) | Existing CA for generated node certs |
| `spec.security.config.*` | KubernetesSecret (`metadata.name`) | Custom security-plugin config trio |
| `spec.dashboards.tls.secret` | KubernetesCertificate (`status.outputs.secret_name`) | Dashboards HTTPS certificate |
| `spec.dashboards.opensearch_credentials_secret` | KubernetesSecret (`metadata.name`) | Dashboards cluster credentials (custom security config) |
| `spec.monitoring.monitoring_user_secret` | KubernetesSecret (`metadata.name`) | Scrape credentials (custom security config) |
| `spec.keystore[].secret` | KubernetesSecret (`metadata.name`) | Secure settings loaded into the node keystore |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `http_endpoint` | In-cluster https API endpoint | Application clients, log shippers, dashboards |
| `admin_credentials_secret_name` | The bootstrapped admin Secret (username/password) — empty when a custom security config replaces the bootstrap | Client authentication |
| `service_name` / `cluster_name` / `namespace` | Cluster identity handles | Ingress/Gateway composition, monitoring |
| `dashboards_service_name` / `dashboards_endpoint` | The console's handles — empty when Dashboards is disabled | Ingress/Gateway exposure for the team |
| `port_forward_command` | Copy-paste developer access | Local exploration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Dev Minimal** — one all-roles pool with 2 replicas (the smallest shape that survives the bootstrap handoff), generated TLS on both layers. Start from the **Dev Minimal** preset.

**Production Cluster** — dedicated manager and data pools, sized heaps, durable storage, disruption budgets. Start from the **Production Cluster** preset.

**S3 Snapshots** — the production shape plus the `repository-s3` plugin, keystore-loaded credentials, and a registered S3 snapshot repository. Start from the **S3 Snapshots** preset.

## Works With

- **Kubernetes OpenSearch Operator** — the engine that reconciles this cluster; deploy it first.
- **Kubernetes Certificate** — cert-manager-issued TLS for bring-your-own postures.
- **Kubernetes Secret** — keystore sources, custom security config, credentials.
- **Kubernetes Ingress / Gateway API kinds** — real exposure for the API and Dashboards.
- **Kubernetes Kube Prometheus Stack** — the ServiceMonitor consumer when monitoring is enabled.

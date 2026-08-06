---
title: "Cloud Native PG Operator"
description: "Cloud Native PG Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetescloudnativepgoperator"
---

# Kubernetes Cloud Native PG Operator

Installs CloudNativePG — the CNCF PostgreSQL operator — from the
official Helm chart, with a typed spec over the chart's meaningful
configuration surface. The operator reconciles `Cluster` custom
resources into highly available PostgreSQL clusters: streaming
replication, automated failover, rolling updates, declarative roles and
storage, and plugin-based backups. Enable the Barman Cloud plugin arm to
add the object-store backup path every KubernetesPostgres backup block
depends on. One installation per cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and
  owned when `create_namespace` is set (`cnpg-system` is the upstream
  convention)
- **Helm Release** (`cnpg`) — the operator Deployment, its
  mutating/validating webhooks, and the CloudNativePG CRDs (Cluster,
  ScheduledBackup, Backup, Pooler, Database, ... — every CRD carries
  `helm.sh/resource-policy: keep`, so uninstalling never cascade-deletes
  the databases)
- **Helm Release** (`plugin-barman-cloud`, when
  `barman_cloud_plugin.enabled`) — the Barman Cloud CNPG-I plugin: a
  separate release in the same namespace, installed after the operator,
  carrying the object-store backup path (WAL archiving, base backups,
  restores) for every database the operator manages

## Prerequisites

- With `barman_cloud_plugin.enabled`: cert-manager on the cluster
  (KubernetesCertManager) — the plugin's internal TLS certificates are
  cert-manager-issued and the release fails to install without it
- With `monitoring.pod_monitor_enabled`: the Prometheus operator CRDs —
  the release fails to install without them
- The cluster must not already run CloudNativePG — the CRDs and webhooks
  are cluster-scoped singletons

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesCloudNativePgOperator
metadata:
  name: cnpg
spec:
  namespace:
    value: cnpg-system
  create_namespace: true
  barman_cloud_plugin:
    enabled: true # requires cert-manager; the backup path for every database
```

The operator becomes Available and starts reconciling; KubernetesPostgres
resources declared afterwards turn into running PostgreSQL clusters.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator (and the plugin, when enabled) runs in |
| `release_name` | Helm release name of the operator (always `cnpg`) |
| `barman_plugin_release_name` | Helm release name of the Barman Cloud plugin when enabled; empty otherwise |

## Next Steps

Declare databases as KubernetesPostgres resources — one per PostgreSQL
cluster; the operator reconciles them into instances, services, and
credential secrets. Enable the plugin arm before any database declares a
backup block. Pin `chart_version` deliberately (default 0.29.0 =
operator 1.30.0) and upgrade the operator on the platform's schedule —
removing the component never deletes the CRDs or the databases behind
them.

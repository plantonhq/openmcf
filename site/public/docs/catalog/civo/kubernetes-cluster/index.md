---
title: "Kubernetes Cluster"
description: "Kubernetes Cluster deployment documentation"
icon: "package"
order: 100
componentName: "civokubernetescluster"
---

# Kubernetes Cluster on Civo

Deploys a managed Kubernetes cluster (K3s) on Civo Cloud with configurable node sizing, optional high availability, automatic version upgrades, and VPC networking. Civo clusters provision in under 90 seconds, making them well-suited for development, CI/CD, and lightweight production workloads. Integrates with Planton's Provider Connections for Civo credential management and ValueFromRef for VPC network dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Kubernetes Cluster** -- a managed K3s-based cluster in the specified Civo region, placed on the referenced VPC network with the configured Kubernetes version and node pool
- **Default Node Pool** -- worker nodes with the configured instance size and count from `defaultNodePool`
- **Kubeconfig** -- generated automatically and stored in stack outputs for downstream access

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo VPC network** in the target region. Provide the network ID directly or reference a CivoVpc Cloud Resource via ValueFromRef.
- **A supported Kubernetes version** -- check available versions via the Civo CLI (`civo kubernetes versions`) or Civo dashboard. Specify the full version string (e.g., `"1.29.2"`).

## Deploy

### Console

Open the deployment store, find **Kubernetes Cluster on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production HA** preset in the [Presets](#presets) tab for a resilient multi-node configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoKubernetesCluster
metadata:
  name: app-cluster
  org: acme-corp
  env: prod
spec:
  clusterName: app-cluster
  region: lon1
  kubernetesVersion: "1.29.2"
  network:
    value: "abc12345-6789-def0-1234-567890abcdef"
  defaultNodePool:
    size: g4s.kube.medium
    nodeCount: 3
```

```shell
planton apply -f civo-cluster.yaml
```

This creates a 3-node Kubernetes cluster on Civo's London region with medium-sized nodes. No HA or auto-upgrade is enabled. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to a VPC network deployed in the same InfraPipeline:

```yaml
spec:
  network:
    valueFrom:
      kind: CivoVpc
      name: app-network
      fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC first, then provisions the Kubernetes cluster on it.

## Key Configuration

These are the most important decisions when configuring a Civo Kubernetes cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Node pool sizing** -- The `defaultNodePool.size` field sets the instance type for worker nodes (e.g., `"g4s.kube.small"` for development, `"g4s.kube.medium"` for production). `defaultNodePool.nodeCount` sets how many nodes to provision -- use at least 3 for production workloads to allow pod scheduling across nodes during rolling updates and node failures.

**High availability** -- Set `highlyAvailable: true` to provision multiple control plane nodes for fault tolerance. Without HA, a single control plane node failure takes the cluster offline. HA increases the cluster's base cost but is required for production SLAs.

**Automatic upgrades** -- Set `autoUpgrade: true` to let Civo apply Kubernetes patch updates automatically. Recommended for production to reduce operational overhead. Disable for environments where you need manual control over the exact Kubernetes version.

**Region** -- The `region` field accepts a Civo region code (e.g., `lon1` for London, `nyc1` for New York, `fra1` for Frankfurt). Choose the region closest to your users or application workloads.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CivoVpc** | `network` | `status.outputs.network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | Unique cluster identifier on Civo | Civo API operations, monitoring dashboards |
| `kubeconfig_b64` | Base64-encoded kubeconfig for cluster access | Kubernetes Provider Connections, CI/CD pipeline configuration |
| `api_server_endpoint` | Kubernetes API server endpoint URL | kubectl configuration, application health checks |
| `created_at_rfc3339` | Cluster creation timestamp in RFC 3339 format | Audit logs, lifecycle tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production HA cluster** -- 3 medium-sized nodes with high availability enabled and automatic Kubernetes patch upgrades. Provides resilience against node failures and zero-downtime upgrades for production workloads. Start from the **Production HA** preset.

**Development cluster** -- Single small node, no HA, no auto-upgrade. Minimal cost for development, CI/CD test clusters, and proof-of-concept deployments. Civo's fast provisioning (under 90 seconds) makes create/destroy cycles practical. Start from the **Development** preset.

## Works With

- [**Civo VPC**](/cloud-catalog/civo-vpc) -- provides the VPC network for cluster networking

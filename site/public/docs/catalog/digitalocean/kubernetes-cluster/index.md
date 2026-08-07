---
title: "Kubernetes Cluster"
description: "Kubernetes Cluster deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceankubernetescluster"
---

# Kubernetes Cluster on DigitalOcean

Deploys a managed Kubernetes cluster (DOKS) on DigitalOcean with configurable node sizing, optional high availability, automatic patch upgrades, control plane firewall restrictions, container registry integration, and VPC networking. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for VPC dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Kubernetes Cluster** -- a managed DOKS cluster in the specified region and VPC, using the configured Kubernetes version, with optional HA control plane, auto-upgrade, and surge upgrade settings
- **Default Node Pool** -- worker nodes with the configured instance size, node count, and optional autoscaling between `minNodes` and `maxNodes`
- **Control Plane Firewall** -- created only when `controlPlaneFirewallAllowedIps` are provided; restricts API server access to the specified CIDR blocks
- **Maintenance Policy** -- created only when `maintenanceWindow` is provided; schedules cluster updates during the specified window
- **Container Registry Integration** -- configured only when `registryIntegration` is enabled; provisions imagePullSecrets for DigitalOcean Container Registry
- **DigitalOcean Tags** -- user-supplied tags applied to the cluster for resource organization

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A VPC network** in the target region. Provide the VPC name directly or reference a DigitalOceanVpc Cloud Resource via ValueFromRef.
- **A supported Kubernetes version** -- check available versions via the DigitalOcean CLI (`doctl kubernetes options versions`) or the DOKS documentation. Specify the version string (e.g., `"1.31"`).

## Deploy

### Console

Open the deployment store, find **Kubernetes Cluster on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production HA** preset in the [Presets](#presets) tab for a resilient multi-node configuration with autoscaling.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanKubernetesCluster
metadata:
  name: app-cluster
  org: acme-corp
  env: prod
spec:
  clusterName: app-cluster
  region: nyc1
  kubernetesVersion: "1.31"
  vpc:
    value: "app-network"
  defaultNodePool:
    size: s-4vcpu-8gb
    nodeCount: 3
```

```shell
planton apply -f do-k8s-cluster.yaml
```

This creates a 3-node Kubernetes cluster in the NYC1 region. No HA, auto-upgrade, autoscaling, or API server firewall is configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to a VPC deployed in the same InfraPipeline:

```yaml
spec:
  vpc:
    valueFrom:
      kind: DigitalOceanVpc
      name: app-network
      fieldPath: metadata.name
```

The InfraPipeline resolves the dependency graph, deploys the VPC first, then provisions the Kubernetes cluster on it.

## Key Configuration

These are the most important decisions when configuring a DOKS cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**High availability** -- Set `highlyAvailable: true` to provision multiple control plane replicas for fault tolerance. Without HA, a single control plane failure takes the cluster offline. HA increases cost but is required for production SLAs.

**Node pool sizing and autoscaling** -- The `defaultNodePool.size` field sets the instance type (e.g., `"s-2vcpu-4gb"` for development, `"s-4vcpu-8gb"` for production). Enable `autoScale` with `minNodes` and `maxNodes` to let DOKS adjust node count based on pod scheduling pressure.

**Automatic upgrades** -- Set `autoUpgrade: true` to let DigitalOcean apply Kubernetes patch updates during the maintenance window. Combine with `maintenanceWindow` (e.g., `"sunday=03:00"`) to control when upgrades occur.

**API server firewall** -- Populate `controlPlaneFirewallAllowedIps` with CIDR blocks to restrict who can reach the Kubernetes API server. When empty, the API server is publicly accessible. Restrict to VPN or office CIDRs for production.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanVpc** | `vpc` | `metadata.name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | Unique cluster identifier (UUID) on DigitalOcean | DigitalOcean API operations, monitoring dashboards |
| `kubeconfig` | Base64-encoded kubeconfig for cluster access | Kubernetes Provider Connections, CI/CD pipeline configuration |
| `api_server_endpoint` | Kubernetes API server endpoint URL | kubectl configuration, application health checks |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production HA cluster** -- HA control plane, autoscaling node pool (2-5 nodes), automatic patch upgrades with a Sunday maintenance window, container registry integration, and API server firewall restrictions. Start from the **Production HA** preset.

**Development cluster** -- Non-HA control plane, 2 fixed-size nodes with smaller instances, no auto-upgrade or API firewall. Minimal cost for dev/test and feature branch clusters. Start from the **Development** preset.

## Works With

- [**DigitalOcean VPC**](/cloud-catalog/digital-ocean-vpc) -- provides the VPC network for cluster placement
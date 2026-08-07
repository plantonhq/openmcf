# Scaleway Kapsule Cluster

Deploys a managed Kubernetes cluster on Scaleway Kapsule as a composite resource that bundles the control plane and a default node pool into a single declarative unit. Configurable CNI plugin (Cilium or Calico), cluster type (shared or dedicated control plane), automatic patch upgrades, cluster-wide autoscaler settings, and node pool sizing provide production-ready Kubernetes from a single manifest. Additional node pools can be added via separate ScalewayKapsulePool resources. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kapsule Cluster** -- a managed Kubernetes control plane (API server, etcd, scheduler, controller-manager) in the specified region with the configured CNI, Kubernetes version, and optional autoscaler and auto-upgrade settings
- **Default Node Pool** -- worker nodes with the configured instance type, count, and optional autoscaling, autohealing, and upgrade policy
- **Kubeconfig** -- generated automatically and stored in stack outputs for downstream cluster access
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway Private Network** in the target region. All Kapsule clusters require a Private Network for node-to-control-plane and node-to-node communication. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef.
- **A supported Kubernetes version** -- check available versions via the Scaleway console or API. Specify a minor version (e.g., `"1.32"`) when auto-upgrade is enabled so Scaleway can apply patch updates.

## Deploy

### Console

Open the deployment store, find **Scaleway Kapsule Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Autoscaling** preset in the [Presets](#presets) tab for a cluster with autoscaling, auto-upgrade, and private nodes.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayKapsuleCluster
metadata:
  name: app-cluster
  org: acme-corp
  env: prod
spec:
  region: fr-par
  kubernetesVersion: "1.32"
  cni: cilium
  privateNetworkId:
    value: "abc12345-6789-def0-1234-567890abcdef"
  defaultNodePool:
    nodeType: DEV1-M
    size: 2
```

```shell
planton apply -f scaleway-kapsule-cluster.yaml
```

This creates a Kapsule cluster with Cilium CNI and a 2-node DEV1-M default pool. No autoscaling, auto-upgrade, or private node configuration is applied. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to a Private Network deployed in the same InfraPipeline:

```yaml
spec:
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC and Private Network first, then provisions the Kapsule cluster with the resolved Private Network ID.

## Key Configuration

These are the most important decisions when configuring a Kapsule cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CNI plugin** -- The `cni` field selects the pod networking implementation. `cilium` (recommended) uses eBPF for high performance, advanced network policies, and Hubble observability. `calico` is a mature alternative for teams already familiar with it. Cannot be changed after creation.

**Cluster type** -- The `type` field defaults to `kapsule` (shared control plane, no additional cost). For production workloads requiring isolated API server SLAs, use `kapsule-dedicated-4`, `kapsule-dedicated-8`, or `kapsule-dedicated-16`. Can be upgraded after creation.

**Default node pool** -- The `defaultNodePool` ships with the cluster and provides immediate compute capacity. Set `nodeType` to `DEV1-M` for development or `PRO2-S` for production. Enable `autoScale` with `minSize` and `maxSize` bounds for elastic workloads. Enable `autohealing` and set `publicIpDisabled` to true for production security.

**Auto-upgrade** -- Configure `autoUpgrade` to let Scaleway apply Kubernetes patch updates during a specified maintenance window (e.g., Sunday 3:00 AM UTC). Only patch versions within the current minor are upgraded automatically; minor version upgrades require manual action.

**Autoscaler configuration** -- The `autoscalerConfig` controls cluster-wide autoscaler behavior (scale-down delays, utilization thresholds, expander strategy). These settings apply to all autoscaling-enabled pools. The defaults (binpacking estimator, 10m scale-down delay, 0.5 utilization threshold) work for most workloads.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | Unique identifier of the Kapsule cluster | ScalewayKapsulePool cluster reference, monitoring dashboards |
| `kubeconfig` | Raw kubeconfig file content for cluster access | Kubernetes Provider Connections, CI/CD pipeline configuration |
| `apiserver_url` | Kubernetes API server endpoint URL | kubectl configuration, external monitoring, IaC provider setup |
| `cluster_ca_certificate` | Base64-encoded CA certificate of the API server | IaC Kubernetes provider configuration for addon deployments |
| `wildcard_dns` | DNS wildcard for ready nodes in the cluster | DNS-based service discovery, CNAME targets |
| `default_pool_id` | Unique identifier of the default node pool | Monitoring, distinguishing default pool from additional pools |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development cluster** -- A shared control plane with a fixed 2-node DEV1-M pool, no auto-upgrade, and public IPs on nodes. Minimal cost for development, testing, and learning Kubernetes on Scaleway. Start from the **Dev Minimal** preset.

**Production autoscaling cluster** -- A shared control plane with PRO2-S nodes that autoscale between 2 and 10, automatic Sunday-morning patch upgrades, autohealing, and private-only nodes. The standard production configuration for elastic workloads. Start from the **Production Autoscaling** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides the network for node-to-control-plane and node-to-node communication
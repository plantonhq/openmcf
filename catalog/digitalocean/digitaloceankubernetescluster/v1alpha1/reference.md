# DigitalOceanKubernetesCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanKubernetesClusterSpec defines the specification for creating a managed Kubernetes cluster on DigitalOcean.
It focuses on essential parameters for a production-grade cluster, following the 80/20 principle to expose only the most commonly used settings.

## Example

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanKubernetesCluster
metadata:
  name: first-cluster                 # K8s resource name (unique per namespace)
spec:
  clusterName: first-cluster          # Must be unique in your DigitalOcean account
  region: blr1                       # Any valid DigitalOceanRegion slug (nyc3, sfo3, blr1, …)
  kubernetesVersion: "1.33"        # Must match a version currently offered by DigitalOcean
  vpc:
    value: b5648f9e-a28a-4760-bb87-b2fad07ae295
  highlyAvailable: false              # HA control plane (extra cost)
  autoUpgrade: false                  # Automatic patch upgrades
  disableSurgeUpgrade: false         # Keep surge nodes for zero‑downtime upgrades
  maintenanceWindow: "sunday=02:00"  # Scheduled maintenance window (optional)
  registryIntegration: false         # Enable DOCR integration (optional)
  # controlPlaneFirewallAllowedIps:  # Restrict API access to specific IPs (optional)
  #   - "203.0.113.5/32"
  #   - "198.51.100.0/24"
  tags:
    - planton
  defaultNodePool:
    size: s-2vcpu-4gb
    nodeCount: 3
    autoScale: true
    minNodes: 1
    maxNodes: 5
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.clusterName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.kubernetesVersion` | `string` | yes |  |  |
| `spec.vpc` | `string \| valueFrom` | yes |  | DigitalOceanVpc (`status.outputs.vpc_id`) |
| `spec.highlyAvailable` | `bool` |  | `false` |  |
| `spec.autoUpgrade` | `bool` |  |  |  |
| `spec.disableSurgeUpgrade` | `bool` |  |  |  |
| `spec.maintenanceWindow` | `string` |  |  |  |
| `spec.registryIntegration` | `bool` |  |  |  |
| `spec.controlPlaneFirewallAllowedIps` | `[]string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.defaultNodePool` | `DigitalOceanKubernetesClusterDefaultNodePool` | yes |  |  |
| `spec.defaultNodePool.size` | `string` | yes |  |  |
| `spec.defaultNodePool.nodeCount` | `uint32` | yes |  |  |
| `spec.defaultNodePool.autoScale` | `bool` |  |  |  |
| `spec.defaultNodePool.minNodes` | `uint32` |  |  |  |
| `spec.defaultNodePool.maxNodes` | `uint32` |  |  |  |

## Field Details

### spec.clusterName

`string` · required

The name of the Kubernetes cluster. This will be the cluster's identifier in DigitalOcean.
Constraints: Must be unique per account. (A maximum length or character set may be enforced by DigitalOcean, e.g., alphanumeric and hyphens.)

- rule: {"required":true}

### spec.region

`enum` · required

The DigitalOcean region where the cluster will be created.
Determines where the cluster's control plane and nodes are provisioned.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

### spec.kubernetesVersion

`string` · required

The Kubernetes version to use for the cluster (semantic versioning).
Must be a supported version on DigitalOcean (e.g., 1.22+).
Example: "1.26.3"

- rule: {"required":true}

### spec.vpc

`string | valueFrom` · required

Reference to the DigitalOcean VPC where the cluster's control plane will reside.
This must be an existing VPC in the same region. The cluster consumes the VPC's
ID (a DigitalOcean UUID), so a reference resolves to the DigitalOceanVpc's
exported vpc_id output rather than its metadata name.

- references: DigitalOceanVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.highlyAvailable

`bool`

Whether to enable a highly available control plane for the cluster.
If true, the cluster is created with a High Availability control plane (multiple masters for increased uptime, additional cost).
Default: false.

- default: `false`

### spec.autoUpgrade

`bool`

Whether to enable automatic patch upgrades for the cluster.
If true, the cluster will automatically upgrade to new patch releases of Kubernetes when available.

### spec.disableSurgeUpgrade

`bool`

Whether to disable surge upgrades for the cluster.
If false(default), cluster upgrades will temporarily provision extra nodes to minimize downtime during updates.

### spec.maintenanceWindow

`string`

Scheduled maintenance window for cluster updates (format: "day=HH:MM" or "any=HH:MM").
Examples: "sunday=02:00" or "any=00:00"
If not specified, DigitalOcean will apply updates at any time.

### spec.registryIntegration

`bool`

Whether to enable DigitalOcean Container Registry (DOCR) integration.
If true, automatically creates imagePullSecrets in the cluster for pulling private images from DOCR.
Default: false.

### spec.controlPlaneFirewallAllowedIps

`[]string`

List of allowed IP addresses (CIDR notation) for control plane firewall.
Restricts Kubernetes API server access to specified IPs for security.
If empty, API server is publicly accessible (not recommended for production).
Example: ["203.0.113.5/32", "198.51.100.0/24"]

### spec.tags

`[]string`

A list of tags to apply to the cluster.
Tags help organize and identify the cluster within DigitalOcean.

### spec.defaultNodePool

`DigitalOceanKubernetesClusterDefaultNodePool` · required

Reference to the default node pool for the cluster.

- rule: {"required":true}

### spec.defaultNodePool.size

`string` · required

The slug identifier for the Droplet size to use for each node (e.g., "s-4vcpu-8gb").
This defines the CPU and memory of the nodes in the pool.

- rule: {"required":true}

### spec.defaultNodePool.nodeCount

`uint32` · required

The number of nodes to provision in the pool.
Must be at least 1. If auto_scale is enabled, this acts as the initial desired node count.

- rule: {"required":true,"uint32":{"gt":0}}

### spec.defaultNodePool.autoScale

`bool`

Enable auto-scaling for this node pool.
If true, the platform will manage node count between min_nodes and max_nodes.

### spec.defaultNodePool.minNodes

`uint32`

Minimum number of nodes when auto-scaling is enabled.
Required if auto_scale is true.

### spec.defaultNodePool.maxNodes

`uint32`

Maximum number of nodes when auto-scaling is enabled.
Required if auto_scale is true.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanKubernetesCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The unique identifier (UUID) of the created Kubernetes cluster. |
| `status.outputs.kubeconfig` | `string` | A base64-encoded Kubernetes config (kubeconfig) for accessing the cluster. |
| `status.outputs.api_server_endpoint` | `string` | The endpoint URL of the Kubernetes API server for the cluster. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpc` | DigitalOceanVpc | `status.outputs.vpc_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| DigitalOceanKubernetesNodePool | `spec.cluster` | `metadata.name` |

## See Also

- [Overview](../README.md)

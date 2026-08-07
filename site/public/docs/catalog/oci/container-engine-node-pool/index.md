---
title: "Container Engine Node Pool"
description: "Container Engine Node Pool deployment documentation"
icon: "package"
order: 100
componentName: "ocicontainerenginenodepool"
---

# Container Engine Node Pool on OCI

Deploys an OKE node pool -- a set of worker nodes within an OKE cluster sharing the same compute shape, OS image, Kubernetes version, and placement rules. Multiple node pools can be attached to a single cluster for heterogeneous workloads (GPU nodes for ML, Arm nodes for cost optimization, preemptible nodes for batch jobs). Supports VCN-native pod networking with dedicated pod subnets and NSGs, multi-AD placement with fault domain constraints, rolling upgrade cycling, and node eviction controls. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for compartment, cluster, subnet, and security group wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **OKE Node Pool** -- a pool of worker nodes attached to the specified OKE cluster, with the configured compute shape, Kubernetes version, placement configs, and sizing
- **Worker Nodes** -- compute instances created and managed by OKE based on `nodeConfigDetails.size` and `placementConfigs`. Distributed across availability domains and fault domains as evenly as possible
- **Node Boot Volumes** -- created from the specified OS image with optional custom boot volume size. Encrypted with KMS when `kmsKeyId` is provided
- **Pod Networking** -- configured only when `podNetworkOptionDetails` is populated; assigns VCN IPs to pods with dedicated pod subnets and NSGs (required for VCN-native CNI clusters)
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the node pool and worker nodes

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the node pool in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- An OKE cluster to attach the node pool to. The cluster must already exist. Provide the cluster OCID directly or reference an OciContainerEngineCluster Cloud Resource via ValueFromRef.
- One or more subnets for node placement -- one per availability domain for regional distribution. Provide subnet OCIDs directly or reference OciSubnet Cloud Resources via ValueFromRef.
- For VCN-native CNI: additional subnets for pod IP allocation and optional NSGs for pod-level traffic control.
- An OS image OCID for the node OS (use `oci ce node-pool-options get` to list available images for a Kubernetes version), or omit to use the default Oracle Linux image.

## Deploy

### Console

Open the deployment store, find **Container Engine Node Pool on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Production** preset in the [Presets](#presets) tab to pre-populate a 3-node pool with VCN-native networking and rolling upgrade cycling.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciContainerEngineNodePool
metadata:
  name: general-purpose-pool
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  clusterId:
    value: "ocid1.cluster.oc1..example"
  nodeShape: VM.Standard.E4.Flex
  nodeShapeConfig:
    ocpus: 2
    memoryInGbs: 32
  nodeConfigDetails:
    size: 3
    placementConfigs:
      - availabilityDomain: "Ixxj:US-ASHBURN-AD-1"
        subnetId:
          value: "ocid1.subnet.oc1..example"
```

```shell
planton apply -f node-pool.yaml
```

This creates a 3-node pool with 2-OCPU flex instances placed in a single availability domain. No pod networking, node eviction settings, or rolling upgrade cycling are configured. The Kubernetes version inherits from the cluster.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the node pool to a compartment, cluster, subnets, and security groups deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
  clusterId:
    valueFrom:
      kind: OciContainerEngineCluster
      name: platform-cluster
      fieldPath: status.outputs.clusterId
  nodeConfigDetails:
    placementConfigs:
      - availabilityDomain: "Ixxj:US-ASHBURN-AD-1"
        subnetId:
          valueFrom:
            kind: OciSubnet
            name: worker-subnet
            fieldPath: status.outputs.subnetId
    nsgIds:
      - valueFrom:
          kind: OciSecurityGroup
          name: worker-nsg
          fieldPath: status.outputs.networkSecurityGroupId
    podNetworkOptionDetails:
      cniType: oci_vcn_ip_native
      podSubnetIds:
        - valueFrom:
            kind: OciSubnet
            name: pod-subnet
            fieldPath: status.outputs.subnetId
      podNsgIds:
        - valueFrom:
            kind: OciSecurityGroup
            name: pod-nsg
            fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, cluster, subnets, and security groups first, then provisions the node pool with the resolved values.

## Key Configuration

These are the most important decisions when configuring a node pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Shape and sizing** -- The `nodeShape` selects the compute shape for all nodes (e.g., `VM.Standard.E4.Flex` for AMD, `VM.Standard.A1.Flex` for Arm, `VM.GPU.A10.1` for GPU). Flex shapes require `nodeShapeConfig` with `ocpus` and `memoryInGbs`. The `nodeConfigDetails.size` sets the total node count, distributed across placement configs.

**Multi-AD placement** -- Each entry in `nodeConfigDetails.placementConfigs` specifies an availability domain, subnet, and optional fault domain constraints. Provide one entry per AD for HA. OKE distributes nodes across placement configs as evenly as possible. Per-placement `preemptibleNodeConfig` enables spot instances for individual ADs.

**VCN-native pod networking** -- When the parent cluster uses `oci_vcn_ip_native` CNI, the node pool must configure `podNetworkOptionDetails` with `cniType`, `maxPodsPerNode` (limited by VNIC attachment count for the shape), and dedicated `podSubnetIds` and `podNsgIds`. This gives each pod a VCN IP address and enables NSG-based pod-level traffic control.

**Rolling upgrade cycling** -- The `nodePoolCyclingDetails` block controls how nodes are replaced during Kubernetes version upgrades or shape changes. Set `isNodeCyclingEnabled: true` with `maximumSurge: "1"` and `maximumUnavailable: "0"` for zero-downtime rolling upgrades. Without cycling, nodes are replaced all at once.

**Node eviction** -- The `nodeEvictionSettings` block controls pod draining during node pool operations. `evictionGraceDuration` (ISO 8601, default PT60M) sets how long OKE tries to gracefully evict pods. `isForceDeleteAfterGraceDuration: true` proceeds with node deletion if pods cannot be evicted in time.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciContainerEngineCluster** | `clusterId` | `status.outputs.clusterId` |
| **OciSubnet** | `nodeConfigDetails.placementConfigs.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `nodeConfigDetails.nsgIds` | `status.outputs.networkSecurityGroupId` |
| **OciSubnet** (optional) | `nodeConfigDetails.podNetworkOptionDetails.podSubnetIds` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `nodeConfigDetails.podNetworkOptionDetails.podNsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `node_pool_id` | OCID of the node pool | Monitoring dashboards, OCI CLI operations, autoscaler configuration |
| `kubernetes_version` | Kubernetes version running on the nodes | Upgrade coordination, compatibility verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard production** -- A 3-node pool with 2-OCPU flex instances across three availability domains, VCN-native pod networking, worker and pod NSGs, node eviction settings, and zero-downtime rolling upgrade cycling. The starting point for most production node pools. Start from the **Standard Production** preset.

**Hardened encrypted** -- A production pool with boot volume KMS encryption, in-transit encryption for paravirtualized attachments, and restrictive NSGs. For environments requiring encryption at rest and in transit. Start from the **Hardened Encrypted** preset.

**Preemptible dev** -- A single-AD, minimal-size pool with preemptible (spot) nodes and flannel networking. Significantly lower cost for development, CI runners, and fault-tolerant batch workloads. Start from the **Preemptible Dev** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this node pool
- [**Container Engine Cluster on OCI**](/cloud-catalog/oci-container-engine-cluster) -- provides the OKE cluster that this node pool is attached to
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides subnets for node placement and VCN-native pod IP allocation
- [**Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides NSGs for node VNICs and pod VNICs
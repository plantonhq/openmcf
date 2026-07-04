---
title: "VPC"
description: "VPC deployment documentation"
icon: "package"
order: 100
componentName: "gcpvpcnetwork"
---

# GCP VPC

Deploys a GCP VPC network in custom subnet mode by default, with configurable dynamic routing and the full provider-floor depth surface: MTU, ULA internal IPv6, firewall-policy enforcement order, network profiles, BGP best-path selection, and default-route suppression. The component automatically enables the Compute Engine API on the target project before creating the network.

Private services access for Google managed services (Cloud SQL, AlloyDB, Memorystore) is not bundled here — compose a `GcpGlobalAddress` (VPC_PEERING range) with a `GcpServiceNetworkingConnection` on this network instead.

## What Gets Created

When you deploy a GcpVpcNetwork resource, Planton provisions:

- **Compute Engine API enablement** — a `google_project_service` resource that activates `compute.googleapis.com` on the target project
- **VPC Network** — a `google_compute_network` resource with the specified name, subnet mode, routing mode, and depth settings in the target project

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId`
- **IAM permissions** to enable APIs and create VPC networks in the target project

## Quick Start

Create a file `vpc.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVpcNetwork
metadata:
  name: my-vpc
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpVpcNetwork.my-vpc
spec:
  projectId:
    value: my-gcp-project-123
  networkName: dev-network
```

Deploy:

```shell
planton apply -f vpc.yaml
```

This creates a custom-mode VPC named `dev-network` with regional routing in the specified GCP project.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `networkName` | `string` | Name of the VPC network in GCP. | 1-63 chars, lowercase letters/numbers/hyphens, must start with a letter and end with a letter or number |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project that owns the VPC. Can reference a GcpProject resource via `valueFrom`. |
| `autoCreateSubnetworks` | `bool` | `false` | When `true`, GCP automatically creates a subnet in every region. When `false` (recommended), subnets are managed separately via GcpSubnetwork resources. |
| `routingMode` | `enum` | `REGIONAL` | Dynamic routing mode for Cloud Routers: `REGIONAL` (routes advertised in one region only) or `GLOBAL` (routes advertised across all regions). Use `GLOBAL` for multi-region or hybrid connectivity. |
| `description` | `string` | — | Human-readable description. Immutable — changing it recreates the network. |
| `mtu` | `int32` | `1460` | Maximum Transmission Unit in bytes, 1300–8896 (jumbo frames). |
| `enableUlaInternalIpv6` | `bool` | `false` | Assigns a /48 ULA internal IPv6 range from fd20::/20 to the network. |
| `internalIpv6Range` | `string` | auto-allocated | Explicit /48 from fd20::/20 when ULA IPv6 is enabled. Immutable. |
| `networkFirewallPolicyEnforcementOrder` | `string` | `AFTER_CLASSIC_FIREWALL` | Order of firewall policy vs classic rule evaluation (`BEFORE_CLASSIC_FIREWALL` or `AFTER_CLASSIC_FIREWALL`). |
| `networkProfile` | `string` | — | Full or partial URL of a network profile applied at creation. Immutable. |
| `bgpBestPathSelection.mode` | `string` | `LEGACY` | BGP best-path selection algorithm (`LEGACY` or `STANDARD`). |
| `bgpBestPathSelection.alwaysCompareMed` | `bool` | `false` | Compare MED across routes from different neighbor ASNs (STANDARD mode only). |
| `bgpBestPathSelection.interRegionCost` | `string` | `DEFAULT` | Inter-regional cost behavior (`DEFAULT` or `ADD_COST_TO_MED`). |
| `deleteDefaultRoutesOnCreate` | `bool` | `false` | Suppress the automatic 0.0.0.0/0 routes at creation. Immutable. |

## Examples

### Custom-Mode VPC with Regional Routing

A basic VPC for a single-region deployment:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVpcNetwork
metadata:
  name: dev-vpc
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.GcpVpcNetwork.dev-vpc
spec:
  projectId:
    value: my-dev-project-123
  networkName: dev-network
```

### Multi-Region VPC with Global Routing

A VPC with global dynamic routing for multi-region workloads or hybrid VPN/Interconnect setups:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVpcNetwork
metadata:
  name: prod-vpc
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpVpcNetwork.prod-vpc
spec:
  projectId:
    value: my-prod-project-456
  networkName: prod-network
  routingMode: GLOBAL
```

### Jumbo-Frame VPC with Hardened Routing Defaults

A VPC tuned for high-throughput east-west traffic that also suppresses the
automatic default routes (workloads reach the internet only through an
explicitly created NAT path):

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpVpcNetwork
metadata:
  name: data-vpc
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.GcpVpcNetwork.data-vpc
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      field: status.outputs.project_id
  networkName: data-network
  mtu: 8896
  deleteDefaultRoutesOnCreate: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `networkSelfLink` | `string` | Full self-link URL of the created VPC network (e.g., `https://www.googleapis.com/compute/v1/projects/my-project/global/networks/my-vpc`) |
| `networkName` | `string` | Name of the VPC network |
| `networkId` | `string` | Self-link identifier of the network (e.g., `projects/my-project/global/networks/my-vpc`) |
| `gatewayIpv4` | `string` | IPv4 address of the network's default internet gateway |
| `internalIpv6Range` | `string` | ULA internal IPv6 range assigned when ULA IPv6 is enabled |

## Related Components

- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project
- [GcpSubnetwork](/docs/catalog/gcp/subnetwork) — creates subnets within this VPC with primary and secondary IP ranges
- [GcpGlobalAddress](/docs/catalog/gcp/global-address) — reserves the VPC_PEERING range used for private services access
- [GcpServiceNetworkingConnection](/docs/catalog/gcp/service-networking-connection) — peers this VPC with Google's managed-service producer network
- [GcpAddress](/docs/catalog/gcp/regional-address) — reserves regional static IPs inside this VPC's subnets
- [GcpRouterNat](/docs/catalog/gcp/router-nat) — provides Cloud NAT for private workload outbound internet access
- [GcpGkeCluster](/docs/catalog/gcp/gke-cluster) — deploys a GKE cluster into this VPC
- [GcpCloudSql](/docs/catalog/gcp/cloud-sql) — deploys Cloud SQL instances that use the composed private-services-access pair for private IP connectivity

---
title: "Subnetwork"
description: "Subnetwork deployment documentation"
icon: "package"
order: 100
componentName: "gcpsubnetwork"
---

# GCP Subnetwork

Deploys a custom-mode subnet in an existing GCP VPC with configurable primary CIDR range, optional secondary IP ranges for GKE alias IPs, and Private Google Access for internal API connectivity. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects and VPCs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Engine API** -- enabled on the target project before creating subnet resources
- **Subnetwork** -- a `compute.Subnetwork` in the specified project, region, and VPC network with the configured primary CIDR range and Private Google Access setting
- **Secondary IP Ranges** -- created only when `secondaryIpRanges` entries are provided; named secondary CIDR blocks used by GKE for Pod and Service IP allocation
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the subnetwork will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **A VPC network** in custom subnet mode. Provide the network self-link directly or reference a GcpVpcNetwork Cloud Resource via ValueFromRef.
- **Compute Engine API** (`compute.googleapis.com`) enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Subnetwork**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **GKE-Ready** preset in the [Presets](#presets) tab to pre-populate a subnet with secondary IP ranges configured for GKE Pod and Service allocation.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSubnetwork
metadata:
  name: gke-subnet
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  vpcSelfLink:
    value: "projects/acme-prod-12345/global/networks/main-vpc"
  region: us-central1
  ipCidrRange: "10.0.0.0/20"
  subnetworkName: "gke-us-central1"
  privateIpGoogleAccess: true
```

```shell
planton apply -f gcp-subnetwork.yaml
```

This creates a subnet with a `/20` primary range and Private Google Access enabled. No secondary ranges are configured -- add them for GKE deployments. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the subnetwork to a GCP project and VPC deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
  vpcSelfLink:
    valueFrom:
      kind: GcpVpcNetwork
      name: main-vpc
      fieldPath: status.outputs.network_self_link
```

The InfraPipeline resolves the dependency graph, deploys the project and VPC first, then provisions the subnetwork with the resolved values.

## Key Configuration

These are the most important decisions when configuring a subnetwork. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Primary CIDR range** -- The `ipCidrRange` must be non-overlapping with other subnets in the VPC. Use `/20` (4,096 IPs) for GKE node subnets or `/24` (256 IPs) for smaller compute workloads. The range is immutable after creation.

**Secondary IP ranges for GKE** -- Add entries to `secondaryIpRanges` with names like `pods` and `services` to enable VPC-native GKE clusters. A `/14` range for pods supports large clusters, while a `/20` range for services covers most Kubernetes deployments. The range names are referenced by GcpGkeCluster's `clusterSecondaryRangeName` and `servicesSecondaryRangeName` fields.

**Private Google Access** -- Set `privateIpGoogleAccess` to true so VMs and GKE nodes without external IPs can reach Google APIs (Container Registry, Cloud Logging, Cloud Storage) through internal routing. Required for private GKE clusters.

**Region selection** -- The `region` field is immutable after creation. Place subnets in regions close to your users or adjacent to other resources that need low-latency connectivity.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpVpcNetwork** | `vpcSelfLink` | `status.outputs.network_self_link` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `subnetwork_self_link` | Self-link URL of the created subnetwork | GcpGkeCluster subnet reference, Compute Engine instance placement |
| `region` | Region where the subnetwork resides | Cloud Router and NAT region alignment |
| `ip_cidr_range` | Primary IPv4 CIDR of the subnet | Network planning, firewall rule source/destination ranges |
| `secondary_ranges` | List of secondary range names and CIDRs | GcpGkeCluster Pod and Service secondary range name references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**GKE-ready subnet** -- A `/20` primary range with two secondary ranges: a `/14` for pods (262,144 IPs) and a `/20` for services (4,096 IPs), plus Private Google Access enabled. The range names `pods` and `services` match the defaults expected by GcpGkeCluster presets. Start from the **GKE-Ready** preset.

**General-purpose subnet** -- A `/24` primary range for Compute Engine VMs or Cloud Run with VPC Access. No secondary ranges, since alias IPs are not needed outside GKE. Private Google Access is still enabled. Start from the **General-Purpose** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the subnetwork is created
- [**GCP VPC**](/cloud-catalog/gcp-vpc) -- provides the VPC network that contains this subnetwork
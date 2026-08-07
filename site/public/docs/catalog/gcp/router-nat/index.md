---
title: "Router NAT"
description: "Router NAT deployment documentation"
icon: "package"
order: 100
componentName: "gcprouternat"
---

# GCP Router NAT

Deploys a Cloud Router with a NAT gateway that provides outbound internet connectivity for VMs and GKE nodes without external IP addresses. The NAT can cover all subnets in a region or be scoped to specific subnets, with automatic or static IP allocation and configurable translation logging. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects and VPCs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloud Router** -- a `compute.Router` in the specified project, region, and VPC network that hosts the NAT configuration
- **Cloud NAT Gateway** -- a `compute.RouterNat` attached to the router, configured with the selected IP allocation strategy and subnet scope
- **Static External IPs** -- created only when `natIpNames` entries are provided; regional `compute.Address` resources used as predictable NAT egress IPs
- **Subnet Scoping** -- when `subnetworkSelfLinks` entries are provided, NAT applies only to the listed subnets; otherwise all subnets in the region are covered
- **NAT Logging** -- translation logging configured with the selected filter (ERRORS_ONLY by default, ALL for auditing, or DISABLED)
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the Cloud Router and NAT will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **A VPC network** with at least one subnet in the target region. Provide the network self-link directly or reference a GcpVpcNetwork Cloud Resource via ValueFromRef.
- **Compute Engine API** (`compute.googleapis.com`) enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Router NAT**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **All-Subnets Auto** preset in the [Presets](#presets) tab to pre-populate a NAT configuration covering all subnets with auto-allocated IPs.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpRouterNat
metadata:
  name: main-nat
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  vpcSelfLink:
    value: "projects/acme-prod-12345/global/networks/main-vpc"
  region: us-central1
  routerName: "main-router"
  natName: "main-nat"
```

```shell
planton apply -f gcp-router-nat.yaml
```

This creates a Cloud Router and NAT covering all subnets in us-central1 with auto-allocated IPs and error-only logging. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Cloud NAT to a GCP project and VPC deployed in the same InfraPipeline:

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

The InfraPipeline resolves the dependency graph, deploys the project and VPC first, then provisions the Cloud Router and NAT with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Cloud NAT. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Subnet scope** -- Leave `subnetworkSelfLinks` empty to cover all subnets in the region (the common case). Provide specific subnet self-links when only certain subnets should have NAT egress -- this is useful for isolating workloads that should not reach the internet.

**IP allocation** -- Leave `natIpNames` empty for auto-allocated IPs that GCP manages automatically. Provide static IP names when you need predictable egress addresses for partner allowlisting, compliance requirements, or API rate limit registration. Static IPs are created as regional `compute.Address` resources.

**Log filter** -- `logFilter` defaults to ERRORS_ONLY, which logs port exhaustion and connection failures without high volume. Use ALL for security auditing or detailed troubleshooting. Use DISABLED for non-production environments to reduce costs.

**Router and NAT naming** -- The `routerName` and `natName` fields set the GCP resource names. Each region in a VPC needs its own Cloud Router. Multiple NAT configurations can share a router, but each must have a unique name.

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
| `name` | Name of the Cloud NAT gateway | Inventory/automation references (private GKE clusters compose Cloud NAT on their network — no spec field links them) |
| `router_self_link` | Self-link URL of the Cloud Router | Network topology documentation, dependency ordering |
| `nat_ips` | Self links of the manual NAT IP reservations (empty under auto-allocation) | Partner allowlisting and compliance records — the literal IPs live on the referenced GcpAddress nodes' `address` outputs |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**All-subnets auto-allocated NAT** -- Covers all subnets in the region with auto-allocated external IPs and error-only logging. The simplest configuration for giving private GKE nodes or VMs outbound internet access. Start from the **All-Subnets Auto** preset.

**Static IP NAT for specific subnets** -- Restricts NAT to listed subnets with manually assigned static external IPs. Use when egress must come from known, stable IP addresses for partner allowlisting or compliance. Start from the **Static IP Specific Subnets** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the Cloud Router and NAT are created
- [**GCP VPC**](/cloud-catalog/gcp-vpc) -- provides the VPC network that the Cloud Router is attached to
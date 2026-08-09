---
title: "VPC"
description: "VPC deployment documentation"
icon: "package"
order: 100
componentName: "civovpc"
---

# VPC on Civo

Deploys an isolated private network (VPC) on Civo Cloud with configurable CIDR allocation and regional placement. Civo VPCs provide network isolation for compute instances, managed databases, and Kubernetes clusters running in the same region. Integrates with Planton's Provider Connections for Civo credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Network** -- a private network in the specified Civo region with the configured label and optional CIDR block. When no CIDR is specified, Civo auto-allocates an available range.
- **Civo Tags** -- metadata tags applied to the network for organizational tracking

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo account** with API access enabled in the target region. No additional prerequisites are required -- VPCs are foundational resources with no upstream dependencies.

## Deploy

### Console

Open the deployment store, find **VPC on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoVpc
metadata:
  name: app-network
  org: acme-corp
  env: prod
spec:
  networkName: app-network
  region: lon1
  ipRangeCidr: "10.0.0.0/24"
```

```shell
planton apply -f civo-vpc.yaml
```

This creates a private network in Civo's London region with a /24 CIDR block providing 254 usable addresses. No default-for-region flag is set. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Civo VPC. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CIDR range** -- The `ipRangeCidr` field sets the IPv4 address block for the network (maximum /24 on Civo). Omit it to let Civo auto-allocate an available range. Specify an explicit CIDR when you need predictable addressing for firewall rules or peering configurations.

**Region** -- The `region` field accepts a Civo region code (e.g., `lon1` for London, `nyc1` for New York, `fra1` for Frankfurt). All resources attached to this VPC must be in the same region.

**Default network** -- Set `isDefaultForRegion` to `true` to make this the default network for the region. Only one network per region can be the default. New instances without an explicit network assignment use the default.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `network_id` | Unique identifier of the network on Civo | CivoKubernetesCluster, CivoComputeInstance, CivoDatabase, CivoFirewall network references |
| `cidr_block` | IPv4 CIDR block of the created network | Firewall rule CIDR configuration, network planning |
| `created_at_rfc3339` | Network creation timestamp in RFC 3339 format | Audit logs, lifecycle tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard VPC** -- a private network with an explicit /24 CIDR range for workload isolation. Suitable for most projects where instances, databases, and Kubernetes clusters need private connectivity. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.
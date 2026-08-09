---
title: "VPC"
description: "VPC deployment documentation"
icon: "package"
order: 100
componentName: "scalewayvpc"
---

# Scaleway VPC

Deploys a regional VPC on Scaleway that serves as a logical container for grouping Private Networks. Configurable inter-network routing enables communication between Private Networks within the same VPC, supporting multi-tier architectures where workloads in separate networks need connectivity. Supports ValueFromRef for dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPC** -- a regional `network.Vpc` in the specified Scaleway region with the configured routing mode and custom routes propagation settings
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair (Access Key + Secret Key). The IaC module authenticates through the Scaleway provider configuration.
- **Choose a region** -- VPCs are regional resources. All Private Networks attached to this VPC must reside in the same region. Available regions include `fr-par` (Paris), `nl-ams` (Amsterdam), and `pl-waw` (Warsaw).

## Deploy

### Console

Open the deployment store, find **Scaleway VPC**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a basic VPC with routing disabled.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayVpc
metadata:
  name: main-vpc
  org: acme-corp
  env: prod
spec:
  region: fr-par
```

```shell
planton apply -f scaleway-vpc.yaml
```

This creates a VPC in the Paris region with routing disabled and no custom routes propagation. Private Networks can be added to this VPC using ScalewayPrivateNetwork. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a VPC. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Region** -- The `region` field determines where the VPC is created (`fr-par`, `nl-ams`, or `pl-waw`). All Private Networks attached to this VPC must be in the same region. Cannot be changed after creation.

**Inter-network routing** -- Set `enableRouting` to true when Private Networks in the VPC need to communicate with each other. Required for multi-tier architectures where a Kapsule cluster in one Private Network reaches an RDB instance or Redis cluster in another. Once enabled, routing cannot be disabled -- plan accordingly.

**Custom routes propagation** -- Set `enableCustomRoutesPropagation` to true when using VPN gateways or network appliances that advertise routes across Private Networks. Only relevant when routing is also enabled. Like routing, this is a one-way toggle -- once enabled, it cannot be disabled.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpc_id` | Unique identifier (UUID) of the created VPC | ScalewayPrivateNetwork VPC reference, network topology tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard VPC** -- A VPC in the Paris region with routing disabled. The simplest starting configuration for single-tier environments or architectures with a single Private Network where cross-network communication is unnecessary. Start from the **Standard** preset.

**Routing-enabled VPC** -- A VPC with inter-Private-Network routing enabled for multi-tier production architectures. Resources in separate Private Networks (e.g., application tier, database tier, cache tier) can communicate over private IPs. Start from the **Routing-Enabled** preset.

## Works With

This component operates independently and does not reference other components.
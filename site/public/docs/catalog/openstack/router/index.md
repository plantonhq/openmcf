---
title: "Router"
description: "Router deployment documentation"
icon: "package"
order: 100
componentName: "openstackrouter"
---

# OpenStack Router

Deploys a Neutron router on OpenStack, providing Layer 3 routing between subnets and optional external network connectivity via SNAT. The router references an external OpenStackNetwork for its gateway and supports Distributed Virtual Router (DVR) mode for eliminating centralized routing bottlenecks.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Router** -- a Layer 3 router with configurable admin state, external gateway, SNAT, DVR mode, and external fixed IP allocations
- **External Gateway** -- created only when `externalNetworkId` is provided; connects the router to a provider network for outbound connectivity and floating IP allocation
- **OpenStack Tags** -- user-defined tags applied to the router for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **External network** (optional) -- for outbound internet access, the router needs an external (provider) network as its gateway. Obtain the external network ID from your OpenStack administrator or reference an OpenStackNetwork Cloud Resource via ValueFromRef.
- **Subnet connectivity** -- routers provide routing between subnets. After creating the router, attach subnets using OpenStackRouterInterface resources to enable inter-subnet traffic.
- **DVR support** -- if you plan to use `distributed: true`, confirm your OpenStack deployment supports DVR (requires the L3 agent configured in DVR mode on compute nodes).

## Deploy

### Console

Open the deployment store, find **OpenStack Router**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Edge with SNAT** preset in the [Presets](#presets) tab to pre-populate a production-ready configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackRouter
metadata:
  name: edge-router
  org: acme-corp
  env: prod
spec:
  externalNetworkId:
    value: "<external-network-id>"
  enableSnat: true
```

```shell
planton apply -f router.yaml
```

This creates a router with an external gateway and SNAT enabled. Tenant instances on connected subnets can reach the internet through this router without individual floating IPs. DVR mode and external fixed IPs are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the router to an external network deployed in the same InfraPipeline:

```yaml
spec:
  externalNetworkId:
    valueFrom:
      kind: OpenStackNetwork
      name: provider-net
      fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph, deploys the external network first, then provisions the router with the resolved network ID.

## Key Configuration

These are the most important decisions when configuring a router. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**External gateway** -- Set `externalNetworkId` to connect the router to a provider network. Without an external gateway, the router provides internal routing only (east-west traffic between subnets). This is the fundamental choice between an edge router and an internal router.

**SNAT** -- `enableSnat` controls whether tenant traffic is NATed to the router's external IP for outbound connectivity. Only valid when an external gateway is configured. Most production routers enable SNAT so instances can reach the internet without individual floating IPs.

**Distributed Virtual Router** -- Set `distributed: true` to distribute routing to each compute node, eliminating the centralized L3 agent as a bottleneck. This is a create-time setting and cannot be changed after creation. DVR improves east-west traffic performance but adds operational complexity.

**External fixed IPs** -- `externalFixedIps` lets you request specific IP addresses or subnets on the external network for the router's gateway. If omitted, OpenStack auto-allocates from the external network. Use this when you need a predictable external IP for DNS or firewall rules.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackNetwork** (optional) | `externalNetworkId` | `status.outputs.network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `router_id` | UUID of the router in OpenStack | Router interface attachments |
| `name` | Name of the router | DNS records, monitoring labels |
| `external_network_id` | ID of the external gateway network (empty if none) | Network topology reference |
| `external_gateway_ip` | Primary external IP allocated to the router's gateway (empty if none) | Floating IP associations, firewall rules |
| `region` | OpenStack region where the router was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Edge router with SNAT** -- Connects to an external network with SNAT enabled. Tenant instances on attached subnets can reach the internet through this router without individual floating IPs. The standard production router configuration. Start from the **Edge with SNAT** preset.

**Internal-only router** -- No external gateway. Provides Layer 3 routing between connected subnets within the tenant. Use for isolated environments, air-gapped networks, or when external access is handled by a separate shared router. Start from the **Internal-Only** preset.

## Works With

- [**OpenStack Network**](/cloud-catalog/openstack-network) -- provides the external (provider) network used as the router's gateway
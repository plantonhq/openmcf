---
title: "Router Interface"
description: "Router Interface deployment documentation"
icon: "package"
order: 100
componentName: "openstackrouterinterface"
---

# OpenStack Router Interface

Attaches an OpenStack Neutron subnet to a router, enabling Layer 3 routing for that subnet. This is a join resource -- it creates a port on the subnet and binds it to the router. Without a router interface, a subnet is an isolated Layer 2 domain with no connectivity to other subnets or external networks. ValueFromRef wiring connects both the router and subnet from their respective Cloud Resources in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Router Interface** -- a port on the specified subnet attached to the specified router, enabling Layer 3 routing between the subnet and any other subnets or external networks connected to that router

## Before You Deploy

### OpenStack Account

- **Router** -- an existing Neutron router to attach the subnet to. Provide the router ID directly or reference an OpenStackRouter Cloud Resource via ValueFromRef.
- **Subnet** -- an existing subnet to connect to the router. Provide the subnet ID directly or reference an OpenStackSubnet Cloud Resource via ValueFromRef.
- **Subnet CIDR conflicts** -- the subnet's CIDR must not overlap with any other subnet already attached to the same router. Overlapping CIDRs cause a creation error.

## Deploy

### Console

Open the deployment store, find **OpenStack Router Interface**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackRouterInterface
metadata:
  name: app-subnet-attach
  org: acme-corp
  env: prod
spec:
  routerId:
    value: "<router-id>"
  subnetId:
    value: "<subnet-id>"
```

```shell
planton apply -f router-interface.yaml
```

This attaches the subnet to the router, enabling instances on the subnet to reach other connected subnets and (if the router has an external gateway) the internet. No region override is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire both the router and subnet deployed in the same InfraPipeline:

```yaml
spec:
  routerId:
    valueFrom:
      kind: OpenStackRouter
      name: edge-router
      fieldPath: status.outputs.router_id
  subnetId:
    valueFrom:
      kind: OpenStackSubnet
      name: app-subnet
      fieldPath: status.outputs.subnet_id
```

The InfraPipeline resolves the dependency graph, deploys the router and subnet first, then provisions the router interface with the resolved IDs.

## Key Configuration

These are the most important decisions when configuring a router interface. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Router selection** -- The `routerId` determines which router provides Layer 3 connectivity. For internet access, reference a router with an external gateway. For internal-only routing, reference a router without an external gateway.

**Subnet selection** -- The `subnetId` determines which subnet gains routing. Each subnet can be attached to only one router. Attaching the same subnet to multiple routers is not supported by OpenStack.

**Immutability** -- All fields are ForceNew. Changing the router or subnet reference recreates the interface, which briefly disrupts routing for instances on the subnet.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackRouter** | `routerId` | `status.outputs.router_id` |
| **OpenStackSubnet** | `subnetId` | `status.outputs.subnet_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `port_id` | UUID of the port created by the router interface attachment | Debugging, network topology inspection |
| `router_id` | UUID of the router this interface is attached to | Audit, topology reference |
| `subnet_id` | UUID of the subnet connected to the router | Audit, topology reference |
| `region` | OpenStack region where the router interface was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard attachment** -- Attaches a subnet to a router with no additional configuration. This is the only pattern -- every router interface is a simple binding between a router and a subnet. Start from the **Standard** preset.

## Works With

- [**OpenStack Router**](/cloud-catalog/openstack-router) -- provides the router ID that the subnet is attached to
- [**OpenStack Subnet**](/cloud-catalog/openstack-subnet) -- provides the subnet ID connected to the router
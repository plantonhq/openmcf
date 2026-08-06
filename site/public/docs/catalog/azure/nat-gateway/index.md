---
title: "NAT Gateway"
description: "NAT Gateway deployment documentation"
icon: "package"
order: 100
componentName: "azurenatgateway"
---

# Azure NAT Gateway

Deploys an Azure NAT Gateway — the managed source-network-address-translation (SNAT) service that gives every workload in its attached subnets stable, scalable outbound internet connectivity. The gateway is deliberately just the gateway: the public IPs and prefixes it SNATs through are referenced first-class resources, and subnets attach themselves to it via `AzureSubnet`'s `natGatewayId`.

## What Gets Created

When you deploy an AzureNatGateway resource, Planton provisions:

- **NAT Gateway** — a `network.NatGateway` in the specified region and resource group, with the configured SKU (Standard or StandardV2), idle timeout, and availability zone
- **Public IP Associations** — a `NatGatewayPublicIpAssociation` for each referenced `AzurePublicIp`, binding that address to the gateway for SNAT
- **Public IP Prefix Associations** — a `NatGatewayPublicIpPrefixAssociation` for each referenced `AzurePublicIpPrefix`, binding that contiguous range to the gateway for SNAT
- **Azure Tags** — resource metadata tags applied to the gateway for tracking and governance

Nothing else is created here. The gateway does not allocate its own IPs or prefixes (reference existing ones), and it does not modify subnets (each `AzureSubnet` declares the attachment itself via `natGatewayId`, matching Azure's model where one gateway serves many subnets).

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the NAT Gateway will be created (can reference an AzureResourceGroup resource)
- **At least one AzurePublicIp or AzurePublicIpPrefix** to associate — a gateway with no addresses deploys but cannot translate anything (a StandardV2 gateway needs StandardV2 addresses)

## Quick Start

Create a file `natgateway.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: my-natgw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureNatGateway.my-natgw
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: egress-nat
  publicIpIds:
    - valueFrom:
        name: my-egress-ip
```

Deploy:

```shell
planton apply -f natgateway.yaml
```

This creates a Standard SKU NAT Gateway that SNATs through the referenced public IP, with Azure's default 4-minute idle timeout. To route a subnet's outbound traffic through it, set `natGatewayId` on that `AzureSubnet` to this gateway's `status.outputs.nat_gateway_id`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the NAT Gateway (e.g., `eastus`, `westeurope`). A gateway only serves subnets in its own region. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `name` | `string` | Name of the NAT Gateway, unique within the resource group. | Required, 1-80 chars, Azure naming rules |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `skuName` | `enum` | `STANDARD` | `STANDARD` (zonal, optionally pinned to one zone) or `STANDARD_V2` (zone-redundant automatically; `zones` must be empty and referenced IPs/prefixes must be StandardV2). Fixed at creation. |
| `idleTimeoutInMinutes` | `int32` | `4` | How long an idle outbound TCP connection's SNAT port stays reserved. Range: 4--120. Raise it only for long-lived idle connections; higher values hold ports longer and hasten SNAT exhaustion. |
| `zones` | `string[]` | `[]` | Availability zone to pin a `STANDARD` gateway to (e.g., `["1"]`). Empty means non-zonal. Must be empty for `STANDARD_V2`. Fixed at creation. |
| `publicIpIds` | `StringValueOrRef[]` | `[]` | ARM IDs of public IPs the gateway SNATs through. Each address adds 64,512 SNAT ports. Defaults to referencing an `AzurePublicIp`'s `public_ip_id` output. |
| `publicIpPrefixIds` | `StringValueOrRef[]` | `[]` | ARM IDs of public IP prefixes (contiguous reserved ranges) the gateway SNATs through — the scalable, allowlistable alternative to individual addresses. Defaults to referencing an `AzurePublicIpPrefix`'s `public_ip_prefix_id` output. |
| `tags` | `map<string, string>` | `{}` | Additional tags applied to the NAT Gateway, merged over Planton-derived tags (user wins on collision). |

## Examples

### Single Public IP

A zonal NAT Gateway with default settings, SNATing through one referenced public IP:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: basic-natgw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureNatGateway.basic-natgw
spec:
  region: eastus
  resourceGroup:
    value: dev-rg
  name: dev-egress-nat
  zones:
    - "1"
  publicIpIds:
    - valueFrom:
        name: dev-egress-ip
```

### Custom Idle Timeout with Tags

A NAT Gateway with a longer idle timeout for workloads that hold long-lived outbound TCP connections:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: long-lived-natgw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNatGateway.long-lived-natgw
spec:
  region: westeurope
  resourceGroup:
    value: prod-rg
  name: prod-egress-nat
  idleTimeoutInMinutes: 30
  publicIpIds:
    - valueFrom:
        name: prod-egress-ip
  tags:
    team: platform
    cost-center: infra
```

### Public IP Prefix for Scale

A NAT Gateway backed by a referenced /28 Public IP Prefix (16 addresses, over 1M SNAT ports) for high-throughput subnets that need multiple outbound IPs to avoid SNAT port exhaustion:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: scale-natgw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNatGateway.scale-natgw
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: prod-scale-nat
  zones:
    - "1"
  idleTimeoutInMinutes: 10
  publicIpPrefixIds:
    - valueFrom:
        name: prod-snat-prefix
  tags:
    workload: high-throughput
```

### Zone-Redundant StandardV2

A next-generation StandardV2 gateway — zone-redundant automatically, so `zones` stays empty. The referenced address must also be StandardV2:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: v2-natgw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNatGateway.v2-natgw
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: prod-egress-nat-v2
  skuName: STANDARD_V2
  publicIpIds:
    - valueFrom:
        name: prod-v2-egress-ip
```

### Production AKS Cluster Egress

A NAT Gateway serving as the outbound egress for an AKS cluster node subnet, with a referenced prefix for SNAT scale and a 120-minute idle timeout for long-running batch jobs. The node subnet attaches itself by setting `natGatewayId`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: aks-egress-natgw
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNatGateway.aks-egress-natgw
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: aks-rg
  name: aks-egress-nat
  zones:
    - "1"
  idleTimeoutInMinutes: 120
  publicIpPrefixIds:
    - valueFrom:
        name: aks-snat-prefix
  tags:
    purpose: aks-egress
    environment: production
---
apiVersion: azure.planton.dev/v1alpha1
kind: AzureSubnet
metadata:
  name: aks-nodes
spec:
  virtualNetworkId:
    valueFrom:
      name: aks-vnet
  name: aks-nodes
  addressPrefixes:
    - 10.0.1.0/24
  natGatewayId:
    valueFrom:
      name: aks-egress-natgw
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `nat_gateway_id` | `string` | Azure Resource Manager ID of the NAT Gateway — the primary output; `AzureSubnet`'s `natGatewayId` references it to attach the gateway to a subnet |
| `nat_gateway_name` | `string` | The NAT Gateway's name as deployed |
| `resource_guid` | `string` | The immutable GUID ARM assigns the gateway — useful when correlating with Azure billing, monitoring, or support data keyed on the GUID |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/resource-group) -- provides the resource group for gateway placement
- [AzurePublicIp](/docs/catalog/azure/public-ip) -- the individual public addresses the gateway SNATs through, referenced via `publicIpIds`
- [AzurePublicIpPrefix](/docs/catalog/azure/public-ip-prefix) -- contiguous, allowlistable address ranges the gateway SNATs through, referenced via `publicIpPrefixIds`
- [AzureSubnet](/docs/catalog/azure/subnet) -- attaches the gateway via its `natGatewayId` field so all outbound traffic from the subnet routes through it
- [AzureAksCluster](/docs/catalog/azure/aks-cluster) -- AKS clusters commonly use a NAT Gateway for predictable outbound IPs and SNAT port scaling

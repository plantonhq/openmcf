---
title: "Public IP Prefix"
description: "Public IP Prefix deployment documentation"
icon: "package"
order: 100
componentName: "azurepublicipprefix"
---

# Azure Public IP Prefix

Creates an Azure Public IP Prefix -- a reserved, contiguous range of public IP addresses. Individual public IPs drawn from the prefix share one allowlistable CIDR, and NAT gateways associate whole prefixes for scalable outbound SNAT.

## What Gets Created

When you deploy an AzurePublicIpPrefix resource, Planton provisions:

- **Public IP Prefix** — an `azurerm_public_ip_prefix` reserving a contiguous public address range in the specified region and resource group

Downstream association is deliberately not part of this resource: `AzurePublicIp` allocates individual addresses from the prefix (`publicIpPrefixId`), and `AzureNatGateway` associates whole prefixes for SNAT (`publicIpPrefixIds`).

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the prefix in (an `AzureResourceGroup` in composed environments)
- **Network write rights**: `Microsoft.Network/publicIPPrefixes/write` (Network Contributor, Contributor, or Owner)

## Quick Start

Create a file `public-ip-prefix.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePublicIpPrefix
metadata:
  name: prod-egress-prefix
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIpPrefix.prod-egress-prefix
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  name: prod-egress-prefix
  prefixLength: 28
  zones:
    - "1"
    - "2"
    - "3"
```

Deploy:

```shell
planton apply -f public-ip-prefix.yaml
```

After deployment, read `status.outputs.ip_prefix` for the CIDR to give partners and firewalls, and `status.outputs.public_ip_prefix_id` for downstream wiring.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region; must match every resource that consumes the prefix. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Prefix name, unique within the resource group. | Required, 1-80 chars, Azure naming rules |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `prefixLength` | `int32` | `28` | CIDR length of the range to reserve. IPv4 spans /21 (2,048 addresses) to /31 (2). Fixed at creation; you pay for every reserved address. |
| `ipVersion` | `enum` | `IPV4` | `IPV4` or `IPV6`. Fixed at creation. |
| `sku` | `enum` | `STANDARD` | `STANDARD` or `STANDARD_V2`. StandardV2 is required for StandardV2 NAT gateways. Fixed at creation. |
| `skuTier` | `enum` | `REGIONAL` | `REGIONAL` (NAT gateway, most workloads) or `GLOBAL` (cross-region load balancer frontends only; requires `STANDARD` SKU). Fixed at creation. |
| `zones` | `string[]` | `[]` | Availability zones. `["1", "2", "3"]` for zone-redundant; a single zone pins the range; omit to let Azure choose. Fixed at creation. |
| `customIpPrefixId` | `string` | `""` | ARM ID of a Custom IP Prefix for BYOIP carve-out. Omit to allocate from Microsoft's pool. Fixed at creation. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). The only field that updates in place. |

## Examples

### NAT Gateway SNAT Range

A zone-redundant /28 (16 addresses, ~1M SNAT ports) for production outbound:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePublicIpPrefix
metadata:
  name: nat-snat-prefix
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIpPrefix.nat-snat-prefix
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: nat-snat-prefix
  prefixLength: 28
  zones:
    - "1"
    - "2"
    - "3"
```

Reference it from a NAT gateway:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureNatGateway
metadata:
  name: prod-nat
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: prod-nat
  publicIpPrefixIds:
    - valueFrom:
        name: nat-snat-prefix
```

### Partner Firewall Allowlist

A smaller /30 (4 addresses) when partners need a tight, cheap allowlist surface:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePublicIpPrefix
metadata:
  name: partner-allowlist
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIpPrefix.partner-allowlist
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: partner-allowlist
  prefixLength: 30
  zones:
    - "1"
    - "2"
    - "3"
```

Share `status.outputs.ip_prefix` with the partner once deployed.

### Allocate a Public IP From the Prefix

Reserve the range first, then draw individual addresses from it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePublicIp
metadata:
  name: lb-frontend-from-prefix
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: lb-frontend-from-prefix
  publicIpPrefixId:
    valueFrom:
      name: partner-allowlist
  zones:
    - "1"
    - "2"
    - "3"
```

The allocated address comes from the prefix's contiguous range and stays
allowlistable as part of `ip_prefix`.

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `public_ip_prefix_id` | `string` | Full ARM ID of the prefix -- referenced by `AzurePublicIp.publicIpPrefixId` and `AzureNatGateway.publicIpPrefixIds` |
| `ip_prefix` | `string` | The actual reserved CIDR (e.g. `20.42.0.16/28`) -- known only after creation; the value partners and firewalls allowlist |
| `public_ip_prefix_name` | `string` | The prefix's name as deployed |

## Related Components

- [AzurePublicIp](/docs/catalog/azure/public-ip) — allocates individual addresses from the prefix via `publicIpPrefixId`
- [AzureNatGateway](/docs/catalog/azure/nat-gateway) — associates whole prefixes for outbound SNAT via `publicIpPrefixIds`
- [AzureSubnet](/docs/catalog/azure/subnet) — attaches a NAT gateway that SNATs through the prefix
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for prefix placement

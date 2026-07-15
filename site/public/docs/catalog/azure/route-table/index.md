---
title: "Route Table"
description: "Route Table deployment documentation"
icon: "package"
order: 100
componentName: "azureroutetable"
---

# Azure Route Table

Creates an Azure route table -- a reusable set of user-defined routes (UDRs) that overrides Azure's default system routing for every subnet attached to it. The building block for forced tunneling, firewall egress, network virtual appliances, and black-hole routing.

## What Gets Created

When you deploy an AzureRouteTable resource, Planton provisions:

- **Route Table** — an `azurerm_route_table` with its user-defined routes managed inline (a route has no life of its own in Azure)

The subnet-side attachment is deliberately not part of this resource: a subnet declares which route table it uses (matching Azure's model), so one table serves many subnets without listing them.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the table in (an `AzureResourceGroup` in composed environments)
- **Network write rights**: `Microsoft.Network/routeTables/write` (Network Contributor, Contributor, or Owner)

## Quick Start

Create a file `route-table.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: egress-via-firewall
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRouteTable.egress-via-firewall
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  name: egress-via-firewall
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      nextHopInIpAddress: "10.0.1.4"
  bgpRoutePropagationEnabled: false
```

Deploy:

```shell
planton apply -f route-table.yaml
```

Every subnet that attaches this table now sends its internet-bound traffic to the firewall at `10.0.1.4`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region; must match the networks whose subnets attach the table. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Table name, unique within the resource group. | Required, 1-80 chars, Azure naming rules |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `routes` | `list(object)` | The user-defined routes. Each: `name`, `addressPrefix` (CIDR or Azure service tag), `nextHopType` (`VIRTUAL_NETWORK_GATEWAY` / `VNET_LOCAL` / `INTERNET` / `VIRTUAL_APPLIANCE` / `NONE`), and `nextHopInIpAddress` (required exactly when the hop is `VIRTUAL_APPLIANCE`). Empty is valid -- attach first, add routes later. |
| `bgpRoutePropagationEnabled` | `bool` | Whether routes learned from on-premises via BGP propagate into attached subnets. Azure defaults to `true`; disable for forced-tunneling designs. |
| `tags` | `map(string)` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Forced Tunneling to On-Premises

Send everything through the VPN/ExpressRoute gateway and stop learned routes from bypassing it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: forced-tunnel
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRouteTable.forced-tunnel
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: forced-tunnel
  routes:
    - name: default-on-premises
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_NETWORK_GATEWAY
  bgpRoutePropagationEnabled: false
```

### Black-Hole Unwanted Prefixes

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: deny-lateral
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRouteTable.deny-lateral
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: deny-lateral
  routes:
    - name: blackhole-data-tier
      addressPrefix: "10.0.20.0/24"
      nextHopType: NONE
```

### Service Tag Routing

Route a whole Azure service's prefix set directly to the internet while everything else goes through the firewall:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: firewall-with-backup-bypass
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRouteTable.firewall-with-backup-bypass
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: firewall-with-backup-bypass
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      nextHopInIpAddress: "10.0.1.4"
    - name: backup-direct
      addressPrefix: "AzureBackup"
      nextHopType: INTERNET
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `route_table_id` | `string` | Full ARM ID of the table -- the join key subnets use to attach it |
| `route_table_name` | `string` | The table's name as deployed |

## Related Components

- [AzureVirtualNetwork](/docs/catalog/azure/virtual-network) — the network whose subnets attach this table
- [AzureSubnet](/docs/catalog/azure/subnet) — declares which route table it uses (the attach side)
- [AzureNetworkSecurityGroup](/docs/catalog/azure/network-security-group) — traffic filtering, complementary to routing

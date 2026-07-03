# AzureRouteTable

## Overview

`AzureRouteTable` provisions an Azure route table: a reusable set of
user-defined routes (UDRs) that overrides Azure's default system routing
for every subnet attached to it. It is the building block for forced
tunneling (send `0.0.0.0/0` to a firewall or on-premises gateway),
steering traffic through network virtual appliances, and black-holing
unwanted prefixes.

## Why a First-Class Resource?

A route table is real infrastructure with its own lifecycle:

- **One policy, many subnets** -- the same egress-via-firewall table is
  attached to every workload subnet; editing its routes changes routing
  for all of them at once
- **Independent lifecycle** -- routing policy evolves (new appliance IP,
  new on-premises prefixes) without touching the subnets that use it
- **Attach from the subnet side** -- matching Azure's model, a subnet
  declares which table it uses; the table never lists its consumers

Routes are folded inside the table because Azure applies them only as part
of it -- a route has no life of its own.

## Key Features

- **Full next-hop surface** -- VirtualNetworkGateway, VnetLocal, Internet,
  VirtualAppliance (with forwarding IP), and None (black-hole)
- **Service-tag destinations** -- routes accept Azure service tags
  ("AzureBackup", "ApiManagement") as well as CIDR prefixes
- **BGP propagation control** -- disable propagation of learned routes for
  forced-tunneling designs (the standard firewall-subnet hardening)
- **Validated pairing** -- the spec enforces that VIRTUAL_APPLIANCE routes
  carry a forwarding IP and other hop types do not
- **Composable** -- the resource group is referenced by name, defaulting
  to an `AzureResourceGroup`'s output in composed environments

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | Azure region (must match the networks whose subnets attach it) |
| `resource_group` | StringValueOrRef | Yes | Resource group name (defaults to an AzureResourceGroup reference) |
| `name` | string | Yes | Table name, unique within the resource group (1-80 chars) |
| `routes` | list | No | User-defined routes (name, address_prefix, next_hop_type, next_hop_in_ip_address) |
| `bgp_route_propagation_enabled` | bool | No | Whether BGP-learned routes propagate into attached subnets (Azure default: true) |
| `tags` | map | No | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `route_table_id` | Full ARM ID of the table -- the join key subnets use to attach it |
| `route_table_name` | The table's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: egress-via-firewall
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: egress-via-firewall
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      nextHopInIpAddress: "10.0.1.4"
  bgpRoutePropagationEnabled: false
```

## Lifecycle Notes

- Routes, BGP propagation, and tags update **in place** -- and take effect
  immediately for every attached subnet
- Name, region, and resource group are the table's ARM identity; changing
  any of them **replaces the table**, detaching it from every subnet until
  the replacement is re-attached
- An empty `routes` list is valid and common: attach the empty table
  first, add routes as the topology grows

---
title: "ExpressRoute Circuit Peering"
description: "ExpressRoute Circuit Peering deployment documentation"
icon: "package"
order: 100
componentName: "azureexpressroutecircuitpeering"
---

# Azure ExpressRoute Circuit Peering

Deploys an ExpressRoute circuit peering -- the BGP routing configuration that makes routes flow through a circuit. Private peering carries your VNets' address space (what an ExpressRoute-type virtual network gateway connects to); Microsoft peering carries Microsoft 365 and Azure public services. A circuit holds at most one peering of each type. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ExpressRoute Circuit Peering** -- the ARM child of the circuit (named by its type) carrying the VLAN, BGP session addressing, and type-specific configuration
- **Global Reach Connections** -- one per `connections` entry: links from this private peering to other circuits' private peerings

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An ExpressRoute circuit** -- the peering may be configured while the circuit's provider state is still NotProvisioned (ARM stores the configuration), but routes only flow after the connectivity provider completes the cross-connect and the circuit reads Provisioned.
- **The session facts from your provider**: the VLAN id and the /30 address pairs (one per physical link), plus -- for Microsoft peering -- your registered public prefixes.

## Deploy

### Console

Open the deployment store, find **Azure ExpressRoute Circuit Peering**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private Peering** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRouteCircuitPeering
metadata:
  name: hq-private-peering
  org: acme-corp
  env: prod
spec:
  resourceGroup:
    value: "network-rg"
  expressRouteCircuitName:
    valueFrom:
      kind: AzureExpressRouteCircuit
      name: hq-circuit
      fieldPath: status.outputs.express_route_circuit_name
  peeringType: AZURE_PRIVATE_PEERING
  vlanId: 100
  primaryPeerAddressPrefix: "192.168.16.0/30"
  secondaryPeerAddressPrefix: "192.168.16.4/30"
```

```shell
planton apply -f azure-express-route-circuit-peering.yaml
```

A Stack Job tracks provisioning in real time; Microsoft's ASN and edge-port identifiers appear in the outputs.

### InfraChart

In a hybrid-connectivity chart the peering sits between the circuit and the gateway: circuit → this peering → EXPRESS_ROUTE-type virtual network gateway → gateway connection, each wiring by reference.

## Key Configuration

These are the most important decisions when configuring a peering. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Peering type** -- `AZURE_PRIVATE_PEERING` for VNet connectivity (the common case) or `MICROSOFT_PEERING` for Microsoft public services. The type IS the peering's ARM identity: one of each per circuit, fixed at creation. Public peering exists only for legacy imports -- Azure deprecated it.

**Session addressing** -- the /30 pair (`primaryPeerAddressPrefix` + `secondaryPeerAddressPrefix`, one per physical link) and the provider-assigned `vlanId` (1-4094, unique on the circuit). Your router takes each /30's first usable address, Microsoft's the second.

**Microsoft peering** -- requires `microsoftPeeringConfig` with public prefixes REGISTERED to you (Microsoft validates ownership against internet routing registries before activating) and typically a `routeFilterId` -- without a route filter, Microsoft peering advertises nothing.

**Global Reach** -- `connections` entries link this private peering to other circuits' private peerings (far side by ARM id, a /29 for tunnel addressing, an authorization key when the far circuit is in another subscription).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureExpressRouteCircuit** | `expressRouteCircuitName` | `status.outputs.express_route_circuit_name` |
| **AzureExpressRouteCircuitPeering** | `connections[].peerPeeringId` | `status.outputs.express_route_circuit_peering_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `express_route_circuit_peering_id` | Azure Resource Manager ID of the peering | Another circuit's Global Reach `peerPeeringId` |
| `azure_asn` | Microsoft's BGP ASN (12076 on public Azure) | Router neighbor configuration |
| `primary_azure_port` / `secondary_azure_port` | Microsoft-edge port identifiers | Provider troubleshooting |
| `connection_ids` | Name-keyed Global Reach connection IDs | Operational tooling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private peering** -- VNet connectivity over the circuit: the standard hybrid path. Start from the **Private Peering** preset.

**Microsoft peering** -- Microsoft 365 / Azure public services with the advertisement contract. Start from the **Microsoft Peering** preset.

## Works With

- [**Azure ExpressRoute Circuit**](/cloud-catalog/azure-express-route-circuit) -- the parent circuit this peering configures
- [**Azure Virtual Network Gateway**](/cloud-catalog/azure-virtual-network-gateway) -- the EXPRESS_ROUTE-type gateway that consumes private peering
- [**Azure Virtual Network Gateway Connection**](/cloud-catalog/azure-virtual-network-gateway-connection) -- the link between a gateway and the circuit

# Azure Virtual Network Gateway Connection

Deploys a gateway connection -- the tunnel object that joins an AzureVirtualNetworkGateway to its far side: an on-premises VPN device (site-to-site IPsec, described by an AzureLocalNetworkGateway), another virtual network gateway (VNet-to-VNet), or an ExpressRoute circuit's private peering. The connection is deliberately its own resource: one gateway carries many tunnels, and each site's tunnel is added or removed without touching the gateway. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Gateway Connection** -- an IPsec, Vnet2Vnet, or ExpressRoute connection on the referenced gateway, with optional custom IPsec/IKE proposal, BGP, DPD timeout, NAT rule opt-ins, and traffic selectors
- **Azure Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with your `tags`, applied to the connection

The gateway, the site description, and the circuit are NOT created here -- they are referenced first-class resources.

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **A provisioned AzureVirtualNetworkGateway** (25-45 minutes to create -- deploy it first).
- **The far side**: an AzureLocalNetworkGateway for IPSEC, a second gateway for VNET_TO_VNET (plus a mirror connection on it with the same shared key), or a circuit id for EXPRESS_ROUTE.
- **Understand "provisioned is not connected"**: the connection object reaches Succeeded when ARM accepts it; the tunnel reaches Connected only when the far side negotiates. A tunnel stuck in Connecting means the device, key, or proposal disagrees -- not a failed deployment.

## Deploy

### Console

Open the deployment store, find **Azure Virtual Network Gateway Connection**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Site-to-Site** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetworkGatewayConnection
metadata:
  name: hq-to-azure
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    value: "network-rg"
  name: hq-to-azure
  type: IPSEC
  virtualNetworkGatewayId:
    value: "/subscriptions/.../virtualNetworkGateways/hub-vpn-gateway"
  localNetworkGatewayId:
    value: "/subscriptions/.../localNetworkGateways/hq-datacenter"
  sharedKey:
    valueFrom:
      kind: Secret
      name: hq-tunnel-psk
```

```shell
planton apply -f azure-gateway-connection.yaml
```

This creates the site-to-site tunnel in seconds (the gateway already exists). A Stack Job tracks provisioning in real time.

### InfraChart

When deploying a complete site-to-site story as one chart, wire all three resources by reference:

```yaml
# On the AzureVirtualNetworkGatewayConnection:
spec:
  virtualNetworkGatewayId:
    valueFrom:
      kind: AzureVirtualNetworkGateway
      name: hub-vpn-gateway
      fieldPath: status.outputs.virtual_network_gateway_id
  localNetworkGatewayId:
    valueFrom:
      kind: AzureLocalNetworkGateway
      name: hq-datacenter
      fieldPath: status.outputs.local_network_gateway_id
```

The InfraPipeline deploys the site description and the gateway first, then the connection that ties them together.

## Key Configuration

These are the most important decisions when configuring a connection. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Type** -- required, and it decides the required far side: IPSEC needs `localNetworkGatewayId`, VNET_TO_VNET needs `peerVirtualNetworkGatewayId` (and a mirror connection on the peer with the same key), EXPRESS_ROUTE needs `expressRouteCircuitId`. All spec-enforced.

**Shared key** -- the IPsec pre-shared key both ends must agree on. Reference a secret; or omit it and Azure generates one (readable back from the connection's shared-key API). Not applicable to ExpressRoute.

**IPsec policy** -- omit for Azure's default proposal set; pin all six algorithms plus SA bounds when the on-premises device demands exact parameters. `usePolicyBasedTrafficSelectors` (for policy-based devices behind a route-based gateway) requires a pinned policy.

**BGP** -- `bgpEnabled` exchanges routes dynamically; both ends must speak BGP (the gateway's `bgpEnabled` and the site's `bgpSettings`). `customBgpAddresses` picks APIPA endpoints per tunnel on active-active gateways.

**NAT opt-in** -- `egressNatRuleIds`/`ingressNatRuleIds` apply the gateway's NAT rules to this tunnel, using the ids the gateway publishes in its `nat_rule_ids` output.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureVirtualNetworkGateway** | `virtualNetworkGatewayId` | `status.outputs.virtual_network_gateway_id` |
| **AzureLocalNetworkGateway** (IPSEC) | `localNetworkGatewayId` | `status.outputs.local_network_gateway_id` |
| **AzureVirtualNetworkGateway** (VNET_TO_VNET peer) | `peerVirtualNetworkGatewayId` | `status.outputs.virtual_network_gateway_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `connection_id` | Azure Resource Manager ID of the connection | Diagnostics, Azure Monitor alert scoping |
| `connection_name` | Name of the connection | Operational tooling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Site-to-site tunnel** -- An IPsec connection to a described on-premises site with a secret-referenced pre-shared key. Start from the **Site-to-Site** preset.

**Pinned-proposal tunnel** -- A site-to-site connection with an explicit IPsec/IKE proposal for devices that demand exact algorithms. Start from the **Custom IPsec Policy** preset.

## Works With

- [**Azure Virtual Network Gateway**](/cloud-catalog/azure-virtual-network-gateway) -- the gateway this tunnel terminates on
- [**Azure Local Network Gateway**](/cloud-catalog/azure-local-network-gateway) -- describes the on-premises side of a site-to-site tunnel
- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- provides the resource group the connection is created in

# Azure Virtual Network Gateway

Deploys an Azure virtual network gateway -- the managed appliance that terminates hybrid connectivity (site-to-site IPsec VPN, point-to-site client VPN, VNet-to-VNet tunnels, or an ExpressRoute circuit's private peering) in a virtual network's dedicated "GatewaySubnet". The gateway is one third of the site-to-site story: an AzureLocalNetworkGateway describes each on-premises site, and an AzureVirtualNetworkGatewayConnection ties a site to this gateway as a tunnel. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Virtual Network Gateway** -- a VPN or ExpressRoute gateway of the chosen SKU in the specified region and resource group, bound to the referenced GatewaySubnet and (for VPN gateways) public IPs
- **NAT Rules** -- one `natRules` entry each, translating overlapping address space; each rule's ARM id surfaces in the `nat_rule_ids` output under its name
- **Azure Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with your `tags`, applied to the gateway

Connections are NOT created here -- each AzureVirtualNetworkGatewayConnection is its own resource pointing back at this gateway.

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **A "GatewaySubnet"** -- ARM requires the gateway's subnet to be named EXACTLY `GatewaySubnet` (/27 or larger recommended, no other workloads, no NSG). Create it as an AzureSubnet and reference its id.
- **A Standard static public IP** per ip configuration for VPN gateways (ExpressRoute gateways must NOT carry one -- Azure manages their addressing). A gateway binds its address exclusively.
- **Time and cost awareness** -- gateways take 25-45 minutes to create, 10-20 to delete, and bill hourly per SKU from the moment they provision (VpnGw1 ≈ $0.19/hour).

## Deploy

### Console

Open the deployment store, find **Azure Virtual Network Gateway**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Site-to-Site VPN** preset in the [Presets](#presets) tab to pre-populate a route-based VpnGw1 gateway.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetworkGateway
metadata:
  name: hub-vpn-gateway
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    value: "network-rg"
  name: hub-vpn-gateway
  sku: VPN_GW_1
  ipConfigurations:
    - subnetId:
        value: "/subscriptions/.../virtualNetworks/hub-vnet/subnets/GatewaySubnet"
      publicIpAddressId:
        value: "/subscriptions/.../publicIPAddresses/vpn-gw-pip"
  bgpEnabled: true
  bgpSettings:
    asn: 65515
```

```shell
planton apply -f azure-virtual-network-gateway.yaml
```

This creates a route-based VpnGw1 VPN gateway with BGP enabled. A Stack Job tracks the 25-45 minute provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the gateway to its subnet and address -- and wire each connection to the gateway:

```yaml
# On the AzureVirtualNetworkGateway:
spec:
  ipConfigurations:
    - subnetId:
        valueFrom:
          kind: AzureSubnet
          name: gateway-subnet
          fieldPath: status.outputs.subnet_id
      publicIpAddressId:
        valueFrom:
          kind: AzurePublicIp
          name: vpn-gw-pip
          fieldPath: status.outputs.public_ip_id

# On each AzureVirtualNetworkGatewayConnection that terminates here:
spec:
  virtualNetworkGatewayId:
    valueFrom:
      kind: AzureVirtualNetworkGateway
      name: hub-vpn-gateway
      fieldPath: status.outputs.virtual_network_gateway_id
```

The InfraPipeline resolves the dependency graph, deploys the subnet and address first, then the gateway, then the connections that attach to it.

## Key Configuration

These are the most important decisions when configuring a gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Type** -- `type` unset deploys VPN (IPsec tunnels, the site-to-site shape); EXPRESS_ROUTE terminates a circuit's private peering instead. Fixed at creation, and it flips the public-IP contract: VPN requires one per ip configuration, ExpressRoute forbids them.

**SKU and generation** -- `sku` is required. VPN_GW_1 is the production entry point (~650 Mbps aggregate); the `_AZ` variants add zone redundancy; GENERATION2 doubles throughput ceilings but starts at VPN_GW_2. BASIC is dev/test only and cannot resize in place. The spec enforces every type/vpn-type/generation/SKU pairing upfront -- a wrong combination fails validation in seconds, not deployment at minute 30.

**Active-active** -- `activeActive` runs the gateway as a two-instance pair (two ip configurations, each with its own public IP): no failover gap, at double the tunnel capacity cost.

**BGP** -- `bgpEnabled` + `bgpSettings` (ASN, APIPA peering addresses) enable dynamic route exchange. Azure's default ASN is 65515; 65515-65520 are Azure-reserved.

**NAT rules** -- `natRules` translate overlapping address space per tunnel; connections opt in via their `egressNatRuleIds`/`ingressNatRuleIds` using the ids this gateway publishes in `nat_rule_ids`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureSubnet** (the GatewaySubnet) | `ipConfigurations[].subnetId` | `status.outputs.subnet_id` |
| **AzurePublicIp** (per configuration) | `ipConfigurations[].publicIpAddressId` | `status.outputs.public_ip_id` |
| **AzureLocalNetworkGateway** (forced tunneling) | `defaultLocalNetworkGatewayId` | `status.outputs.local_network_gateway_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `virtual_network_gateway_id` | Azure Resource Manager ID of the gateway | AzureVirtualNetworkGatewayConnection's `virtualNetworkGatewayId` (and `peerVirtualNetworkGatewayId` for VNet-to-VNet) |
| `virtual_network_gateway_name` | Name of the gateway | Diagnostics and operational tooling |
| `nat_rule_ids` | ARM ids of the gateway's NAT rules, keyed by rule name | Connections' `egressNatRuleIds`/`ingressNatRuleIds` (supplied as literals) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Site-to-site VPN gateway** -- A route-based VpnGw1 gateway with BGP: the standard datacenter-to-Azure anchor. Start from the **Site-to-Site VPN** preset.

**Active-active zone-redundant gateway** -- Two instances on VpnGw2AZ with per-instance addresses and APIPA BGP: the high-availability posture. Start from the **Active-Active Zone-Redundant** preset.

**Point-to-site with Entra ID** -- A gateway whose VPN clients authenticate with Entra ID over OpenVPN. Start from the **Point-to-Site Entra ID** preset.

## Works With

- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- provides the resource group the gateway is created in
- [**Azure Subnet**](/cloud-catalog/azure-subnet) -- the dedicated "GatewaySubnet" the gateway lives in
- [**Azure Public IP**](/cloud-catalog/azure-public-ip) -- the addresses tunnels terminate on
- [**Azure Local Network Gateway**](/cloud-catalog/azure-local-network-gateway) -- describes each on-premises site
- [**Azure Virtual Network Gateway Connection**](/cloud-catalog/azure-virtual-network-gateway-connection) -- the tunnels that terminate on this gateway

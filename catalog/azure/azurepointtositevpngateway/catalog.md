# Azure Point-to-Site VPN Gateway

Deploys a Point-to-Site VPN Gateway -- the managed receiver inside a Virtual WAN hub that individual devices dial into from anywhere. HOW users authenticate lives on the VPN Server Configuration the gateway references; WHAT addresses connected clients get comes from the gateway's connection configurations. The gateway bills from creation and creates in 30-45 minutes. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Point-to-Site VPN Gateway** -- the managed instance pair in the hub, with its scale units, client address pools, and optional per-pool hub routing

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Virtual Hub** the gateway deploys into (one point-to-site gateway per hub -- a slot separate from the hub's site-to-site VPN gateway).
- **An Azure VPN Server Configuration** defining how users authenticate (both references are fixed at creation).
- **A client address pool** that overlaps nothing the hub reaches.

## Deploy

### Console

Open the deployment store, find **Azure Point-to-Site VPN Gateway**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Remote Users** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePointToSiteVpnGateway
metadata:
  name: remote-users-gw
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: network-rg
      fieldPath: status.outputs.resource_group_name
  name: remote-users-gw
  virtualHubId:
    valueFrom:
      kind: AzureVirtualHub
      name: hub-eastus
      fieldPath: status.outputs.virtual_hub_id
  vpnServerConfigurationId:
    valueFrom:
      kind: AzureVpnServerConfiguration
      name: remote-workforce
      fieldPath: status.outputs.vpn_server_configuration_id
  connectionConfigurations:
    - name: default-clients
      addressPrefixes:
        - "172.16.201.0/24"
```

```shell
planton apply -f azure-point-to-site-vpn-gateway.yaml
```

The gateway creates in 30-45 minutes and bills from creation.

### InfraChart

In a remote-access chart the order is: WAN → hub → server configuration → **point-to-site gateway**, each wiring to the previous by reference.

## Key Configuration

These are the most important decisions when configuring the gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The two references** -- `virtualHubId` (where the gateway lives) and `vpnServerConfigurationId` (how users authenticate) are both fixed at creation; changing either replaces the gateway. Changes WITHIN the referenced configuration apply in place.

**The address pool** -- size each `connectionConfigurations` pool for expected concurrent connections and overlap nothing the hub reaches. Most gateways carry one pool; multiple pools require OpenVPN on the server configuration and map to its policy groups.

**Scale units** -- 500 concurrent connections each. `scaleUnit` updates in place; unset applies 1 (the smallest gateway).

**Internet security** -- per pool: on, the hub advertises 0.0.0.0/0 so client internet traffic rides the tunnel (pair with a hub firewall via routing intent); off (the default), clients keep local internet egress.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureVirtualHub** | `virtualHubId` | `status.outputs.virtual_hub_id` |
| **AzureVpnServerConfiguration** | `vpnServerConfigurationId` | `status.outputs.vpn_server_configuration_id` |
| **AzureVirtualHub** (route block) | `route.associatedRouteTableId` | `status.outputs.default_route_table_id` or `status.outputs.route_table_ids.<table-name>` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `point_to_site_vpn_gateway_id` | ARM ID of the gateway | Operational tooling, diagnostics |
| `point_to_site_vpn_gateway_name` | Name of the gateway | Operational tooling |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard remote users** -- one pool, split tunneling (clients keep local internet). Start from the **Standard Remote Users** preset.

**Forced-tunnel clients** -- internet security on, all client traffic into the hub. Start from the **Forced-Tunnel Clients** preset.

## Works With

- [**Azure VPN Server Configuration**](/cloud-catalog/azure-vpn-server-configuration) -- the authentication policy the gateway attaches
- [**Azure Virtual Hub**](/cloud-catalog/azure-virtual-hub) -- where the gateway lives
- [**Azure Virtual WAN**](/cloud-catalog/azure-virtual-wan) -- the managed network umbrella
- [**Azure VPN Gateway**](/cloud-catalog/azure-vpn-gateway) -- the hub's site-to-site sibling (branches, not people)

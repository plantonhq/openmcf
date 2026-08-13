# Azure VPN Site

Deploys a VPN Site -- the Virtual WAN address-book entry for one branch location: its internet links (each with a public endpoint and optional BGP speaker), the address space reachable behind it, and the device that terminates the tunnels. The site is free and deploys nothing at the branch; a VPN Gateway Connection points at it to build the actual tunnels. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPN Site** -- the ARM description of the branch, including its links (ARM assigns each link an ID the `link_ids` output republishes by name)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Virtual WAN** the site belongs to (sites are WAN-scoped so any hub's VPN gateway in the WAN can connect to them).
- **The branch's real details**: each link's public IP or FQDN, and either the reachable prefixes (`addressCidrs`) or per-link BGP.

## Deploy

### Console

Open the deployment store, find **Azure VPN Site**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Single-Link Branch** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVpnSite
metadata:
  name: branch-london
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: network-rg
      fieldPath: status.outputs.resource_group_name
  name: branch-london
  virtualWanId:
    valueFrom:
      kind: AzureVirtualWan
      name: corp-wan
      fieldPath: status.outputs.virtual_wan_id
  addressCidrs:
    - "192.168.10.0/24"
  links:
    - name: primary-isp
      ipAddress: "203.0.113.10"
      speedInMbps: 200
```

```shell
planton apply -f azure-vpn-site.yaml
```

The site provisions in seconds and is free.

### InfraChart

In a branch-connectivity chart the order is: WAN → hub → VPN gateway, plus one **site** per branch → one connection per site, each wiring to the previous by reference.

## Key Configuration

These are the most important decisions when configuring a site. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Links** -- the connectable unit. Each link needs a public endpoint (IP or FQDN); its ARM ID surfaces in `link_ids` keyed by the link's name, which is exactly what a connection's `vpnLinks` reference. Two links model a dual-ISP branch.

**Routing source** -- static `addressCidrs` (Azure routes those prefixes into the tunnels) or per-link `bgp` (the branch advertises its prefixes), or both. A site with neither routes nothing.

**Device metadata** -- `deviceVendor`/`deviceModel` are informational (portal display, SD-WAN partners); they change no behavior.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `resourceGroup` | `status.outputs.resource_group_name` |
| **AzureVirtualWan** | `virtualWanId` | `status.outputs.virtual_wan_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpn_site_id` | ARM ID of the site | A connection's `remoteVpnSiteId` |
| `vpn_site_name` | Name of the site | Operational tooling |
| `link_ids` | ARM ID of each link, keyed by link name | A connection link's `vpnSiteLinkId` (`status.outputs.link_ids.<link-name>`) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single-link branch** -- one ISP, static prefixes. Start from the **Single-Link Branch** preset.

**Dual-link BGP branch** -- two ISPs with per-link BGP for active-active connectivity. Start from the **Dual-Link BGP Branch** preset.

## Works With

- [**Azure Virtual WAN**](/cloud-catalog/azure-virtual-wan) -- the WAN the site belongs to
- [**Azure VPN Gateway**](/cloud-catalog/azure-vpn-gateway) -- the hub gateway branches connect to
- [**Azure VPN Gateway Connection**](/cloud-catalog/azure-vpn-gateway-connection) -- the tunnels that point at this site

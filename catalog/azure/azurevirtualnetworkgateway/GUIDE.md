# Azure Virtual Network Gateway -- Operational Guide

Judgment that does not fit in field references: how to choose shapes,
what the slow and expensive edges are, and where deployments go wrong.

## The three-resource shape

A working site-to-site tunnel is always three resources, and keeping
them separate is what makes the design composable:

1. **This gateway** -- the Azure-side appliance. One per VNet (of VPN
   type), expensive and slow to build. Treat it as long-lived hub
   infrastructure.
2. **AzureLocalNetworkGateway** -- the DESCRIPTION of one on-premises
   site (device endpoint + reachable prefixes). Free, instant, one per
   site.
3. **AzureVirtualNetworkGatewayConnection** -- the tunnel joining them.
   Cheap and fast; one per site (or per redundant device pair).

Adding a second branch office never touches the gateway: add one local
network gateway and one connection.

## Time and money are the design constraints

- Creation runs **25-45 minutes**; deletion 10-20. Anything ForceNew
  (type, vpn_type, generation, edge zone, private-IP enablement, any ip
  configuration -- including its subnet or public IP) costs you that
  full cycle plus tunnel downtime for every connected site. Get the SKU
  family and generation right the first time; resizing WITHIN a family
  (VpnGw1AZ -> VpnGw2AZ) is in-place, crossing families (Basic ->
  anything, non-AZ -> AZ for ExpressRoute) is not.
- The gateway bills hourly from provisioning -- the SKU sets the
  rate, and Azure reduced AZ pricing as part of the SKU consolidation
  -- whether or not any tunnel is connected. Do not leave test
  gateways running overnight. Verified per-SKU figures live in the
  generated estimate at
  `catalog/_pricing/estimates/azurevirtualnetworkgateway.yaml`.

## Choosing the SKU

- **Only the AZ tiers accept new VPN gateway creates.** Azure retired
  the non-AZ VpnGw1-5 SKUs (creates rejected with
  NonAzSkusNotAllowedForVPNGateway since November 2025) and the legacy
  Standard/HighPerformance VPN tiers before that; existing non-AZ
  gateways are being migrated to their AZ twins. AZ SKUs deploy in
  EVERY region: zone-redundant where the region has availability
  zones, plain regional where it does not -- there is no region where
  "AZ" is the wrong choice.
- **AZ SKUs demand ZONED public IPs.** ARM rejects an AZ-SKU gateway
  whose Standard public IP carries no zones
  (VmssVpnGatewayPublicIpsMustHaveZonesConfigured) -- give the
  AzurePublicIp zones `["1","2","3"]` unless you are deliberately
  pinning a single zone. The rejection happens at deploy, not at
  validation: the zones live on the referenced address.
- **VPN_GW_1_AZ** is the right production default: 650 Mbps aggregate,
  30 tunnels, BGP. Move up when aggregate throughput or tunnel count
  says so.
- **BASIC is a trap** for anything durable: no zone redundancy, no
  IKEv2 point-to-site, no BGP, and NO in-place upgrade path -- leaving
  it means rebuilding the gateway and every tunnel. Use it only for
  short-lived experiments (it is also the only SKU policy-based VPN
  supports). It is not retiring, but new BASIC gateways must use a
  Standard public IP.
- **GENERATION2** doubles throughput ceilings but starts at
  VPN_GW_2_AZ; Azure picks the default generation when unset -- pin it
  explicitly if you ever plan to resize, because generation is
  ForceNew.

## The GatewaySubnet contract

ARM rejects any subnet not named exactly `GatewaySubnet`. Size it /27
or larger (Microsoft's guidance; /24 leaves room for coexisting
ExpressRoute + VPN gateways). Never attach an NSG or a route table
sending 0.0.0.0/0 away from Azure infrastructure -- both break the
gateway's control plane in ways that surface as mysterious tunnel
failures, not clear errors.

## Provisioned is not connected

The gateway (and its connections) reach ARM `Succeeded` when Azure
provisions them -- NOT when tunnels establish. A tunnel stuck in
`Connecting` means the far side disagrees (shared key, IKE version,
IPsec proposal, or the on-premises device is not configured yet).
Check `az network vpn-connection show` for the connection status and
the gateway's BGP peer status before suspecting the deployment.

## BGP judgment

- Leave the ASN at Azure's 65515 default unless your on-premises AS
  design says otherwise; 65515-65520 are Azure-reserved and the spec's
  comment documents the safe private ranges.
- APIPA peering addresses (169.254.21.0-169.254.22.255) exist for peers
  that demand link-local BGP endpoints -- AWS site-to-site VPN being
  the canonical case. Configure them per ip configuration, and name the
  ip configuration explicitly when the gateway has more than one.

## Active-active honestly

Active-active removes the failover gap (single-instance gateways take
tunnels down for minutes during Azure maintenance) at the cost of two
public IPs, two tunnels per site, and on-premises devices that must
handle both. If the on-premises side cannot run two tunnels, an
active-passive single instance is the honest configuration -- do not
enable active-active for a device that will only ever talk to one
endpoint.

## NAT rules

Model overlapping-address translation on the gateway (`natRules`), then
opt each connection in via `egressNatRuleIds`/`ingressNatRuleIds` with
the ids from this gateway's `nat_rule_ids` output. Rules that no
connection references do nothing -- the opt-in is per tunnel by design.

## Engine note

ER_GW_SCALE (the autoscaling ExpressRoute SKU) deploys via the
Terraform engine only -- the Pulumi engine cannot express the autoscale
bounds and fails loudly rather than dropping them. Everything else
deploys identically on both engines.

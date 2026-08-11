# Active-Active Zone-Redundant Gateway

This preset creates a Generation2 VpnGw2AZ gateway running as an ACTIVE-ACTIVE pair: two gateway instances, each with its own public IP and APIPA BGP endpoint, both terminating tunnels simultaneously. Azure maintenance on one instance no longer drops connectivity -- the high-availability posture for tunnels that carry production traffic.

## When to Use

- Production site-to-site connectivity where a minutes-long failover gap is unacceptable
- Peers that require link-local BGP endpoints (AWS site-to-site VPN is the canonical case -- each AWS tunnel pairs with one APIPA address)
- Regions with availability zones (the `_AZ` SKU spreads instances across them)

## Key Configuration Choices

- **`activeActive: true` + two ip configurations** -- each instance needs its own configuration and exclusive public IP (spec-enforced)
- **`VPN_GW_2_AZ` + `GENERATION2`** -- Generation2 doubles throughput ceilings and starts at VpnGw2AZ (only the AZ tiers accept new creates). Generation is FIXED at creation
- **Named configurations + per-instance APIPA addresses** -- with multiple configurations, every `peeringAddresses` entry must name its `ipConfigurationName`; Azure public regions accept 169.254.21.0-169.254.22.255
- **On-premises devices must run BOTH tunnels** -- if the device can only manage one, use the single-instance Site-to-Site preset instead; active-active with a one-tunnel device buys nothing

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region with availability zones | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-gateway-subnet-resource-name>` | Planton metadata name of the `AzureSubnet` named "GatewaySubnet" | Your subnet resource |
| `<your-first-public-ip-resource-name>` / `<your-second-public-ip-resource-name>` | Two `AzurePublicIp` resources, one per instance | Your public IP resources |

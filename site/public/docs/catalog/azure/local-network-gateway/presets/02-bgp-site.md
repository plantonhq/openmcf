---
title: "BGP Site"
description: "This preset describes a site whose routes arrive dynamically over BGP instead of a static prefix list -- the posture for estates where sites multiply or prefixes churn. The description carries the..."
type: "preset"
rank: "02"
presetSlug: "02-bgp-site"
componentSlug: "local-network-gateway"
componentTitle: "Local Network Gateway"
provider: "azure"
icon: "package"
order: 2
---

# BGP Site

This preset describes a site whose routes arrive dynamically over BGP instead of a static prefix list -- the posture for estates where sites multiply or prefixes churn. The description carries the on-premises BGP speaker's identity; no `addressSpaces` are declared, so learned routes win unambiguously.

## When to Use

- Many sites, or sites whose reachable prefixes change (adding a subnet on-premises needs no Azure change)
- Devices on dynamic public addressing (the FQDN endpoint re-resolves)
- Active-active gateway designs where BGP steers per-tunnel routing

## Key Configuration Choices

- **BGP requires all three ends configured** -- the gateway (`bgpEnabled` + `bgpSettings`), this site description (`bgpSettings`), AND the connection (`bgpEnabled`). Any one missing and routes silently do not flow
- **`bgpPeeringAddress` is tunnel-interior** -- the device's tunnel interface (often a loopback or APIPA), NEVER its public address: the single most common BGP misconfiguration
- **ASN discipline** -- must differ from the Azure gateway's ASN (default 65515); 65515-65520 are Azure-reserved
- **No `addressSpaces`** -- deliberately empty so route provenance stays unambiguous; the spec permits both, but mixing complicates incident debugging

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (by convention, the connecting gateway's) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-device-fqdn>` | The VPN device's public DNS name | Your network team |
| `<your-onprem-asn>` | The on-premises BGP AS number | Your network team's AS design |
| `<your-device-tunnel-interface-ip>` | The BGP speaker's tunnel-interior address | The device's tunnel/loopback interface configuration |

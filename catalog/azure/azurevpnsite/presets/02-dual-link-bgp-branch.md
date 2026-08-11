# Dual-Link BGP Branch

This preset describes a resilient branch: two ISPs, each a separate link with its own BGP speaker, and no static prefix list -- the branch advertises its routes. A connection then builds one tunnel per link for active-active connectivity.

## When to Use

- Branches with two internet uplinks that must survive an ISP outage
- Branches whose internal prefixes change (BGP keeps Azure current without manifest edits)

## Key Configuration Choices

- **No `addressCidrs`** -- with BGP on every link, learned routes carry the routing; leaving the static list empty avoids a surprising union of both sources
- **Same ASN, different peering addresses** -- one branch device speaks for both links; each link peers from its own tunnel-inside address
- **Device metadata** -- vendor/model are informational but help support and SD-WAN partner automation identify the device

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the site object | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-virtual-wan-arm-id>` | ARM ID of the Virtual WAN | `AzureVirtualWan` status outputs (`virtual_wan_id`), or reference it with valueFrom |
| `<primary-isp-name>` / `<backup-isp-name>` | The ISPs behind each link | Informational -- the carrier names |

Replace the example endpoints, ASN, and peering addresses with the branch's real values (the branch ASN must avoid Azure-reserved 65515-65520).

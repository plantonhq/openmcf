---
title: "NAT for Overlapping Branches"
description: "This preset creates a gateway prepared for the classic acquisition problem: two branches that both use `192.168.10.0/24`. Each gets a static ingress NAT rule translating it to a distinct..."
type: "preset"
rank: "02"
presetSlug: "02-nat-overlapping-branches"
componentSlug: "vpn-gateway"
componentTitle: "VPN Gateway"
provider: "azure"
icon: "package"
order: 2
---

# NAT for Overlapping Branches

This preset creates a gateway prepared for the classic acquisition problem: two branches that both use `192.168.10.0/24`. Each gets a static ingress NAT rule translating it to a distinct `100.64.x.0/24` space, so Azure (and every other branch) sees them as different networks. Rules are inert until a connection's tunnel opts in.

## When to Use

- Branch estates with colliding private address spaces (mergers, franchises, default router configs)
- Ahead of onboarding branches you know will collide -- rules are free to declare

## Key Configuration Choices

- **INGRESS_SNAT + STATIC_NAT** -- translate the branch-side source one-to-one as it enters Azure; static NAT keeps ports intact and translation predictable
- **One rule per colliding branch** -- each branch maps to its own external space; the rule's NAME is the key its connection references via `status.outputs.nat_rule_ids.<rule-name>`
- **BGP route translation ON** -- with BGP branches, advertisements must carry the post-NAT prefixes or the overlap re-leaks
- **Scale unit 2** -- NAT estates are usually bigger estates; size to measured aggregate demand

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the gateway | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-virtual-hub-arm-id>` | ARM ID of the hub | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |

Replace the example internal/external mappings with the real colliding and translated spaces.

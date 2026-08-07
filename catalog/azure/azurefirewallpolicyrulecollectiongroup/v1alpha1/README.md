# AzureFirewallPolicyRuleCollectionGroup

## Overview

`AzureFirewallPolicyRuleCollectionGroup` provisions a rule collection
group inside an Azure Firewall Policy: an ordered document of application,
network, and DNAT rule collections that the policy's firewalls enforce. A
policy carries many groups -- typically one per team or per application --
each deployed and updated independently.

## Why a First-Class Resource?

- **Many per policy, independent lifecycles** -- the platform team's
  baseline rules and an application team's rules deploy on their own
  schedules without touching each other (or the policy)
- **The deployment unit Azure itself defines** -- the group is a real
  nested ARM resource; rules fold inside it because nothing references an
  individual rule
- **Serialized safely** -- concurrent group deployments against one policy
  queue (the provider locks the parent policy) rather than conflict

## Processing Order (worth internalizing)

1. Groups run by **group priority**; collections by **collection
   priority** within a type. Lower numbers first (100 is highest).
2. Across collection TYPES, Azure always processes **DNAT → network →
   application**, regardless of priorities.
3. A matching DNAT rule **implicitly allows** the translated flow -- no
   companion network rule needed.

## Key Features

- **Application rules** -- allow/deny by HTTP/HTTPS/MSSQL destination:
  FQDNs, URLs (Premium, needs TLS termination), Azure-curated FQDN tags,
  web categories (Premium), plus header injection
- **Network rules** -- classic L3/L4 by protocol/source/destination/port,
  with FQDN destinations (needs the policy's DNS proxy)
- **DNAT rules** -- publish internal services through the firewall's
  public IP (TCP/UDP only; exactly one translation target, address XOR
  FQDN -- validated at authoring time)
- **IP Group references** -- every source/destination address list has a
  reusable `AzureIpGroup` counterpart field

## Spec Fields (top level)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `firewall_policy_id` | StringValueOrRef | Yes | Parent policy (defaults to AzureFirewallPolicy reference); ForceNew |
| `name` | string | Yes | Group name, 2-80 chars, unique within the policy; ForceNew |
| `priority` | int | Yes | 100-65000, lower first |
| `application_rule_collections` | list | No | L7 collections (name, priority, ALLOW/DENY, rules) |
| `network_rule_collections` | list | No | L3/L4 collections |
| `nat_rule_collections` | list | No | DNAT collections (action is always Azure's "Dnat" -- not modeled) |

## Outputs

| Output | Description |
|--------|-------------|
| `rule_collection_group_id` | Full nested ARM ID of the group |
| `rule_collection_group_name` | The group's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFirewallPolicyRuleCollectionGroup
metadata:
  name: platform-baseline
  org: mycompany
  env: production
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  name: platform-baseline
  priority: 500
  networkRuleCollections:
    - name: core-egress
      priority: 200
      action: ALLOW
      rules:
        - name: allow-dns
          protocols: [UDP]
          sourceAddresses: ["10.0.0.0/16"]
          destinationAddresses: ["8.8.8.8"]
          destinationPorts: ["53"]
```

## Lifecycle Notes

- `name` and `firewall_policy_id` are **fixed at creation**; priority and
  every collection update in place
- **Priorities**: leave gaps (100, 200, 300...) so future collections
  slot in without renumbering
- **Premium destinations** (`destination_urls`, `web_categories`,
  `terminate_tls`) require a PREMIUM policy -- ARM rejects them on
  Standard at apply time
- **FQDN network rules** require the policy's DNS proxy so the firewall
  and clients resolve identically
- An **empty group** (no collections) is legal -- a placeholder that
  reserves a priority slot

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

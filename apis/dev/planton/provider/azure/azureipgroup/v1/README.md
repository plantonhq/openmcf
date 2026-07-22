# AzureIpGroup

## Overview

`AzureIpGroup` provisions an Azure IP Group: a named, reusable set of IP
addresses and CIDR ranges that Azure Firewall and Firewall Policy rules
reference by ARM id instead of repeating literal address lists. A rule can
say "allow the branch offices to reach the datacenter ranges" rather than
enumerating twenty CIDRs inline in every rule -- update the group once and
every rule that references it follows.

## Why a First-Class Resource?

An IP Group is real infrastructure with its own lifecycle:

- **Created once, referenced by many** -- firewall policy application,
  network, and NAT rules list source/destination IP Groups, and
  intrusion-detection traffic bypasses reference them the same way;
  nothing lives inside the group but addresses
- **Independent lifecycle** -- the address set is curated on its own
  schedule (a new branch office CIDR lands in one place) and outlives the
  individual rules referencing it
- **Intent-carrying composition seam** -- `ip_group_id` is the join key
  every rule references

The inversion -- consumption declared from the rule's side, not the
group's -- is what makes the IP Group composable. Updating the group's
addresses immediately retargets every referencing rule without touching
the rules themselves.

## Key Features

- **Address sets by intent** -- name the group after what the addresses
  mean ("branch-offices", "on-prem-datacenter", "blocked-scanners") so
  firewall rules read as policy statements
- **In-place updates** -- adding or removing an address updates every
  referencing rule at once; the group itself never recreates for an
  address change
- **Mixed entries** -- single addresses ("203.0.113.7") and CIDR blocks
  ("10.10.0.0/16") coexist in one group (up to 5,000 entries)
- **Composable** -- the resource group is referenced by name, defaulting
  to an `AzureResourceGroup`'s output in composed environments

## Spec Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | -- | Azure region (regional resource; referable from policies in any region) |
| `resource_group` | StringValueOrRef | Yes | -- | Resource group name (defaults to an AzureResourceGroup reference) |
| `name` | string | Yes | -- | Group name, unique within the resource group (1-80 chars); fixed at creation |
| `cidrs` | list(string) | No | `[]` | IP addresses and CIDR ranges; updates in place |
| `tags` | map | No | -- | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `ip_group_id` | Full ARM ID -- the join key firewall policy rules and IDPS bypasses reference |
| `ip_group_name` | The group's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureIpGroup
metadata:
  name: branch-offices
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: branch-offices
  cidrs:
    - "198.51.100.0/24"
    - "203.0.113.0/24"
```

Target it from a firewall policy rule collection group:

```yaml
spec:
  networkRuleCollections:
    - name: allow-branches-to-dc
      priority: 200
      action: ALLOW
      rules:
        - name: branches-to-dc-any
          protocols: [ANY]
          sourceIpGroups:
            - valueFrom:
                name: branch-offices
          destinationIpGroups:
            - valueFrom:
                name: on-prem-datacenter
          destinationPorts: ["*"]
```

## Lifecycle Notes

- `name` and `region` are **fixed at creation**; changing either replaces
  the group, and every rule that referenced it must be re-pointed
- `cidrs` and `tags` update **in place** -- address changes propagate to
  every referencing rule immediately
- An IP Group can be referenced from at most **100 firewall policies**
  and holds at most **5,000 entries** (Azure Firewall limits)
- An **empty group is legal** -- a placeholder anchor rules can reference
  before the address plan is final

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).

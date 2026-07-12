---
title: "Workload Tier Group"
description: "This preset creates a single application security group named after a workload role (\"web-tier\"). It is the building block of address-independent micro-segmentation: instead of writing NSG rules..."
type: "preset"
rank: "01"
presetSlug: "01-workload-tier"
componentSlug: "application-security-group"
componentTitle: "Application Security Group"
provider: "azure"
icon: "package"
order: 1
---

# Workload Tier Group

This preset creates a single application security group named after a
workload role ("web-tier"). It is the building block of address-independent
micro-segmentation: instead of writing NSG rules against CIDRs, you group
each tier's network interfaces into an ASG and write rules that reference
the groups by name.

Create one group per tier (web, app, data) and the security policy reads as
intent -- "web reaches app on 8080", "app reaches data on 5432" -- and
survives every scale event without a firewall rewrite.

## When to Use

- Multi-tier applications where each tier should reach only its neighbors
- Auto-scaling workloads whose instance addresses change constantly
- Any design where NSG rules should express roles, not IP ranges

## Key Configuration Choices

- **`name`** -- name it after the tier it represents; NSG rules reference
  this value, so a good name makes the policy self-documenting
- **`region`** -- must match the region of the network interfaces that will
  join the group; an ASG cannot span regions

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the group in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Join a network interface to this group:

```yaml
spec:
  applicationSecurityGroupIds:
    - valueFrom:
        name: my-workload-tier
```

Target it in an NSG rule (as source or destination):

```yaml
spec:
  securityRules:
    - name: allow-web-to-app
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      priority: 100
      destinationPortRange: "8080"
      sourceApplicationSecurityGroupIds:
        - valueFrom:
            name: my-workload-tier
```

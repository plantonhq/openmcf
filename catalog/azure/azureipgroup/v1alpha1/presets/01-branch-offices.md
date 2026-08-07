# Branch-Office Address Set

This preset creates an IP Group carrying the CIDR ranges of an
organization's branch offices. It is the building block of maintainable
firewall policy: instead of enumerating branch CIDRs inline in every rule
that means "the branches", rules reference this group by id -- and when a
new office opens, one address lands here and every rule follows.

Create one group per address-set intent (branch offices, on-prem
datacenter ranges, partner networks, scanner deny-lists) and firewall
rules read as policy statements.

## When to Use

- Hub-spoke egress policies where many rules share the same source or
  destination address sets
- Organizations whose network footprint changes (new offices, new VPN
  ranges) more often than their security policy does
- Any design where a firewall rule should express intent, not a CIDR
  enumeration

## Key Configuration Choices

- **`name`** -- name it after what the addresses mean; firewall rules
  reference this group, so a good name makes the policy self-documenting
- **`cidrs`** -- single addresses and CIDR blocks both work; up to 5,000
  entries per group, updated in place

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the group in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Reference the group from a firewall policy network rule:

```yaml
spec:
  networkRuleCollections:
    - name: allow-branches-to-dc
      priority: 200
      action: ALLOW
      rules:
        - name: branches-to-dc
          protocols: [ANY]
          sourceIpGroups:
            - valueFrom:
                name: branch-offices
          destinationAddresses: ["10.0.0.0/8"]
          destinationPorts: ["*"]
```

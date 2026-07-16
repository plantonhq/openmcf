# On-Premises Datacenter Ranges

This preset creates an IP Group carrying the private ranges of an
on-premises datacenter -- the destination half of the classic hybrid
egress policy ("branches and Azure workloads may reach the datacenter;
nothing else may"). Pairing it with a branch-offices group lets a single
firewall rule express the whole site-to-site policy through two
referenced groups.

## When to Use

- Hybrid connectivity (ExpressRoute/VPN) where Azure workloads reach
  on-premises systems through an Azure Firewall
- Migration programs where the on-prem footprint shrinks over time --
  retiring a range is one edit here, not a rule audit
- Separating address curation (a network team's job) from rule authoring
  (a security team's job): each team owns its own resource

## Key Configuration Choices

- **`cidrs`** -- carry every on-prem range the policy should treat as "the
  datacenter"; rules referencing the group follow automatically as ranges
  are added or retired
- **`tags`** -- tag the group with its scope so governance tooling can
  distinguish hybrid address sets from cloud-only ones

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the group in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Use the group as a rule destination:

```yaml
spec:
  networkRuleCollections:
    - name: allow-azure-to-onprem
      priority: 300
      action: ALLOW
      rules:
        - name: workloads-to-dc
          protocols: [TCP]
          sourceAddresses: ["10.0.0.0/16"]
          destinationIpGroups:
            - valueFrom:
                name: on-prem-datacenter
          destinationPorts: ["443", "1433"]
```

# Governed Data-Tier Group

This preset creates a data-tier application security group carrying a full
governance tag set -- cost center, owning team, and data classification.
Tags are Azure's governance surface: Azure Policy enforces required tags,
and Cost Management groups spend by them, so a group that participates in a
regulated data path should carry its classification from creation.

## When to Use

- Data-tier or regulated workloads that must be tagged for compliance
- Organizations enforcing tag policies via Azure Policy
- Any group whose spend or ownership must be attributable

## Key Configuration Choices

- **`data-classification` tag** -- signals the sensitivity of the workloads
  this group fronts; pair it with tight NSG rules that only allow the app
  tier to reach it
- **`owner` tag** -- the team accountable for the segment, so security
  reviews have a clear escalation path

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the group in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
| `<team-email>` | The owning team's contact | Your org directory |

## Downstream Wiring

Restrict the data tier to only the app tier in an NSG rule:

```yaml
spec:
  securityRules:
    - name: allow-app-to-data
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      priority: 100
      destinationPortRange: "5432"
      sourceApplicationSecurityGroupIds:
        - valueFrom:
            name: app-tier
      destinationApplicationSecurityGroupIds:
        - valueFrom:
            name: my-governed-asg
```

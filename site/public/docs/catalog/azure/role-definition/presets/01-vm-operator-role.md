---
title: "Explicit-Actions Operator Role"
description: "This preset defines the most common custom-role shape: an explicit list of control-plane actions capturing exactly what an operations team may do -- here, observe and restart existing VMs without any..."
type: "preset"
rank: "01"
presetSlug: "01-vm-operator-role"
componentSlug: "role-definition"
componentTitle: "Role Definition"
provider: "azure"
icon: "package"
order: 1
---

# Explicit-Actions Operator Role

This preset defines the most common custom-role shape: an explicit list of
control-plane actions capturing exactly what an operations team may do -- here,
observe and restart existing VMs without any right to create, modify, or
delete them. No built-in role expresses this; Virtual Machine Contributor can
delete VMs, Reader cannot restart them.

Explicit action lists are the least-privilege grain: unlike wildcard grants,
they do not silently grow as ARM adds new operations. Start from this shape
for any role whose purpose is "do these specific things and nothing else".

Grant the role with an `AzureRoleAssignment` binding the definition's output:

```yaml
spec:
  scope:
    valueFrom:
      name: platform-rg
  roleDefinitionId: <the definition's status.outputs.role_definition_id>
  principalId:
    valueFrom:
      name: ops-identity
```

## When to Use

- Operations teams that manage workload lifecycle but must not alter infrastructure
- Break-glass roles whose exact powers must be auditable at a glance
- Any role reviewed like code, where each permitted operation is explicit

## Key Configuration Choices

- **Subscription scope** -- the role is visible and assignable subscription-wide;
  narrow `scope` to a resource group (by reference) to keep a team-specific
  role invisible outside its environment
- **Action granularity** -- find operations per provider with
  `az provider operation show --namespace Microsoft.Compute`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<org-prefix>` | Your organization's role-name prefix (names are tenant-unique) | Your naming convention |
| `<subscription-arm-id>` | `/subscriptions/{subscription-id}` | `az account show --query id` (prepend `/subscriptions/`) |

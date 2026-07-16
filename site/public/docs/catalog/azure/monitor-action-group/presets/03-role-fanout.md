---
title: "ARM Role Fan-Out"
description: "This preset notifies every user assigned an ARM role on the subscription -- here the built-in Owner role. RBAC membership becomes the notification list, so joiners and leavers are handled by access..."
type: "preset"
rank: "03"
presetSlug: "03-role-fanout"
componentSlug: "monitor-action-group"
componentTitle: "Monitor Action Group"
provider: "azure"
icon: "package"
order: 3
---

# ARM Role Fan-Out

This preset notifies every user assigned an ARM role on the subscription -- here the built-in Owner role. RBAC membership becomes the notification list, so joiners and leavers are handled by access management instead of an address list nobody maintains.

## When to Use

- Subscription-level alerts (quota exhaustion, Service Health, budget events) that must reach whoever owns the subscription
- Estates where team rosters change faster than alert configuration

## Key Configuration Choices

- **Built-in role by well-known GUID** -- Owner is `8e3af657-a8ff-443c-a75c-2fe8c4bcb635` (built-in role GUIDs are a vendor catalog: look them up, never infer). Contributor: `b24988ac-6180-42a0-ab88-20f7382dd24c`; Monitoring Reader: `43d0d8ad-25c7-4714-9337-8ba259a9fe05`
- **Custom roles compose by FK** -- reference an `AzureRoleDefinition`'s `role_definition_guid` output with valueFrom instead of a literal
- **Email-only channel** -- ARM-role receivers email the role members; pair with an on-call group for paging

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the group | `AzureResourceGroup` status outputs |

## Related Presets

- **01-oncall-team** -- Direct paging channels
- **02-automation-hooks** -- Machine receivers

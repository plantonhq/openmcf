---
title: "On-Call Team Notifications"
description: "This preset creates the human-notification action group: email, SMS, and Azure mobile app push for the on-call rotation. Metric alerts, query alerts, and Service Health alerts all reference this one..."
type: "preset"
rank: "01"
presetSlug: "01-oncall-team"
componentSlug: "monitor-action-group"
componentTitle: "Monitor Action Group"
provider: "azure"
icon: "package"
order: 1
---

# On-Call Team Notifications

This preset creates the human-notification action group: email, SMS, and Azure mobile app push for the on-call rotation. Metric alerts, query alerts, and Service Health alerts all reference this one group -- keep it stable and let the alert rules churn around it.

## When to Use

- The first action group of any environment -- every alert needs somewhere to fire
- Team-scoped routing (one group per team, referenced by that team's alert rules)

## Key Configuration Choices

- **Short name shows on phones** (`shortName: PltOnCall`, max 12 chars) -- it is the SMS sender identity; make it recognizable on a lock screen
- **Common alert schema on email** -- one consistent JSON payload across alert types; SMS and push are not payload-aware, so they carry no toggle
- **Group-level kill switch** -- set `enabled: false` during maintenance windows to silence everything without touching alert rules

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the group | `AzureResourceGroup` status outputs |
| `oncall@example.com` | The rotation's address | Your paging setup |
| `5555550100` / `"1"` | Phone number and country code | Your paging setup |

## Related Presets

- **02-automation-hooks** -- Machine receivers (webhooks, Functions, runbooks)
- **03-role-fanout** -- Notify everyone holding an ARM role, no address list to maintain

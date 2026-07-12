---
title: "Service-Health Incident Alert"
description: "This preset creates a subscription-scoped alert that fires when Azure reports a service-health incident or maintenance event. Service-health alerts are the ones you cannot get any other way -- they..."
type: "preset"
rank: "01"
presetSlug: "01-service-health"
componentSlug: "monitor-activity-log-alert"
componentTitle: "Monitor Activity Log Alert"
provider: "azure"
icon: "package"
order: 1
---

# Service-Health Incident Alert

This preset creates a subscription-scoped alert that fires when Azure
reports a service-health incident or maintenance event. Service-health
alerts are the ones you cannot get any other way -- they tell you when Azure
itself (not your code) is having a problem in a region or service you depend
on, so your on-call learns from your own alerting channel instead of the
status page.

## When to Use

- Every production subscription -- know about Azure incidents affecting your
  services before customers report them
- Coordinating maintenance windows against Azure's planned maintenance

## Key Configuration Choices

- **`category: SERVICE_HEALTH`** with a `serviceHealth` block -- narrow to
  `INCIDENT`/`MAINTENANCE`/`SECURITY` events, and optionally to specific
  `services` and `locations`
- **Subscription scope** -- service health is subscription-wide, so scope
  the alert at the subscription, not a resource group

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the alert in | The resource group's `status.outputs.resource_group_name` |
| `<subscription-id>` | The subscription to watch | `az account show --query id` |
| `<action-group>` | The AzureMonitorActionGroup to notify | Your action group's Planton resource name |

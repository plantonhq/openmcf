# AzureMonitorActionGroup

## Overview

`AzureMonitorActionGroup` provisions an Azure Monitor action group -- the
notification and automation hub every alert fires into. When a metric alert,
log-query alert, activity-log alert, or Service Health event triggers, Azure
notifies every receiver in the referenced group.

Action groups are a GLOBAL resource: they live in a resource group but not in
a region, so notifications keep flowing during regional outages -- exactly
when they matter most. One group is typically shared by many alert rules
(an on-call group per team), so treat the group as the stable routing node
and the alert rules as the volatile edge.

## Key Features

- **All eleven receiver families** -- email, SMS, voice, Azure mobile app
  push, webhooks (with Entra authentication for the keyless posture),
  Automation runbooks, Logic Apps, Azure Functions, ARM-role fan-out,
  Event Hubs, and ITSM connections
- **Common alert schema** -- the one consistent JSON payload, modeled on
  exactly the seven receiver types that carry it (SMS, voice, push, and ITSM
  are not payload-aware)
- **Composable receivers** -- the Function receiver references
  `AzureFunctionApp`, the ARM-role receiver references `AzureRoleDefinition`,
  the Event Hub receiver references `AzureEventHubNamespace`
- **Kill switch** -- `enabled: false` silences everything during maintenance
  windows without touching alert rules

## When to Use

- The first observability resource after the workspace -- every alert needs a
  routing target
- Team-scoped notification policies referenced by that team's alert rules
- Automation fan-out (webhooks into incident platforms, remediation functions)

## Spec Highlights

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `resource_group` | StringValueOrRef | Yes | Resource group (action groups are global -- no region) |
| `short_name` | string | Yes | 1-12 chars; the SMS/push sender identity |
| `enabled` | bool | No (true) | The group-level kill switch |
| `email_receivers[]` ... `itsm_receivers[]` | messages | No | Any number of receivers across the 11 families |
| `tags` | map | No | User tags |

## Outputs

| Output | Description |
|--------|-------------|
| `action_group_id` | ARM resource ID -- the FK seam `AzureMonitorMetricAlert` and `AzureMonitorScheduledQueryAlert` reference |
| `action_group_name` | The group name |

## Composition

```yaml
actionGroupId:
  valueFrom:
    kind: AzureMonitorActionGroup
    name: my-platform-oncall
    fieldPath: status.outputs.action_group_id
```

See `presets/` for on-call, automation, and role-fan-out starting points.

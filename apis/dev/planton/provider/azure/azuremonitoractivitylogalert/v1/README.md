# AzureMonitorActivityLogAlert

## Overview

`AzureMonitorActivityLogAlert` provisions an Azure Monitor Activity Log
Alert: a rule that fires an action (notify an action group, post a webhook)
when a matching entry appears in the Azure Activity Log. It is the only way
to alert on the events that never appear as metrics -- a resource someone
deleted, a region-wide service-health incident, a resource going
Unavailable, a policy denial, an Advisor recommendation.

## Why a First-Class Resource?

- **The only alerting path for control-plane and health events** -- metric
  alerts watch numbers; this watches the log of what happened to your
  infrastructure
- **Composes with action groups** -- the same `AzureMonitorActionGroup`
  used by metric and query alerts, so notification routing is uniform
- **Watches any scope** -- a subscription, a resource group, or specific
  resources

## Key Features

- **Seven categories** -- Administrative, Autoscale, Policy, Recommendation,
  ResourceHealth, Security, ServiceHealth; `category` selects the log slice
- **Rich narrowing** -- operation, caller, levels, resource
  providers/types/groups/ids, statuses/sub-statuses
- **Service-health alerts** -- incident/maintenance/security events by
  region and service
- **Resource-health alerts** -- current/previous state transitions and their
  reason
- **Advisor recommendation alerts** -- by category+impact or a specific type
- **Composable** -- scope defaults to a resource group, actions default to
  action groups

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `resource_group` | StringValueOrRef | Yes | Resource group (defaults to AzureResourceGroup) |
| `name` | string | Yes | Alert name, unique in the resource group; fixed at creation |
| `location` | enum | No | Alert-resource region (default GLOBAL) |
| `scopes` | repeated StringValueOrRef | Yes (≥1) | What the alert watches (default AzureResourceGroup) |
| `criteria` | message | Yes | The matching criteria (see below) |
| `actions` | repeated message | No | Action groups to notify + webhook properties |
| `description` | string | No | Human-readable description |
| `enabled` | bool | No | Active (default true) |
| `tags` | map | No | User tags, merged over Planton-derived tags |

### `criteria`

`category` (required) selects the Activity Log slice; the rest narrow within
it. Plural list fields (`levels`, `resource_providers`, `resource_types`,
`resource_groups`, `resource_ids`, `statuses`, `sub_statuses`) carry the
single case as a one-element list. `resource_health` and `service_health`
are category-specific blocks; `recommendation_category`/`_impact`/`_type`
apply to the Recommendation category. `caller`, `resource_health`, and
`service_health` are mutually exclusive.

## Outputs

| Output | Description |
|--------|-------------|
| `activity_log_alert_id` | Full ARM ID of the alert |
| `activity_log_alert_name` | The alert's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorActivityLogAlert
metadata:
  name: vm-delete-alert
spec:
  resourceGroup:
    valueFrom:
      name: platform-rg
  name: vm-delete-alert
  scopes:
    - valueFrom:
        name: platform-rg
  criteria:
    category: ADMINISTRATIVE
    operationName: Microsoft.Compute/virtualMachines/delete
    levels: [CRITICAL]
    statuses: [Succeeded]
  actions:
    - actionGroupId:
        valueFrom:
          name: ops-action-group
```

## Lifecycle Notes

- `name` is **fixed at creation**
- `location` must be one of GLOBAL / WEST_EUROPE / NORTH_EUROPE /
  EAST_US_2_EUAP (Azure only supports the alert resource in these); GLOBAL is
  correct for virtually every alert
- Scopes, criteria, actions, description, `enabled`, and `tags` update in
  place

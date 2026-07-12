# Azure Monitor Activity Log Alert

Creates an Azure Monitor Activity Log Alert — a rule that fires actions when a matching entry appears in the Azure Activity Log. It is the only way to alert on control-plane operations, service-health incidents, resource-health transitions, policy events, and Advisor recommendations — the events that never appear as metrics.

## What Gets Created

When you deploy an AzureMonitorActivityLogAlert resource, Planton provisions:

- **Activity Log Alert** — an `azurerm_monitor_activity_log_alert` watching the configured scopes, matching on the criteria, and firing the configured actions

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the alert in (an `AzureResourceGroup`)
- **An action group** to notify (an `AzureMonitorActionGroup`)
- **Monitoring write rights**: `Microsoft.Insights/activityLogAlerts/write`

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorActivityLogAlert
metadata:
  name: vm-delete-alert
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureMonitorActivityLogAlert.vm-delete-alert
spec:
  resourceGroup:
    value: platform-rg
  name: vm-delete-alert
  scopes:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg
  criteria:
    category: ADMINISTRATIVE
    operationName: Microsoft.Compute/virtualMachines/delete
    levels: [CRITICAL]
  actions:
    - actionGroupId:
        value: /subscriptions/.../Microsoft.Insights/actionGroups/ops
```

Deploy:

```shell
planton apply -f activity-log-alert.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `resourceGroup` | `StringValueOrRef` | Resource group. Defaults to an `AzureResourceGroup` reference. |
| `name` | `string` | Alert name, unique in the resource group. Fixed at creation. |
| `scopes` | `StringValueOrRef[]` | What the alert watches (subscription, RG, or resources). Defaults to `AzureResourceGroup`. |
| `criteria` | `object` | The matching criteria; `category` (required) selects the log slice. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `location` | `enum` | `GLOBAL` | Alert-resource region (GLOBAL / WEST_EUROPE / NORTH_EUROPE / EAST_US_2_EUAP). |
| `actions` | `object[]` | `[]` | Action groups to notify (`actionGroupId` + optional `webhookProperties`). |
| `description` | `string` | `""` | Human-readable description. |
| `enabled` | `bool` | `true` | Whether the alert is active. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins). |

### Criteria

`category` is one of `ADMINISTRATIVE`, `AUTOSCALE`, `POLICY`, `RECOMMENDATION`, `RESOURCE_HEALTH`, `SECURITY`, `SERVICE_HEALTH`. Narrow with `operationName`, `caller`, `levels`, `resourceProviders`, `resourceTypes`, `resourceGroups`, `resourceIds`, `statuses`, `subStatuses`; the `resourceHealth` and `serviceHealth` blocks and the `recommendation*` fields apply to their categories. `caller`, `resourceHealth`, and `serviceHealth` are mutually exclusive.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `activity_log_alert_id` | `string` | Full ARM ID of the alert |
| `activity_log_alert_name` | `string` | The alert's name as deployed |

## Related Components

- [AzureMonitorActionGroup](/docs/catalog/azure/monitor-action-group) — the notification target the alert fires
- [AzureMonitorMetricAlert](/docs/catalog/azure/monitor-metric-alert) — the metric-based counterpart
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group and a common scope

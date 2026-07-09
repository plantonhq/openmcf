---
title: "Monitor Scheduled Query Alert"
description: "Monitor Scheduled Query Alert deployment documentation"
icon: "package"
order: 100
componentName: "azuremonitorscheduledqueryalert"
---

# Azure Monitor Scheduled Query Alert

Deploys an Azure Monitor scheduled query alert -- a KQL query run on a schedule against a Log Analytics Workspace (or Application Insights resource), firing action groups when its condition holds. The alerting half of the logging pipeline.

## What Gets Created

When you deploy an AzureMonitorScheduledQueryAlert resource, Planton provisions:

- **Scheduled query rule** -- an `azurerm_monitor_scheduled_query_rules_alert_v2` in the queried resource's region, carrying the KQL criteria, evaluation cadence, noise controls, optional managed identity, and action-group wiring
- **Azure Tags** -- Planton-derived governance tags merged with your own (your values win)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureLogAnalyticsWorkspace** (or Application Insights resource) whose logs you are alerting on
- **An AzureMonitorActionGroup** to fire into

## Quick Start

Create a file `query-alert.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorScheduledQueryAlert
metadata:
  name: error-spike
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMonitorScheduledQueryAlert.error-spike
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: observability-rg
      fieldPath: status.outputs.resource_group_name
  alertName: app-error-spike
  scope:
    valueFrom:
      kind: AzureLogAnalyticsWorkspace
      name: platform-logs
      fieldPath: status.outputs.workspace_id
  severity: 2
  evaluationFrequency: PT5M
  windowDuration: PT10M
  criteria:
    - query: |
        AppExceptions
        | where TimeGenerated > ago(10m)
      timeAggregationMethod: COUNT
      operator: GREATER_THAN
      threshold: 10
  autoMitigationEnabled: true
  action:
    actionGroupIds:
      - valueFrom:
          kind: AzureMonitorActionGroup
          name: platform-oncall
          fieldPath: status.outputs.action_group_id
```

Deploy it:

```bash
planton pulumi up --manifest query-alert.yaml
```

## Common Configurations

- **Metric measurement**: a non-COUNT aggregation over a projected column (`timeAggregationMethod: AVERAGE` + `metricMeasureColumn`) for latency-style conditions
- **Per-service splits**: `dimensions` with `INCLUDE` + `"*"` -- each value alerts independently
- **Flap damping**: `failingPeriods` requiring N of the last M evaluations to breach
- **Absence of data**: `COUNT` + `LESS_THAN` (missing heartbeats); pair with `muteActionsAfterAlertDuration` instead of auto-mitigation -- a dead agent cannot report recovery
- **Entra-only workspaces**: add an `identity` and grant it workspace read access

## Key Outputs

| Output | Use |
| --- | --- |
| `scheduled_query_alert_id` | The rule's ARM ID |
| `identity_principal_id` | Grant target when the workspace enforces Entra-only queries |

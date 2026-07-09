---
title: "Monitor Metric Alert"
description: "Monitor Metric Alert deployment documentation"
icon: "package"
order: 100
componentName: "azuremonitormetricalert"
---

# Azure Monitor Metric Alert

Deploys an Azure Monitor metric alert rule -- static thresholds, machine-learned dynamic thresholds, or Application Insights web-test availability conditions over platform metrics, firing action groups when they hold.

## What Gets Created

When you deploy an AzureMonitorMetricAlert resource, Planton provisions:

- **Metric alert rule** -- an `azurerm_monitor_metric_alert` (global) carrying the scopes, exactly one condition family, the evaluation cadence, and the action-group wiring
- **Azure Tags** -- Planton-derived governance tags merged with your own (your values win)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource to monitor** -- any kind's `*_id` output becomes a scope
- **An AzureMonitorActionGroup** to fire into

## Quick Start

Create a file `metric-alert.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorMetricAlert
metadata:
  name: storage-availability
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMonitorMetricAlert.storage-availability
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: observability-rg
      fieldPath: status.outputs.resource_group_name
  alertName: storage-availability-low
  scopes:
    - valueFrom:
        kind: AzureStorageAccount
        name: app-storage
        fieldPath: status.outputs.storage_account_id
  severity: 1
  staticCriteria:
    - metricNamespace: Microsoft.Storage/storageAccounts
      metricName: Availability
      aggregation: AVERAGE
      operator: LESS_THAN
      threshold: 99.9
  actions:
    - actionGroupId:
        valueFrom:
          kind: AzureMonitorActionGroup
          name: platform-oncall
          fieldPath: status.outputs.action_group_id
```

Deploy it:

```bash
planton pulumi up --manifest metric-alert.yaml
```

## Common Configurations

- **Dynamic anomaly detection**: replace `staticCriteria` with `dynamicCriteria` (operator `GREATER_OR_LESS_THAN`, an `alertSensitivity`) -- Azure learns the normal band, seasonality included
- **Web-test availability**: `webTestAvailabilityCriteria` with the test's ARM ID, the Application Insights FK, and a `failedLocationCount`
- **Dimension splits**: add `dimensions` with `INCLUDE` + `"*"` to alert per dimension value independently
- **Multi-resource rules**: several scopes (or a resource group / subscription scope) plus `targetResourceType` and `targetResourceLocation`

## Key Outputs

| Output | Use |
| --- | --- |
| `metric_alert_id` | The rule's ARM ID |

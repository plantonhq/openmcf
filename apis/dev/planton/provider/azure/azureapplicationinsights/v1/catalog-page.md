# Azure Application Insights

Deploys a workspace-based Azure Application Insights resource -- the APM layer for web apps, functions, and container apps. Telemetry lands in a referenced Log Analytics Workspace; the resource's connection string is what applications are configured with.

## What Gets Created

When you deploy an AzureApplicationInsights resource, Planton provisions:

- **Application Insights component** -- an `azurerm_application_insights` bound to your Log Analytics Workspace, carrying the application type, retention, sampling, cost caps, and privacy/auth/network posture
- **Azure Tags** -- Planton-derived governance tags merged with your own (your values win)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureLogAnalyticsWorkspace** to store the telemetry (classic, workspace-less mode was retired by Azure) -- referenced through `workspaceId`
- **An Azure Resource Group** (can reference an AzureResourceGroup resource)

## Quick Start

Create a file `app-insights.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationInsights
metadata:
  name: web-app-insights
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureApplicationInsights.web-app-insights
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: observability-rg
      fieldPath: status.outputs.resource_group_name
  applicationInsightsName: web-app-insights
  applicationType: WEB
  workspaceId:
    valueFrom:
      kind: AzureLogAnalyticsWorkspace
      name: platform-logs
      fieldPath: status.outputs.workspace_id
```

Deploy it:

```bash
planton pulumi up --manifest app-insights.yaml
```

## Common Configurations

- **Production cost control**: `samplingPercentage: 50` + `dailyDataCapInGb: 10` -- representative telemetry at a fraction of the volume; the cap email warns before data drops silently
- **Keyless ingestion**: `localAuthenticationEnabled: false` -- SDKs authenticate with Entra identities; a bare instrumentation key stops working
- **Private-link only**: `internetIngestionEnabled: false` + `internetQueryEnabled: false` (requires Azure Monitor Private Link Scope)
- **Real client IPs**: `ipMaskingEnabled: false` -- only with privacy review; Azure masks to 0.0.0.0 by default

## Key Outputs

| Output | Use |
| --- | --- |
| `connection_string` | What apps are configured with -- Function Apps, Web Apps, and Container Apps reference it |
| `application_insights_id` | The ARM ID web-test metric alerts and diagnostic settings target |
| `app_id` | The identifier for REST-API telemetry queries |

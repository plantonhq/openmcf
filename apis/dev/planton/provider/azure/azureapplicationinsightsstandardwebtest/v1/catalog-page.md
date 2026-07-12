# Azure Application Insights Standard Web Test

Creates an Application Insights Standard Web Test — a synthetic availability monitor that pings a URL from one or more Azure regions on a schedule and records the endpoint's health into an Application Insights component. Pair it with a metric alert to get paged when the endpoint goes down.

## What Gets Created

When you deploy an AzureApplicationInsightsStandardWebTest resource, Planton provisions:

- **Standard Web Test** — an `azurerm_application_insights_standard_web_test` bound to the referenced Application Insights component, running the configured request from the configured regions

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Application Insights component** (an `AzureApplicationInsights`) to store results
- **A resource group** (an `AzureResourceGroup`)
- **Monitoring write rights**: `Microsoft.Insights/webTests/write`

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationInsightsStandardWebTest
metadata:
  name: homepage-health
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureApplicationInsightsStandardWebTest.homepage-health
spec:
  resourceGroup:
    value: platform-rg
  name: homepage-health
  region: eastus
  applicationInsightsId:
    valueFrom:
      name: platform-appinsights
  request:
    url: https://www.example.com/health
  geoLocations:
    - us-tx-sn1-azr
    - us-va-ash-azr
```

Deploy:

```shell
planton apply -f web-test.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `resourceGroup` | `StringValueOrRef` | Resource group. Defaults to an `AzureResourceGroup` reference. |
| `name` | `string` | Test name, unique in the resource group. Fixed at creation. |
| `region` | `string` | The test resource's region (usually the component's). Fixed at creation. |
| `applicationInsightsId` | `StringValueOrRef` | The component storing results. Defaults to `AzureApplicationInsights`. Fixed at creation. |
| `request` | `object` | The HTTP request; `url` required, plus `httpVerb`, `body`, `headers`, redirect/dependent-request toggles. |
| `geoLocations` | `string[]` | Azure web-test location ids to run FROM (at least one). |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `frequency` | `int32` | `300` | Run interval in seconds (300 / 600 / 900). |
| `timeout` | `int32` | `30` | Per-run timeout in seconds. |
| `enabled` | `bool` | `true` | Whether the test runs. |
| `retryEnabled` | `bool` | provider default | Retry a failed run before counting a failure. |
| `validationRules` | `object` | — | `expectedStatusCode`, `sslCheckEnabled`, `sslCertRemainingLifetime`, and a `content` match block. |
| `description` | `string` | `""` | Human-readable description. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins). |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `web_test_id` | `string` | Full ARM ID — referenced by `AzureMonitorMetricAlert` web-test criteria |
| `web_test_name` | `string` | The test's name as deployed |
| `synthetic_monitor_id` | `string` | The synthetic monitor id Azure assigns |

## Related Components

- [AzureApplicationInsights](/docs/catalog/azure/application-insights) — the component that stores the test's results
- [AzureMonitorMetricAlert](/docs/catalog/azure/monitor-metric-alert) — alerts on the test's availability via `webTestId`
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group

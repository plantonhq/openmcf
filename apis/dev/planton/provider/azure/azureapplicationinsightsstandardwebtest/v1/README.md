# AzureApplicationInsightsStandardWebTest

## Overview

`AzureApplicationInsightsStandardWebTest` provisions an Application Insights
Standard Web Test: a synthetic availability monitor that issues an HTTP
request to a URL from one or more Azure regions on a schedule and records
whether the endpoint responded correctly. Paired with a metric alert on the
test's availability, it is how you get paged when an endpoint goes down.

## Why a First-Class Resource?

- **Outside-in availability** -- proves an endpoint is reachable and healthy
  from real Azure regions, not just from inside your network
- **The metric alert's target** -- `AzureMonitorMetricAlert`'s web-test
  availability criteria reference this test by ARM ID
- **Bound to a component** -- results and the availability metric live in an
  `AzureApplicationInsights` component

## Key Features

- **Multi-region** -- run from several `geo_locations` so one region's
  network blip does not read as an outage
- **Schedule + timeout + retry** -- 5/10/15-minute frequency, a per-run
  timeout, and optional retry to suppress false alarms
- **Full request control** -- verb, headers, body, redirect and
  dependent-request behavior
- **Response validation** -- expected status code, SSL certificate lifetime
  floor, and body-content match (pass or fail on found text)
- **Composable** -- component and resource group are references defaulting
  to their Planton kinds

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `resource_group` | StringValueOrRef | Yes | Resource group (defaults to AzureResourceGroup) |
| `name` | string | Yes | Test name, unique in the resource group; fixed at creation |
| `region` | string | Yes | The test resource's region (usually the component's); fixed at creation |
| `application_insights_id` | StringValueOrRef | Yes | The component (defaults to AzureApplicationInsights); fixed at creation |
| `frequency` | int32 | No | 300 / 600 / 900 seconds (default 300) |
| `timeout` | int32 | No | Per-run timeout in seconds (default 30) |
| `enabled` | bool | No | Active (default true) |
| `retry_enabled` | bool | No | Retry a failed run before counting a failure |
| `request` | message | Yes | The HTTP request (url required; verb/headers/body/redirects) |
| `validation_rules` | message | No | Status code / SSL / content assertions |
| `geo_locations` | repeated string | Yes (≥1) | Azure web-test location ids to run FROM |
| `description` | string | No | Human-readable description |
| `tags` | map | No | User tags, merged over Planton-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `web_test_id` | Full ARM ID -- the seam a metric alert references |
| `web_test_name` | The test's name as deployed |
| `synthetic_monitor_id` | The synthetic monitor id Azure assigns |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationInsightsStandardWebTest
metadata:
  name: homepage-health
spec:
  resourceGroup:
    valueFrom:
      name: platform-rg
  name: homepage-health
  region: eastus
  applicationInsightsId:
    valueFrom:
      name: platform-appinsights
  request:
    url: https://www.example.com/health
  validationRules:
    expectedStatusCode: 200
    content:
      contentMatch: "healthy"
  geoLocations:
    - us-tx-sn1-azr
    - us-va-ash-azr
```

## Lifecycle Notes

- `name`, `region`, and `application_insights_id` are **fixed at creation**
- `geo_locations` are Azure web-test location ids (e.g. `us-tx-sn1-azr`),
  not standard region names
- Setting `validation_rules.ssl_cert_remaining_lifetime` also enables the
  SSL check

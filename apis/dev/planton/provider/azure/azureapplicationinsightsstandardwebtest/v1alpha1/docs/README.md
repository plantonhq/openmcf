# AzureApplicationInsightsStandardWebTest -- Design Research

## The Resource

An Application Insights Standard Web Test
(`Microsoft.Insights/webTests`, "standard" kind) is a synthetic
availability monitor: it issues an HTTP request to a URL from Azure regions
on a schedule and records success/failure into an Application Insights
component. The component maps onto
`azurerm_application_insights_standard_web_test` (azurerm v4.x,
`internal/services/applicationinsights/application_insights_standard_web_test_resource.go`),
parity-verified against pulumi-azure v6 (`appinsights.StandardWebTest`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `application_insights_id` | `application_insights_id` | FK → AzureApplicationInsights, ForceNew |
| `location` | `region` | Required, ForceNew |
| `frequency` | `frequency` | 300/600/900; default 300 |
| `timeout` | `timeout` | default 30 |
| `enabled` / `retry_enabled` | `enabled` / `retry_enabled` | |
| `request` (block, MaxItems 1) | `request` message | url required; verb/headers/body/toggles |
| `validation_rules` (block, MaxItems 1) | `validation_rules` message | status/ssl/content |
| `geo_locations` (list, MinItems 1) | `geo_locations` | Azure web-test location ids |
| `description` | `description` | |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `id` (computed) | `web_test_id` output | Metric-alert seam |
| `synthetic_monitor_id` (computed) | `synthetic_monitor_id` output | |

## Key Design Decisions

- **Classic web test is the recorded skip.** azurerm also has
  `azurerm_application_insights_web_test` (the legacy Ping/Multistep test
  configured via an XML `configuration` blob). The standard web test is its
  modern, structured successor -- the fields are first-class instead of raw
  XML -- so Planton models only the standard one. The classic kind is a
  deliberate omission (its XML-blob shape is hostile to a self-service
  form).
- **frequency and http_verb are validated in-spec.** The provider's
  IntInSlice{300,600,900} and StringInSlice verb list become CELs, so the
  error surfaces at validation rather than apply. Modeled as scalar fields
  (not enums) because they are numeric/verb literals the provider already
  round-trips as-is.
- **The metric-alert seam is the reason this ships now.** The metric alert's
  web-test availability criteria previously carried a plain-string
  `web_test_id`; forging this kind converts that to a
  `default_kind`-annotated FK, so an availability alert composes cleanly
  onto a Planton-managed test.

## Composition Seams

The test consumes `AzureApplicationInsights` (the component) and
`AzureResourceGroup` (its RG). Its `web_test_id` output is referenced by
`AzureMonitorMetricAlert.web_test_availability_criteria.web_test_id`.

## Live E2E

Live dual-engine E2E deploys the fixture Application Insights component
(registry prerequisite, bringing its workspace and RG), then a web test
pinging microsoft.com from one region -- a public endpoint, so the run
needs no subscription-owned URL. Web tests provision in seconds; the run
verifies via a generic ARM GetByID and tears down cleanly.

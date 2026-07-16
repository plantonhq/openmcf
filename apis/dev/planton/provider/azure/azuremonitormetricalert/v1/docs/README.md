# AzureMonitorMetricAlert -- Design Research

## The Resource

A metric alert (`Microsoft.Insights/metricAlerts`) evaluates platform
metrics against static, dynamic, or web-test conditions. The component maps
onto `azurerm_monitor_metric_alert` (azurerm v4.x,
`internal/services/monitor/monitor_metric_alert_resource.go`),
parity-verified against pulumi-azure v6 (`monitoring.MetricAlert`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `alert_name` | ForceNew |
| `scopes` | `scopes[]` | Repeated bare StringValueOrRef, NO default kind (any metric-emitting resource; the diagnostic-setting target precedent). Min 1 |
| `criteria` | `static_criteria[]` | Multiple allowed (AND); the shared operator enum restricted to the five directional comparisons by CEL |
| `dynamic_criteria` | `dynamic_criteria` | MaxItems 1 in the provider -- a singular message here; operator restricted to GREATER_THAN/LESS_THAN/GREATER_OR_LESS_THAN by CEL |
| `application_insights_web_test_location_availability_criteria` | `web_test_availability_criteria` | component_id FK -> AzureApplicationInsights.application_insights_id; web_test_id a literal ARM id (availability tests are not a catalog kind) |
| `frequency` / `window_size` | same | The provider's exact ISO 8601 vocabularies as in-list rules on strings (the catalog's duration-vocabulary pattern); defaults PT1M / PT5M |
| `severity` | same | 0-4, default 3 |
| `auto_mitigate` | same | optional-bool default true |
| `target_resource_type` / `target_resource_location` | same | Required by Azure for multi-resource/group/subscription scopes -- an apply-time contract documented on the fields (whether it applies depends on the resolved scope shape) |
| `action` | `actions[]` | action_group_id FK -> AzureMonitorActionGroup.action_group_id + webhook_properties map |
| `tags` | `tags` | User tags merged over metadata tags |

## The Shared Operator Enum

Static and dynamic criteria overlap on GREATER_THAN/LESS_THAN but diverge on
the rest (EQUALS et al. are static-only; GREATER_OR_LESS_THAN dynamic-only).
Proto enum values share a namespace within the package, so ONE
`AzureMonitorMetricAlertOperator` enum carries all six values and
per-criteria CELs enforce the split -- the Front Door rule-set
shared-operator precedent.

## Front-Loaded Contracts (provider-verified)

- **Exactly one condition family** -- the provider's `ExactlyOneOf` across
  the three criteria fields, as a message CEL.
- **Operator direction splits** -- the provider's per-family
  `StringInSlice` sets, as per-field CELs.
- **Dimension values min 1** -- the provider's MinItems.

## Deliberate Skips (recorded reasons)

- **Legacy single-resource criteria wire type** -- the provider silently
  keeps `MetricAlertSingleResourceMultipleMetricCriteria` for pre-existing
  rules (backward compatibility); new rules always use the multi-resource
  type. Not a spec concern -- both engines create new rules.
- **`evaluation_failure_count <= evaluation_total_count`** (dynamic) --
  ARM's contract, but the provider does not validate it statically;
  documented on the fields rather than invented as CEL.

## Design Notes

- `threshold` is a plain double: zero is a meaningful threshold ("any failed
  request") and the field is required within its criterion.
- Rules are global (no location); the provider stamps "Global" server-side.

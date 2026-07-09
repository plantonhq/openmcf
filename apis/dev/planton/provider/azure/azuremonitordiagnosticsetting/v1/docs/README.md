# AzureMonitorDiagnosticSetting -- Design Research

## The Resource

A diagnostic setting (`Microsoft.Insights/diagnosticSettings`) is an ARM
EXTENSION resource: it lives on its target and routes the target's platform
logs/metrics to destinations. The component maps onto
`azurerm_monitor_diagnostic_setting` (azurerm v4.x,
`internal/services/monitor/monitor_diagnostic_setting_resource.go`),
parity-verified against pulumi-azure v6 (`monitoring.DiagnosticSetting`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `setting_name` | The provider's forbidden-character set (`<>*%&:\?+/`) as CEL. ForceNew |
| `target_resource_id` | `target_resource_id` | Bare StringValueOrRef, NO default kind -- any ARM resource can be a target and none dominates (the VM disk-encryption-set precedent). ForceNew |
| `enabled_log` | `enabled_logs[]` | category XOR category_group per entry -- the provider enforces it at expand time; front-loaded as a message CEL |
| `enabled_metric` | `enabled_metrics[]` | category required |
| `log_analytics_workspace_id` | same | FK -> AzureLogAnalyticsWorkspace.workspace_id |
| `log_analytics_destination_type` | same | Closed DEDICATED/AZURE_DIAGNOSTICS enum; absent lets Azure pick (some resource types support only one layout) |
| `storage_account_id` | same | FK -> AzureStorageAccount.storage_account_id |
| `eventhub_authorization_rule_id` | same | Literal ARM id -- a NAMESPACE-level authorization rule; no auth-rule kind exists (the Event Hub namespace kind predates its rework) |
| `eventhub_name` | same | Rides with the rule (CEL); empty = Azure's default hub per category |
| `partner_solution_id` | same | Literal ARM id (Azure Native ISV resources are outside the catalog) |

## Front-Loaded Contracts (all provider-verified)

- **At-least-one destination** -- the provider's `AtLeastOneOf` across the
  four destination fields, as a message CEL.
- **At-least-one category** -- the provider re-enforces at create AND update
  ("at least one type of Log or Metric must be enabled") because the API
  "creates" an empty setting and then 404s on read.
- **category XOR category_group** -- the provider's expand-time error.

## Deliberate Skips (recorded reasons)

- **`metric` block + `retention_policy`** -- deprecated (removed in azurerm
  v5; retention moved to `azurerm_storage_management_policy`); only the
  v5-clean `enabled_log`/`enabled_metric` shapes are modeled.
- **Tags** -- ARM does not support tags on diagnostic settings.

## The Composite-ID Trap

The provider's state ID is `"{target_resource_id}|{name}"` -- NOT an ARM id;
no Azure API consumes it. Both modules therefore CONSTRUCT the real
extension-resource id
(`{target}/providers/Microsoft.Insights/diagnosticSettings/{name}`) for the
`diagnostic_setting_id` output, which is also what the E2E verifier reads.

## Provider Behaviors Leaned On

- Azure's API is eventually consistent here: the provider polls after create
  and delete until the setting is readable/gone.
- `log_analytics_destination_type` is unconfigurable for some resource types
  (for example Key Vault) -- Azure decides; documented on the field.

# AzureMonitorScheduledQueryAlert -- Design Research

## The Resource

A scheduled query rule (`Microsoft.Insights/scheduledQueryRules`, kind
LogAlert) runs KQL on a schedule and alerts on the result. The component
maps onto `azurerm_monitor_scheduled_query_rules_alert_v2` (azurerm v4.x,
`internal/services/monitor/monitor_scheduled_query_rules_alert_v2_resource.go`),
parity-verified against pulumi-azure v6
(`monitoring.ScheduledQueryRulesAlertV2`).

## Kind Naming

The kind drops the provider's "v2" artifact: ARM's own type is
scheduledQueryRules, and the v1 provider resource is a superseded older-API
shape (below), not a different Azure resource.

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `alert_name` | ForceNew |
| `location` | `region` | REQUIRED -- the rule must live in the queried resource's region (unlike metric alerts, which are global) |
| `scopes` (MaxItems 1, ForceNew) | `scope` | A singular reference; FK-defaults to AzureLogAnalyticsWorkspace.workspace_id (query alerts overwhelmingly target workspaces; App Insights scopes via explicit valueFrom) |
| `severity` | same | 0-4; the provider requires it -- middleware default 3 keeps the friendly-default grain |
| `evaluation_frequency` / `window_duration` | same | The provider's exact ISO 8601 vocabularies as in-list rules; middleware default PT5M (the portal's default cadence) |
| `criteria` | `criteria[]` | Min 1; per-criterion notes below |
| `query_time_range_override` | same | The provider's exact vocabulary as CEL |
| `auto_mitigation_enabled` | same | plain bool (provider default false) |
| `mute_actions_after_alert_duration` | same | Vocabulary CEL; XOR with auto-mitigation (the provider's create/update error, front-loaded as message CEL) |
| `workspace_alerts_storage_enabled` | same | plain bool |
| `skip_query_validation` | same | plain bool -- needed for custom-log tables that appear only after data flows |
| `identity` | `identity` | SystemAssigned / UserAssigned (no combined model on this resource) |
| `target_resource_types` | same | Pairs with a criterion's resource_id_column |
| `action` | `action` | action_groups FKs + custom_properties + email_subject (serialized to ActionProperties["Email.Subject"] by the provider) |
| `tags` | `tags` | User tags merged over metadata tags |

## Criterion Notes

- **Operator vocabulary** -- this API spells equality "Equal" (the provider
  fixed "Equals" as a breaking bug); the enum value is EQUAL and both
  modules' maps carry the exact wire string.
- **metric_measure_column** -- required for non-COUNT aggregations and
  forbidden for COUNT (Azure's apply-time pairing; the provider only
  documents it) -- documented on the field, not invented as CEL.
- **failing_periods** -- both ints 1-6 (provider-validated); the
  min <= total bound is front-loaded as a message CEL (a rule violating it
  can never fire -- Azure rejects it).

## Deliberate Skips (recorded reasons)

- **`azurerm_monitor_scheduled_query_rules_alert` (v1)** -- the older
  2018-04-16 API shape, documented by the provider as superseded by v2 and
  known to cause problems; only v2 is modeled.
- **`azurerm_monitor_scheduled_query_rules_log` (metric-measurement v1
  sibling)** -- same superseded family.
- **Computed read-backs** (`created_with_api_version`,
  `is_a_legacy_log_analytics_rule`, `is_workspace_alerts_storage_configured`)
  -- server-side reflection with no configuration value.

## Cross-Resource Contracts (documented, not CEL'd)

- `evaluation_frequency` must not exceed `window_duration` (and not exceed
  the mute duration when set) -- Azure validates at apply.
- Entra-only workspaces require the rule's identity to hold read access on
  the workspace -- granted through `AzureRoleAssignment`.

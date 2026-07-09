# AzureApplicationInsights -- Design Research

## The Resource

An Application Insights component (`Microsoft.Insights/components`) is
Azure's APM store and portal experience. The component maps onto
`azurerm_application_insights` (azurerm v4.x,
`internal/services/applicationinsights/application_insights_resource.go`),
parity-verified against pulumi-azure v6 (`appinsights.Insights`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `application_insights_name` | Kind-authentic rename; 1-260. ForceNew |
| `application_type` | `application_type` | Closed 8-value enum; unspecified deploys "web". Azure's strings are CASE-SENSITIVE and irregular ("Node.JS", "MobileCenter") and an unmatched value is silently treated as ASP.NET -- the closed vocabulary makes that class unrepresentable. ForceNew |
| `workspace_id` | `workspace_id` | REQUIRED here though the provider marks it Optional+Computed: classic mode was retired Feb 2024, Azure auto-attaches a default workspace to unmigrated legacy resources, and an explicit binding is the only honest new-resource shape. Repointable, never removable (provider-enforced) |
| `retention_in_days` | same | Closed 9-value set; default 90 |
| `daily_data_cap_in_gb` | same | >= 0, default 100; applied through a SEPARATE billing API in a follow-up call by the provider |
| `daily_data_cap_notifications_enabled` | same | v5-positive form; the wire property is the inverted StopSendNotificationWhenHitCap |
| `sampling_percentage` | same | 0-100, default 100 |
| `ip_masking_enabled` | same | v5-positive form of `disable_ip_masking` |
| `local_authentication_enabled` | same | v5-positive form of `local_authentication_disabled` |
| `internet_ingestion_enabled` / `internet_query_enabled` | same | optional-bool default true |
| `force_customer_storage_for_profiler` | same | plain bool |
| `tags` | `tags` | User tags merged over metadata tags |

## Deliberate Skips (recorded reasons)

- **Deprecated negative-form aliases** (`disable_ip_masking`,
  `local_authentication_disabled`, `daily_data_cap_notifications_disabled`)
  -- removed in azurerm v5; only the positive forms are modeled.

## PARITY-EXCEPTION (bridge lag)

pulumi-azure v6.38 bridges ONLY the deprecated negative-form booleans
(`disableIpMasking`, `localAuthenticationDisabled`,
`dailyDataCapNotificationsDisabled`). The Pulumi module inverts the spec's
positive booleans; the wire property is identical for each pair, behavior
and outputs match exactly. Documented in the Pulumi module with the
re-alignment trigger (the bridge shipping the positive forms).

## Provider Behaviors Leaned On

- Azure auto-creates a noisy "Failure Anomalies" smart-detector rule with
  every component; the provider disables it by default (deleting it just
  resurrects it server-side).
- The data-cap pair lives on a separate billing API -- expect the extra
  follow-up call at create/update.

## Outputs Design

`connection_string` is the composition seam (`AzureFunctionApp`,
`AzureLinuxWebApp`, and `AzureContainerAppEnvironment` reference it);
`instrumentation_key` remains for unmigrated SDKs. Both are secret-bearing
outputs (they authorize ingestion while local auth is enabled), documented
in prose per the established access-key precedent.

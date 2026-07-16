# AzureMonitorActionGroup -- Design Research

## The Resource

An action group (`Microsoft.Insights/actionGroups`) is the notification hub
of Azure Monitor alerting. The component maps onto
`azurerm_monitor_action_group` (azurerm v4.x,
`internal/services/monitor/monitor_action_group_resource.go`),
parity-verified against pulumi-azure v6 (`monitoring.ActionGroup`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | metadata.name | The resource name basis, per the catalog grain |
| `short_name` | `short_name` | 1-12 chars (provider-validated length as CEL); the SMS sender identity |
| `enabled` | `enabled` | optional-bool default true |
| `location` | -- | CONSTANT "global" (recorded skip below) |
| receiver blocks | `*_receivers[]` | All 11 families modeled; per-family notes below |
| `tags` | `tags` | User tags merged over metadata tags |

## Receiver-Family Notes (provider-verified)

- **`use_common_alert_schema`** exists on exactly seven families (email,
  webhook, automation runbook, logic app, function, ARM role, event hub) --
  SMS, voice, app push, and ITSM are not payload-aware; the spec models the
  flag only where it exists.
- **Webhook `aad_auth`** -- object_id required (UUID CEL); identifier_uri and
  tenant_id optional (Azure derives them). The keyless-webhook posture.
- **ITSM `ticket_configuration`** -- the provider's expand rejects JSON
  lacking the PayloadRevision and WorkItemType keys; front-loaded as a
  contains-check CEL (CEL cannot parse JSON; the check catches the
  documented failure class).
- **ITSM `workspace_id`** -- a `{subscription}|{workspace_customer_id}`
  composite (Azure's ITSM addressing), NOT an ARM id; documented on the
  field.
- **ARM-role `role_id`** -- FK-defaults to
  `AzureRoleDefinition.role_definition_guid`; built-in roles by well-known
  literal GUID (the vendor-catalog look-up-never-infer class).
- **Event Hub receiver** -- addressed by namespace NAME + hub NAME (not ARM
  id); the namespace FK-defaults to
  `AzureEventHubNamespace.namespace_name`; tenant/subscription optional
  UUIDs (Azure defaults them to the home tenant/subscription).
- **Function receiver** -- `function_app_resource_id` FK-defaults to
  `AzureFunctionApp.function_app_id`.

## Deliberate Skips (recorded reasons)

- **`location`** -- the provider defaults to "global" and global is the only
  value that makes sense for alerting infrastructure (regional action groups
  exist for niche data-residency cases); a one-value knob is a constant.
  Both modules let the provider default apply. Re-enable trigger: a real
  data-residency requirement surfacing.

## Design Notes

- A receiver-less action group is legal and useful (a "null" routing target
  alert rules can point at before channels are decided) -- no invented
  min-receivers rule.
- Receiver names must be unique within the group case-insensitively (an
  Azure-side contract, documented rather than CEL'd -- it spans eleven
  repeated fields).

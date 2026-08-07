# AzureMonitorActivityLogAlert -- Design Research

## The Resource

An Azure Monitor Activity Log Alert
(`Microsoft.Insights/activityLogAlerts`) fires actions when a matching entry
appears in the subscription's Activity Log. The component maps onto
`azurerm_monitor_activity_log_alert` (azurerm v4.x,
`internal/services/monitor/monitor_activity_log_alert_resource.go`),
parity-verified against pulumi-azure v6 (`monitoring.ActivityLogAlert`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `location` | `location` enum | global/westeurope/northeurope/eastus2euap only (CustomizeDiff) |
| `scopes` (set) | `scopes` | FK list; polymorphic (default AzureResourceGroup) |
| `criteria` (block, MaxItems 1) | `criteria` | category + narrowing fields |
| `action` (block list) | `actions` | action_group_id FK + webhook_properties |
| `description` | `description` | |
| `enabled` | `enabled` | Azure default true |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## Key Design Decisions

- **Singular forms folded into plural.** azurerm exposes both `level`/
  `levels`, `resource_provider`/`resource_providers`, etc., as mutually
  exclusive `ConflictsWith` pairs. The spec keeps only the plural form (a
  one-element list expresses the single case), eliminating an entire class
  of conflict validation and reading cleaner. Every list field is optional
  (empty = match any).
- **The location allowlist is an enum, not a free region.** azurerm's
  CustomizeDiff rejects any region outside
  {global, westeurope, northeurope, eastus2euap}; modeling location as a
  four-value enum (defaulting GLOBAL) makes the constraint unmissable and
  removes a whole failure mode. Nearly every alert uses GLOBAL -- the alert
  evaluates the subscription-global Activity Log regardless of where its
  definition lives.
- **The three exclusivity CELs mirror the provider.** `caller`,
  `resource_health`, and `service_health` are mutually exclusive in the
  provider (a caller has no meaning for a health event); the spec
  front-loads that as message CELs (they dereference only scalar/message
  presence, not StringValueOrRef sub-fields, so they are expressible), plus
  the recommendation_type-vs-category/impact exclusivity.
- **Category-specific blocks stay nested.** `resource_health` and
  `service_health` are their own messages because they carry structured
  sub-fields (state sets, event/location/service lists) that only apply to
  their category; folding them into the flat criteria would muddle the
  contract.

## Composition Seams

The alert consumes `AzureResourceGroup` (its own RG + the common scope),
any resource id as a scope (polymorphic), and `AzureMonitorActionGroup`
(each action). Its `activity_log_alert_id` output is available for
reference.

## Live E2E

Live dual-engine E2E deploys the fixture action group (registry
prerequisite, bringing the fixture RG), then an administrative alert scoped
to the fixture RG firing that action group. The alert is global metadata
and provisions in seconds; verification is a generic ARM GetByID.

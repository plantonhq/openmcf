# AzureLogAnalyticsWorkspace -- Design Research

## The Resource

A Log Analytics Workspace (`Microsoft.OperationalInsights/workspaces`) is
Azure Monitor's log store: the ingestion target for diagnostic settings and
agents, the KQL query surface, and the scope log alerts evaluate against.
The component maps onto `azurerm_log_analytics_workspace` (azurerm v4.x,
`internal/services/loganalytics/log_analytics_workspace_resource.go`),
parity-verified against pulumi-azure v6
(`operationalinsights.AnalyticsWorkspace`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `workspace_name` | Kind-authentic rename; the provider's 4-63 letters/digits/hyphens rule as CEL. ForceNew |
| `sku` | `sku` | Closed 4-value enum; unspecified deploys PerGB2018 (the provider's own default when omitted) |
| `reservation_capacity_in_gb_per_day` | same | Closed 11-value set; CapacityReservation pairing enforced both directions as message CELs (the provider enforces it at create AND update) |
| `retention_in_days` | same | 30-730; middleware default 30 |
| `daily_quota_gb` | same | -1 sentinel = unlimited; always sent so both engines match |
| `identity` | `identity` | SystemAssigned XOR UserAssigned -- the workspace REJECTS the combined model (caught by the offline plan; the schema's generic identity block over-promises) |
| `local_authentication_enabled` | same | optional-bool default true (v5-positive form; the deprecated `local_authentication_disabled` is not modeled) |
| `internet_ingestion_enabled` / `internet_query_enabled` | same | optional-bool default true |
| `allow_resource_only_permissions` | same | optional-bool default true; the provider applies it via a follow-up update + poll (an Azure REST quirk) |
| `cmk_for_query_forced` | same | plain bool |
| `immediate_data_purge_on_30_days_enabled` | same | plain bool |
| `data_collection_rule_id` | same | Literal ARM id (no DCR kind exists); the provider applies it via a follow-up update (ARM rejects it at create) |
| `tags` | `tags` | User tags merged over metadata tags |

## Deliberate Skips (recorded reasons)

- **`Standard` / `Premium` skus** -- the provider's CustomizeDiff blocks creating
  workspaces on them ("no longer supported by Azure"); modeling them would ship
  guaranteed apply failures.
- **`LACluster` sku** -- a server-side state Azure sets when the workspace links
  to a dedicated Log Analytics cluster (creation is blocked outside a
  soft-delete-relink edge case); not a tier a user chooses. Re-enable trigger:
  a dedicated-cluster kind joining the catalog.
- **`Unlimited` sku** -- retired legacy tier (removed from the service's SDK; the
  provider carries it with a "may no longer be valid" TODO).
- **`local_authentication_disabled`** -- the deprecated negative-form alias
  (removed in azurerm v5); only the positive form is modeled.

## SKU Transition Semantics (documented on the spec)

The provider's CustomizeDiff makes SKU changes ForceNew EXCEPT
PerGB2018 <-> CapacityReservation (in-place, Azure's commitment entry/exit)
and transitions out of LACluster. Raising a commitment tier restarts Azure's
31-day commitment period.

## Outputs Design

`workspace_id` carries the ARM resource ID -- the seam every FK consumer
(App Insights, AKS addons, Container App Environments, diagnostic settings,
scheduled query alerts) references, and what `ValidateWorkspaceID` accepts.
The provider's `workspace_id` ATTRIBUTE is the customer GUID agents
authenticate against -- exported as `workspace_customer_id` so the two can
never be confused. Shared keys are secret-bearing outputs (documented in
prose; unusable as credentials under the Entra-only posture).

## Cross-Resource Contracts (apply-time, documented not CEL'd)

- Private-link-only posture (`internet_*_enabled: false`) requires an Azure
  Monitor Private Link Scope wired to the workspace -- AMPLS is a separate
  resource family outside this kind.
- The identity's usefulness (cluster CMK, linked storage) depends on grants
  made through `AzureRoleAssignment`.
